package utils

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Generate a JWT Token
func GenerateJWT(userId string) (string, error) {
	JWT_SECRET := []byte(os.Getenv("JWT_SECRET"))

	// Define an expiration time for the token
	// expTime := time.Now().Add(24 * time.Hour)
	expTime := time.Now().Add(15 * time.Minute) // 15 minutes, will change it for production

	// Create the JWT claims that include the userId and the expiration time
	claims := &jwt.MapClaims{
		"userId": userId,
		"exp":    expTime.Unix(),
	}

	// Create a token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	tokenString, err := token.SignedString(JWT_SECRET)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT verifies the given JWT token and returns the claims
func VerifyJWT(tokenString string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))


	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}