package auth

import (
	"github.com/alexedwards/argon2id"
)

func CheckPasswordHash(pw, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(pw, hash)
	if err != nil {
		return false, err
	}

	return result, nil
}
