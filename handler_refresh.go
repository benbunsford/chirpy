package main

import (
	"github.com/benbunsford/chirpy/internal/auth"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or invalid refresh token in request header.", err)
		return
	}

	dbUser, err := cfg.db.GetUserByRefreshToken(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error: Missing or invalid refresh token in request header.", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error creating access token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Server Error: Error creating token.", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	resp := response{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, resp)

}
