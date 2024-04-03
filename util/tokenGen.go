package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, name string, user_type string) (signedToken, signedRefreshToken string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"name":      name,
		"user_type": user_type,
		"exp":       time.Now().Add(time.Hour * time.Duration(24)).Unix(),
	})

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     email,
		"name":      name,
		"user_type": user_type,
		"exp":       time.Now().Local().Add(time.Hour * time.Duration(169)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshString, err := refresh.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshString, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	const op = "util.ValidateToken"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Header["alg"])
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return claims, nil
}
