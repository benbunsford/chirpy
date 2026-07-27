package main

import (
	"github.com/ben-rw/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"sort"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error
	if authorId != "" {
		authorUuid, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error: Malformed author uuid.", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUuid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Server Error: Error retrieving chirps from database.", err)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Server Error: Error retrieving chirps from database.", err)
			return
		}
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}

		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sorting == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
