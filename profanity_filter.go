package main

import (
	"strings"
)

var profanity map[string]struct{} = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func filterProfanity(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		if _, ok := profanity[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
