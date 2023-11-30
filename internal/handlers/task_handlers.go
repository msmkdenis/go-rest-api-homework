package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/apperrors"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/handlers/dto"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/mapper"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/utils"
)

type TaskService interface {
	Add(task model.Task) (*model.Task, error)
	GetByID(id string) (*model.Task, error)
	GetAll() ([]model.Task, error)
	DeleteByID(id string) error
}

type TaskHandlers struct {
	taskService TaskService
	logger      *zap.Logger
	router      *chi.Mux
}

func NewTaskHandlers(taskService TaskService, logger *zap.Logger, router *chi.Mux) *TaskHandlers {
	h := &TaskHandlers{
		taskService: taskService,
		logger:      logger,
		router:      router,
	}

	h.router.Post("/tasks", h.Add)
	h.router.Get("/tasks/{id}", h.GetByID)
	h.router.Get("/tasks", h.GetAll)
	h.router.Delete("/tasks/{id}", h.DeleteByID)

	return h
}

func (h *TaskHandlers) Add(w http.ResponseWriter, r *http.Request) {
	contentHeader := r.Header.Get("Content-Type")
	if contentHeader != "application/json" {
		h.logger.Warn("Bad Request: wrong content type", zap.Error(fmt.Errorf("caller: %s", utils.Caller())))
		http.Error(w, "Wrong content type", http.StatusBadRequest)
		return
	}

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		h.logger.Warn("Bad Request: unknown error", zap.Error(readBodyErr))
		http.Error(w, "Unable to read request", http.StatusBadRequest)
		return
	}

	var task dto.TaskRequest
	unmarshalErr := json.Unmarshal(body, &task)
	if unmarshalErr != nil {
		h.logger.Warn("Bad Request: unknown error", zap.Error(unmarshalErr))
		http.Error(w, "Unable to unmarshal request", http.StatusBadRequest)
		return
	}

	taskValidator := validator.New()
	if err := taskValidator.Struct(task); err != nil {
		h.logger.Warn("Bad Request: invalid request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	savedTask, err := h.taskService.Add(mapper.ToTaskModel(task))
	if err != nil {
		h.logger.Warn("Bad Request: unknown error", zap.Error(err))
		http.Error(w, "Unable to add task", http.StatusBadRequest)
		return
	}

	taskResponseDto := mapper.ToTaskResponse(*savedTask)
	response, marshalErr := json.Marshal(taskResponseDto)
	if marshalErr != nil {
		h.logger.Error("Bad Request: unknown error", zap.Error(marshalErr))
		http.Error(w, "Unable to marshal response", http.StatusBadRequest)
		return
	}

	_, writeErr := w.Write(response)
	if writeErr != nil {
		h.logger.Error("Bad Request: unknown error", zap.Error(writeErr))
		http.Error(w, "Unable to write response", http.StatusBadRequest)
		return
	}
}

func (h *TaskHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.taskService.GetByID(id)
	if err != nil && !errors.Is(err, apperrors.ErrTaskNotFound) {
		h.logger.Error("Bad Request: unknown error", zap.Error(err))
		http.Error(w, "Unable to get task", http.StatusBadRequest)
		return
	}

	if errors.Is(err, apperrors.ErrTaskNotFound) {
		h.logger.Warn("Bad Request: task not found", zap.String("stacktrace", err.Error()))
		http.Error(w, fmt.Sprintf("Task with id %s not found", id), http.StatusBadRequest)
		return
	}

	taskResponseDto := mapper.ToTaskResponse(*task)
	response, marshalErr := json.Marshal(taskResponseDto)
	if marshalErr != nil {
		h.logger.Error("Bad Request: unknown error", zap.Error(marshalErr))
		http.Error(w, "Unable to marshal response", http.StatusBadRequest)
		return
	}

	_, writeErr := w.Write(response)
	if writeErr != nil {
		h.logger.Error("Bad Request: unknown error", zap.Error(writeErr))
		http.Error(w, "Unable to write response", http.StatusBadRequest)
		return
	}
}

func (h *TaskHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAll()
	if err != nil {
		h.logger.Error("Internal Server Error: unknown error", zap.Error(err))
		http.Error(w, "Unable to get tasks", http.StatusInternalServerError)
		return
	}

	tasksDto := make([]dto.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		tasksDto = append(tasksDto, mapper.ToTaskResponse(task))
	}

	response, marshalErr := json.Marshal(tasksDto)
	if marshalErr != nil {
		h.logger.Error("Internal Server Error: unknown error", zap.Error(marshalErr))
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	_, writeErr := w.Write(response)
	if writeErr != nil {
		h.logger.Error("Internal Server Error: unknown error", zap.Error(writeErr))
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandlers) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.taskService.DeleteByID(id)
	if err != nil && !errors.Is(err, apperrors.ErrTaskNotFound) {
		h.logger.Error("Bad Request: unknown error", zap.Error(err))
		http.Error(w, "Unable to get task", http.StatusBadRequest)
		return
	}

	if errors.Is(err, apperrors.ErrTaskNotFound) {
		h.logger.Error("Bad Request: task not found", zap.Error(err))
		http.Error(w, fmt.Sprintf("Task with id %s not found", id), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
