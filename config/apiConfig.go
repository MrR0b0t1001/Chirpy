package config

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/MrR0b0t1001/Chirpy/internal/database"
	"github.com/MrR0b0t1001/Chirpy/utils"
	"github.com/google/uuid"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type CreateChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
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

func (cfg *APIConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	req := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     req.Email,
	})
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	return utils.WriteJSON(w, http.StatusCreated, User{
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

func (cfg *APIConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) error {
	req := CreateChirpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	if ok := validateChirp(req.Body); !ok {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			utils.ApiError{Error: "Length of body exceeds limit"},
		)
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      req.Body,
		UserID:    req.UserID,
	})
	if err != nil {
		return utils.WriteJSON(
			w,
			http.StatusBadRequest,
			err,
		)
	}

	return utils.WriteJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) bool {
	if len(body) > 140 {
		return false
	}
	return true
}
