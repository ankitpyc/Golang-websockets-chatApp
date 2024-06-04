package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type TokenServiceInf interface {
	GetToken(secret string) string
	ValidateToken(secret string, tok string) (bool, error)
}

func ValidateToken(secret string, tok string) (bool, error) {
	claims := &Claims{}
	log.Printf("Validating token")
	jwtToken, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return false, err
	}

	if !jwtToken.Valid {
		return false, errors.New("invalid token")
	}

	return true, nil
}

func GetToken(secret string) string {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 120).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("error :  ", err)
	}
	return tok
}
