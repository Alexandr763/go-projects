package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"gamification/internal/transport/http/apigen"
)

const maxBody = 1 << 20

func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "invalid JSON")
		return false
	}
	var extra struct{}
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, apigen.ErrorBodyErrorCodeValidationError, "invalid JSON")
		return false
	}
	return true
}

func validEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" || strings.Contains(email, " ") || !strings.Contains(email, "@") {
		return false
	}
	return true
}

func validPage(limit, offset *int) bool {
	if limit != nil && (*limit < 1 || *limit > 100) {
		return false
	}
	if offset != nil && *offset < 0 {
		return false
	}
	return true
}
