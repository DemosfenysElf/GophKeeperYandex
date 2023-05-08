package service

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const keyJWT = "pqla3zxjonfgwouhf"
const (
	Authorization = "Authorization"
	Bearer        = "Bearer"
	User          = "user"
	Type          = "dataType"
	Card          = "card"
	Text          = "text"
	Password      = "password"
)

type ClaimsUser struct {
	Login string `json:"Login"`
	jwt.StandardClaims
}

func DecodeJWT(headertoken string) (Claims *ClaimsUser, err error) {
	token, err := jwt.ParseWithClaims(headertoken, &ClaimsUser{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(keyJWT), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsUser)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}

func EncodeJWT(login string) (token string, err error) {
	userClaims := ClaimsUser{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 9999).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Login: login,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	token, err = t.SignedString(keyJWT)
	if err != nil {
		return "", err
	}
	return token, nil
}
