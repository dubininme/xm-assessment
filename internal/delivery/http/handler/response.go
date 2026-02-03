package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func writeErr(w http.ResponseWriter, status int, code oapi.ErrorCode, message string) {
	writeJSON(w, status, oapi.Error{Code: code, Message: message})
}
