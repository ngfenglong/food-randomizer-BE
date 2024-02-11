package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ngfenglong/food-randomizer-BE/pkg/config"
	"github.com/ngfenglong/food-randomizer-BE/pkg/utils"
	"github.com/pascaldekloe/jwt"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func CheckToken(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// could set an annonymous user
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			utils.ErrorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			utils.ErrorJSON(w, errors.New("unauthorize - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(cfg.JWT.Secret))
		if err != nil {
			utils.ErrorJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			utils.ErrorJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			utils.ErrorJSON(w, errors.New("unauthorized - invalid audience"), http.StatusForbidden)
			return
		}

		if claims.Issuer != "mydomain.com" {
			utils.ErrorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			utils.ErrorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
		}

		log.Println("Valid user", userID)
		next.ServeHTTP(w, r)
	})
}
