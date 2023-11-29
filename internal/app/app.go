package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/handlers"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/repository/memory"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/service"
)

func TaskManagerRun() {
	router := chi.NewRouter()
	logger, loggerErr := zap.NewProduction()
	if loggerErr != nil {
		logger.Fatal("Failed to create zap logger", zap.Error(loggerErr))
	}

	taskRepository := memory.NewTaskStorage(logger)
	taskService := service.NewTaskService(taskRepository, logger)
	handlers.NewTaskHandlers(taskService, logger, router)

	logger.Info("Starting task manager")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal("Failed to start task manager", zap.Error(err))
	}
}
