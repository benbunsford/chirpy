package main

import "errors"

func validateChirp(msg *string) (string, error) {
	if len(*msg) > 140 {
		return "", errors.New("Chirp is too long - must be 140 characters or less.")
	}

	return filterProfanity(*msg), nil
}
