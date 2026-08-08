package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sukrutham/reliability-platform/internal/store"
)

func newTestHandler() *TodoHandler {
	return NewTodoHandler(store.NewTodoStore())
}

func TestCreateAndGetTodo(t *testing.T) {
	h := newTestHandler()

	body := bytes.NewBufferString(`{"title":"write tests"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", body)
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "write tests") {
		t.Errorf("expected response to contain title, got %q", rec.Body.String())
	}
}

func TestCreateTodoMissingTitle(t *testing.T) {
	h := newTestHandler()

	body := bytes.NewBufferString(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", body)
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestListTodos(t *testing.T) {
	h := newTestHandler()
	h.Store.Create("first")
	h.Store.Create("second")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	rec := httptest.NewRecorder()

	h.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "first") || !strings.Contains(rec.Body.String(), "second") {
		t.Errorf("expected both todos in response, got %q", rec.Body.String())
	}
}

func TestGetTodoNotFound(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos/999", nil)
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()

	h.Get(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
