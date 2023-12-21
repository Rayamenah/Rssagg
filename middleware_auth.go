package main

import (
	"fmt"
	"net/http"

	"github.com/Rayamenah/rssagg/internal/auth"
	"github.com/Rayamenah/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			writeErr(w, 403, fmt.Sprintf("auth error: %s", err))
			return
		}
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			writeErr(w, 400, fmt.Sprintf("couldnt get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
