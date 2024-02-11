package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var accessTokenSecret = viper.GetString("JTW_ACCESS_SECRET")
var refreshTokenSecret = viper.GetString("JTW_REFRESH_SECRET")

type TokenDetail struct {
	ID       int
	Email    string
	Username string
}

func GenerateAccessToken(td *TokenDetail) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 3)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = td.ID
	claims["Email"] = td.Email
	claims["Username"] = td.Username
	claims["exp"] = expiry

	signedToken, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}

func GenerateRefreshToken(td *TokenDetail) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 7)

	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = td.ID
	claims["exp"] = expiry

	signedToken, err := token.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}
