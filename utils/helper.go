package utils

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func ValidateChirp(body string) bool {
	if len(body) > 140 {
		return false
	}
	return true
}

func ExtractID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("chirpID")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("No token provided")
	}

	chirpID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid Format")
	}

	return chirpID, nil
}
