package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	pw1 := "123abc!@#"
	pw2 := "9lkjdlf289374(*&^)"
	hash1, _ := HashPassword(pw1)
	hash2, _ := HashPassword(pw2)

	tests := map[string]struct {
		pw            string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		"correct":                         {pw: pw1, hash: hash1, wantErr: false, matchPassword: true},
		"incorrect":                       {pw: "wrongpw", hash: hash1, wantErr: false, matchPassword: false},
		"pw doesn't match different hash": {pw: pw1, hash: hash2, wantErr: false, matchPassword: false},
		"empty pw":                        {pw: "", hash: hash1, wantErr: false, matchPassword: false},
		"invalid hash":                    {pw: pw1, hash: "invalidhash", wantErr: true, matchPassword: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.pw, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && match != tc.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tc.matchPassword, match)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	id := uuid.New()
	rightSecret := "123"
	wrongSecret := "1234"

	tests := map[string]struct {
		id               uuid.UUID
		signingSecret    string
		validatingSecret string
		expiresIn        time.Duration
		wantErr          bool
	}{
		"success":      {id: id, signingSecret: rightSecret, validatingSecret: rightSecret, expiresIn: 1 * time.Second, wantErr: false},
		"expired":      {id: id, signingSecret: rightSecret, validatingSecret: rightSecret, expiresIn: -1 * time.Second, wantErr: true},
		"wrong secret": {id: id, signingSecret: rightSecret, validatingSecret: wrongSecret, expiresIn: 1 * time.Second, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			token, _ := MakeJWT(tc.id, tc.signingSecret, tc.expiresIn)

			validatedID, err := ValidateJWT(token, tc.validatingSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr && validatedID != uuid.Nil {
				t.Errorf("ValidateJWT() should return a nil uuid on error, instead returned %v", validatedID)
			}
			if !tc.wantErr && tc.id != validatedID {
				t.Errorf("ValidateJWT() expected to return %v, instead returned %v", tc.id, validatedID)
			}
		})
	}
}
