package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Payload struct {
	Name  string
	Phone string
	Role  string
	Id    uint
}

type Claims struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	Id    uint   `json:"id"`
	jwt.StandardClaims
}

var JWT_SECRET string

func GenerateJwtToken(payload Payload) (string, error) {
	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)

	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)

	claims := &Claims{
		Id:    payload.Id,
		Name:  payload.Name,
		Phone: payload.Phone,
		Role:  payload.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	UnsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := UnsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}

func VerifyJwtToken(strToken string) (*Claims, error) {
	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, fmt.Errorf("invalid token signature")
		}
	}

	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}

	return claims, nil
}
