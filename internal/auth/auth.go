package auth

import (
	"github.com/alexedwards/argon2id"
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
