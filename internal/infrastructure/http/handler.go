package http

import (
	"encoding/json"
	"net/http"

	"github.com/attaleast/idp-task-manager-demo/internal/application"
	"github.com/attaleast/idp-task-manager-demo/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	uc *application.TaskUseCase
}

func NewTaskHandler(uc *application.TaskUseCase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) Register(r chi.Router) {
	r.Post("/api/v1/tasks", h.CreateTask)
	r.Get("/api/v1/tasks", h.ListTasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID   string `json:"project_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	projID, _ := uuid.Parse(req.ProjectID)
	task, err := h.uc.CreateTask(r.Context(), application.CreateTaskInput{
		ProjectID:   projID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "invalid or missing project_id", http.StatusBadRequest)
		return
	}

	tasks, err := h.uc.ListTasks(r.Context(), projID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
