package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/foolish06/gin-essential/model"
	"time"
)

var jwtKey = []byte("wxy-secret_token")

type MyClaims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken (user model.User) (string, error){
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &MyClaims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "wxy",
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, *MyClaims, error){
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}