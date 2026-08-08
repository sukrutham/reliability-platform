package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sukrutham/reliability-platform/internal/store"
)

type TodoHandler struct {
	Store *store.TodoStore
}

func NewTodoHandler(s *store.TodoStore) *TodoHandler {
	return &TodoHandler{Store: s}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos := h.Store.List()
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	todo := h.Store.Create(req.Title)
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.Store.Get(id)
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
