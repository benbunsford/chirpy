package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ben-rw/chirpy/internal/auth"
	"github.com/ben-rw/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string
		Password string
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or malformed access token in request header.", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or malformed access token in request header.", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server error: Error decoding request body.", err)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server error: Error hashing password.", err)
		return
	}

	updateUserParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
		ID:             id,
	}

	dbUser, err := cfg.db.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server error: Error updating user.", err)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
}
