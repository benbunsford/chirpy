package main

import (
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDstring := r.PathValue("ID")
	chirpID, err := uuid.Parse(chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The provided path parameter is not a valid uuid.", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "The provided uuid was not found in the database.", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
