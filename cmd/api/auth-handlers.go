package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var loginCredential models.LoginDto
	err := json.NewDecoder(r.Body).Decode(&loginCredential)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.DB.GetUserByEmail(loginCredential.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid login crendential"), http.StatusUnauthorized)
		return
	}

	validPassword, err := app.PasswordMatches(user.Password, loginCredential.Password)
	if !validPassword || err != nil {
		app.errorJSON(w, errors.New("invalid login crendential"), http.StatusUnauthorized)
		return
	}

	tokenDetail := &TokenDetail{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.UserName,
	}
	accessToken, accessExpiry, err := app.GenerateAccessToken(tokenDetail)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	refreshToken, refreshExpiry, err := app.GenerateRefreshToken(tokenDetail)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = app.models.DB.InsertToken(user.ID, refreshToken, refreshExpiry)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var payload models.LoginResponseDto
	payload.AccessToken = accessToken
	payload.RefreshToken = refreshToken
	payload.Expiry = accessExpiry
	payload.UserName = user.UserName

	err = app.WriteJSON(w, http.StatusOK, payload, "data")
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	var logoutRequestDto struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&logoutRequestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Todo: Add Validation Handling
	err = app.models.DB.DeleteToken(logoutRequestDto.RefreshToken)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = "Logout successful"

	err = app.WriteJSON(w, http.StatusOK, payload, "data")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) forgetPassword(w http.ResponseWriter, r *http.Request) {
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	fmt.Print("err 0", r.Body)
	var registerInput models.RegisterUserDto
	err := json.NewDecoder(r.Body).Decode(&registerInput)
	if err != nil {
		fmt.Print("err 1", err)
		app.errorJSON(w, err)
		return
	}

	if registerInput.SecretCode != app.config.secretCode {
		fmt.Print("err 2", err)
		app.errorJSON(w, errors.New("insert a valid secret code"))
		return
	}

	usernameExists, emailExists, err := app.models.DB.CheckIfUserExists(registerInput)
	if err != nil {
		fmt.Print("err 3", err)
		app.errorJSON(w, err)
		return
	}

	if usernameExists {
		app.errorJSON(w, errors.New("usename already exists"), http.StatusConflict)
		return
	}

	if emailExists {
		app.errorJSON(w, errors.New("email already exists"), http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), 10)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	registerInput.Password = string(hashedPassword)

	err = app.models.DB.RegisterUser(registerInput)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = "User registered successfully"

	err = app.WriteJSON(w, http.StatusOK, payload, "data")
	if err != nil {
		app.errorJSON(w, err)
	}
}
