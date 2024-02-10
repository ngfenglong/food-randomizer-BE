package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ngfenglong/food-randomizer-BE/pkg/config"
	"github.com/ngfenglong/food-randomizer-BE/pkg/utils"
)

func Login(repo AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginCredential LoginDto
		err := json.NewDecoder(r.Body).Decode(&loginCredential)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		user, err := repo.GetUserByEmail(ctx, loginCredential.Email)
		if err != nil {
			utils.ErrorJSON(w, errors.New("invalid login crendential"), http.StatusUnauthorized)
			return
		}

		validPassword, err := utils.PasswordMatches(user.Password, loginCredential.Password)
		if !validPassword || err != nil {
			utils.ErrorJSON(w, errors.New("invalid login crendential"), http.StatusUnauthorized)
			return
		}

		tokenDetail := &TokenDetail{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.UserName,
		}
		accessToken, accessExpiry, err := GenerateAccessToken(tokenDetail)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		refreshToken, refreshExpiry, err := GenerateRefreshToken(tokenDetail)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		err = repo.InsertToken(ctx, user.ID, refreshToken, refreshExpiry)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		var payload LoginResponseDto
		payload.AccessToken = accessToken
		payload.RefreshToken = refreshToken
		payload.Expiry = accessExpiry
		payload.UserName = user.UserName

		err = utils.WriteJSON(w, http.StatusOK, payload, "data")
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func Logout(repo AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logoutRequestDto struct {
			RefreshToken string `json:"refresh_token"`
		}

		err := json.NewDecoder(r.Body).Decode(&logoutRequestDto)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Todo: Add Validation Handling
		err = repo.DeleteToken(ctx, logoutRequestDto.RefreshToken)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "Logout successful"

		err = utils.WriteJSON(w, http.StatusOK, payload, "data")
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	}
}

func ForgetPassword(repo AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func Register(repo AuthRepository, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerInput RegisterUserDto
		err := json.NewDecoder(r.Body).Decode(&registerInput)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		if registerInput.SecretCode != cfg.SecretCode {
			utils.ErrorJSON(w, errors.New("insert a valid secret code"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		cancel()

		usernameExists, emailExists, err := repo.CheckIfUserExists(ctx, registerInput)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		if usernameExists {
			utils.ErrorJSON(w, errors.New("usename already exists"), http.StatusConflict)
			return
		}

		if emailExists {
			utils.ErrorJSON(w, errors.New("email already exists"), http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), 10)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		registerInput.Password = string(hashedPassword)

		err = repo.RegisterUser(ctx, registerInput)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "User registered successfully"

		err = utils.WriteJSON(w, http.StatusOK, payload, "data")
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	}
}

func Request_access(repo AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var adminRequest AdminRequestDto
		err := json.NewDecoder(r.Body).Decode(&adminRequest)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		isRequestPending, err := repo.IsAdminRequestPending(ctx, adminRequest)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		if isRequestPending {
			utils.ErrorJSON(w, errors.New("An admin access request for this Telegram ID has already been submitted and is pending review."))
			return
		}

		err = repo.RegisterRequest(ctx, adminRequest.TelegramID, adminRequest.TelegramUsername)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		var payload struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}
		payload.Error = false
		payload.Message = "Request submitted successfully"

		err = utils.WriteJSON(w, http.StatusOK, payload, "data")
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	}
}
