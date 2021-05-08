package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type LoginClaims struct {
	UserName string
	UserID   string
	jwt.StandardClaims
}

func GenerateToken(userName string, userID string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LoginClaims{
		UserName: userName,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	})

	return token.SignedString([]byte(ManagerConfig.SecretKey))
}
