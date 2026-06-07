package main

import (
	"github.com/benbunsford/chirpy/internal/auth"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or malformed access token.", err)
		return
	}

	RequestUserId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Error: Unable to validate token.", err)
		return
	}

	chirpIdString := r.PathValue("ID")
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The provided path parameter is not a valid uuid.", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error: Chirp not found.", err)
		return
	}

	if chirp.UserID != RequestUserId {
		respondWithError(w, http.StatusForbidden, "Error: Only the owner of a chirp can delete it.", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error deleting chirp.", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
