package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benbunsford/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string
		Password         string
		ExpiresInSeconds *int
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error decoding parameters.", err)
		return
	}

	expires := 3600 //default expiration, 1 hour
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds < 3600 {
		expires = *params.ExpiresInSeconds
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password.", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		log.Printf("Error validating password: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password.", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password.", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(expires)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error creating token.", err)
		return
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	resp := response{
		User:  user,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
