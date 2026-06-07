package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"to-do-list/internal/service"
)

type TodosHandler struct {
	service *service.TodosService
}

func NewTodosHandler(service *service.TodosService) *TodosHandler {
	return &TodosHandler{service: service}
}

func (h *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	switch {
	case path == "/todos":
		h.handleCollection(w, r)
	case strings.HasPrefix(path, "/todos/"):
		h.handleItem(w, r)
	default:
		h.respondError(w, http.StatusNotFound, "not found")
	}
}

func (h *TodosHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TodosHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	id, err := h.extractID(r.URL.Path)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTodo(w, r, id)
	case http.MethodPut, http.MethodPatch:
		h.updateTodo(w, r, id)
	case http.MethodDelete:
		h.deleteTodo(w, r, id)
	default:
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TodosHandler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.ListTodos()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, todos)
}

func (h *TodosHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(payload.Title) == "" {
		h.respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := h.service.CreateTodo(payload.Title)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusCreated, todo)
}

func (h *TodosHandler) getTodo(w http.ResponseWriter, r *http.Request, id int64) {
	todo, err := h.service.GetTodo(id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if todo == nil {
		h.respondError(w, http.StatusNotFound, "todo not found")
		return
	}
	h.respondJSON(w, http.StatusOK, todo)
}

func (h *TodosHandler) updateTodo(w http.ResponseWriter, r *http.Request, id int64) {
	var payload struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing, err := h.service.GetTodo(id)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing == nil {
		h.respondError(w, http.StatusNotFound, "todo not found")
		return
	}

	if payload.Title != nil {
		existing.Title = strings.TrimSpace(*payload.Title)
	}
	if payload.Completed != nil {
		existing.Completed = *payload.Completed
	}

	if err := h.service.UpdateTodo(existing); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondJSON(w, http.StatusOK, existing)
}

func (h *TodosHandler) deleteTodo(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.service.DeleteTodo(id); err != nil {
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TodosHandler) extractID(path string) (int64, error) {
	trimmed := strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/todos/")
	return strconv.ParseInt(trimmed, 10, 64)
}

func (h *TodosHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TodosHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

func (h *TodosHandler) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
