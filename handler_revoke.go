package main

import (
	"github.com/ben-rw/chirpy/internal/auth"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or invalid refresh token in request header.", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), headerToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error revoking refresh token.", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
