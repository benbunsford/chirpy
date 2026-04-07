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
		respondWithError(w, 500, "Server Error: Error decoding parameters.")
		return
	}

	validChirp, err := validateChirp(&params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	dbChirpParams := database.CreateChirpParams{
		Body:   validChirp,
		UserID: params.UserID,
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), dbChirpParams)
	if err != nil {
		respondWithError(w, 500, "Server Error: Error adding chirp to database.")
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		User_ID:   dbChirp.UserID,
	}

	respondWithJSON(w, 201, chirp)
}
