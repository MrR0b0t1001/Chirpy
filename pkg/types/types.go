package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

type CreateUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginUserRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds string `json:"expires_in_seconds"`
}

type LoginReponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

/////////////////////////////////////////////////////

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type CreateChirpRequest struct {
	Body string `json:"body"`
}

type GetChirpsReq struct {
	Email string `json:"email"`
}

/////////////////////////////////////////////////////

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

/////////////////////////////////////////////////////

type UpdateUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

/////////////////////////////////////////////////////

type PolkaUpdateRequest struct {
	Event string    `json:"event"`
	Data  PolkaData `json:"data"`
}

type PolkaData struct {
	UserID string `json:"user_id"`
}

/////////////////////////////////////////////////////
