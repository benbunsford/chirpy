package auth

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func HashPassword(pw string) (string, error) {
	hashed_pw, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashed_pw, nil
}

func CheckPasswordHash(pw, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(pw, hash)
	if err != nil {
		return false, err
	}

	return result, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   fmt.Sprintf("%v", userID),
	})

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	stringID := claims.Subject

	convertedID, err := uuid.Parse(stringID)
	if err != nil {
		return uuid.Nil, err
	}

	return convertedID, nil
}
