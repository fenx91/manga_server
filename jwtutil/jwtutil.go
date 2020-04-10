package jwtutil

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

var Secret []byte

type Claims struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	jwt.StandardClaims
}

func VerifyTokenAndGetUsername(r *http.Request) (UserEmail string, Nickname string, err error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", "", err
	}

	tokenString := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})

	if err != nil {
		return "", "", err
	}
	if token.Valid == false {
		return "", "", errors.New("token not valid")
	}

	return claims.Email, claims.Nickname, nil
}

func CreateJwt(email, nickname string) (string, error) {
	claims := &Claims{
		Email:          email,
		Nickname:       nickname,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(Secret)
	return signedToken, err
}
