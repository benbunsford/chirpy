package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ben-rw/chirpy/internal/auth"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerUpgradeUserToChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Malformed request header", err)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "API Key mismatch", err)
		return
	}

	type parameters struct {
		Event string
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		}
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request body: %s", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user: %s", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
