package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/galiullindo/todolist-rest-api/internal/model"
)

type Service interface {
	GetAll() ([]model.Task, error)
	GetByID(id int64) (model.Task, error)
	Create(title string, status string) (model.Task, error)
}

type handler struct {
	svc Service
}

func New(svc Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.GetAll)
	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	mux.HandleFunc("POST /tasks", h.Create)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll()
	if err != nil {
		log.Println("handler:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("handler:", err)
	}
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("handler:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.svc.GetByID(id)
	if err != nil {
		log.Println("handler:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Println("handler:", err)
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Status string
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("handler:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.svc.Create(input.Title, input.Status)
	if err != nil {
		log.Println("handler:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Println("handler:", err)
	}
}
