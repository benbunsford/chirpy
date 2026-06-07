package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benbunsford/chirpy/internal/auth"
	"github.com/benbunsford/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string
		Password string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error decoding parameters.", err)
		return
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error creating access token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error creating token.", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	dbRefreshParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(1440 * time.Hour), //60 days
	}

	err = cfg.db.CreateRefreshToken(r.Context(), dbRefreshParams)
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error creating refresh token.", err)
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	resp := response{
		User:         user,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
