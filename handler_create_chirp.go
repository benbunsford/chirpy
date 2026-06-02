package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/benbunsford/chirpy/internal/auth"
	"github.com/benbunsford/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retrieving token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error: Error retrieving token.", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error: Error validating token.", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error decoding parameters.", err)
		return
	}

	validChirp, err := validateChirp(&params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}

	dbChirpParams := database.CreateChirpParams{
		Body:   validChirp,
		UserID: id,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), dbChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error adding chirp to database.", err)
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
