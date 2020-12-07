package utils

import (
	"SMS/model"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

//ValidateToken funtion to parse the token and compare it with the JWT_SECRET
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		godotenv.Load()
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unauthenticated")
	}

	return token, nil
}

//ValidateUser function
func ValidateUser(tokenString string) (model.User, bool) {
	var user model.User
	token, err := ValidateToken(tokenString)
	if err != nil {
		return user, false
	}
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		user.ID = uint(claims["user_id"].(float64))
		return user, true
	}
	return user, false
}
