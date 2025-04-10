package config

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/MrR0b0t1001/Chirpy/internal/auth"
	"github.com/MrR0b0t1001/Chirpy/internal/database"
	"github.com/MrR0b0t1001/Chirpy/pkg/types"
	"github.com/MrR0b0t1001/Chirpy/utils"
	"github.com/google/uuid"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	JWTSecret      string
}

func (cfg *APIConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) error {
	hits := cfg.FileserverHits.Load()
	t, err := template.ParseFiles("admin/admin.html")
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Error Occurred"})
	}

	t.Execute(w, hits)
	return nil
}

func (cfg *APIConfig) ResetHandler(w http.ResponseWriter, r *http.Request) error {
	cfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	return nil
}

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (cfg *APIConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	req := types.CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	return utils.WriteJSON(w, http.StatusCreated, types.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *APIConfig) HandleDeleteUsers(w http.ResponseWriter, r *http.Request) error {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "Forbidden!!!"})
	}

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusOK, utils.ApiDelResponse{Message: "Reset Successful"})
}

func (cfg *APIConfig) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	req := types.LoginUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		return utils.WriteJSON(w, http.StatusNotFound, err)
	}

	if err := auth.CheckHashPassword(user.HashedPassword, req.Password); err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: "Unable to create JWT Token"},
		)
	}

	rToken, err := auth.MakeRefreshToken()
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			utils.ApiError{Error: "Unable to create Refresh Token"},
		)
	}

	if _, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     rToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		UserID:    user.ID,
	}); err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusOK, types.LoginReponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: rToken,
	})
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (cfg *APIConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) error {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusUnauthorized,
			utils.ApiError{Error: "Provided token is not valid"},
		)
	}

	id, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusUnauthorized,
			utils.ApiError{Error: "JWT token invalid"},
		)
	}

	req := types.CreateChirpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	if ok := utils.ValidateChirp(req.Body); !ok {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: "Length of body exceeds limit of 140 characters"},
		)
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      req.Body,
		UserID:    id,
	})
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	return utils.WriteJSON(w, http.StatusCreated, types.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *APIConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) error {
	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	chirpsArray := []types.Chirp{}
	for _, chirp := range chirps {
		chirpsArray = append(chirpsArray, types.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	return utils.WriteJSON(w, http.StatusOK, chirpsArray)
}

func (cfg *APIConfig) HandleGetChirpByID(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.ExtractID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusNotFound, err)
	}

	chirp, err := cfg.DB.GetChirpByID(r.Context(), id)
	if err != nil {
		return utils.WriteJSON(w, http.StatusNotFound, err)
	}

	return utils.WriteJSON(w, http.StatusOK, types.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *APIConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) error {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	refreshToken, err := cfg.DB.GetUserFromRefreshToken(r.Context(), authToken)
	if err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return utils.WriteJSON(w, http.StatusUnauthorized, utils.ApiError{Error: "Expired token"})
	}

	if refreshToken.RevokedAt.Valid {
		return utils.WriteJSON(w, http.StatusUnauthorized, utils.ApiError{Error: "Token Revoked"})
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.JWTSecret)
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusOK, types.RefreshTokenResponse{
		Token: token,
	})
}

func (cfg *APIConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) error {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	if err := cfg.DB.RevokeUserToken(r.Context(), authToken); err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (cfg *APIConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	userID, err := auth.ValidateJWT(authToken, cfg.JWTSecret)
	if err != nil {
		return utils.WriteJSON(w, http.StatusUnauthorized, err)
	}

	req := types.UpdateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	if err := cfg.DB.UpdateUserCreds(r.Context(), database.UpdateUserCredsParams{
		ID:             userID,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}); err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusOK, types.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
