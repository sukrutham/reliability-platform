package handlers

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeStatus(w, "ok")
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	writeStatus(w, "ready")
}

func writeStatus(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResponse{Status: status})
}
