package main

import (
	"encoding/json"
	"github.com/benbunsford/chirpy/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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
		UserID: params.UserID,
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
