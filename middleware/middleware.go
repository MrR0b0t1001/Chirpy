package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/MrR0b0t1001/Chirpy/config"
	"github.com/MrR0b0t1001/Chirpy/internal/auth"
	"github.com/MrR0b0t1001/Chirpy/utils"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			utils.PermissionDenied(w)
			return
		}

		userID, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			utils.PermissionDenied(w)
			return
		}

		_, err = cfg.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.PermissionDenied(w)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		handlerFunc(w, r.WithContext(ctx))
	}
}

func WithAPIAuth(handlerFunc http.HandlerFunc, cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIToken(r.Header)
		if err != nil {
			utils.PermissionDenied(w)
			return
		}

		if apiKey != cfg.APIKey {
			utils.PermissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}
