package util

import (
	"errors"
	"strconv"

	"github.com/antiloger/nhostel-go/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Jwt_Secret)

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

	// Handling `id` which might be unmarshalled as float64
	if idFloat, ok := claims["id"].(float64); ok {
		userID = strconv.Itoa(int(idFloat)) // Convert float64 to int to string
	} else {
		return "", "", errors.New("user_id claim missing or not a number")
	}

	// Handling `role` which is expected to be a string
	if roleVal, ok := claims["role"].(string); ok {
		role = roleVal
	} else {
		return "", "", errors.New("role claim missing or not a string")
	}

	return userID, role, nil
}
