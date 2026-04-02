package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFilterProfanity(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":             {input: "There was a kerfuffle today.", want: "There was a **** today."},
		"attached character": {input: "Oh, shartbert!", want: "Oh, shartbert!"},
		"capitalization":     {input: "I AM FORNAX", want: "I AM ****"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := filterProfanity(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
