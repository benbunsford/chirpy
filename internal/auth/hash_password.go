package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(pw string) (string, error) {
	hashed_pw, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	if err != nil {
		return "", nil
	}

	return hashed_pw, nil
}
