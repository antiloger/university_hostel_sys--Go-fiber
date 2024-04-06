package util

import (
	"errors"

	"github.com/antiloger/nhostel-go/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Jwt_Secret) // Assuming your secret key is stored in the config

// ValidateTokenAndExtractClaims takes a JWT token as input and returns the user ID, role, and error
func ValidateTokenAndExtractClaims(tokenString string) (userID string, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok = claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("user_id claim missing")
	}

	role, ok = claims["role"].(string)
	if !ok {
		return "", "", errors.New("role claim missing")
	}

	return userID, role, nil
}
