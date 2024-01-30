package service

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/utils"
)

type TaskRepository interface {
	Insert(task model.Task) (*model.Task, error)
	SelectByID(id string) (*model.Task, error)
	SelectAll() ([]model.Task, error)
	DeleteByID(id string) error
}

type TaskService struct {
	taskRepository TaskRepository
	logger         *zap.Logger
}

func NewTaskService(taskRepository TaskRepository, logger *zap.Logger) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
		logger:         logger,
	}
}

func (s *TaskService) Add(task model.Task) (*model.Task, error) {
	savedTask, err := s.taskRepository.Insert(task)
	if err != nil {
		return nil, fmt.Errorf("%s %w", utils.Caller(), err)
	}

	return savedTask, nil
}

func (s *TaskService) GetByID(id string) (*model.Task, error) {
	task, err := s.taskRepository.SelectByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s %w", utils.Caller(), err)
	}

	return task, nil
}

func (s *TaskService) GetAll() ([]model.Task, error) {
	tasks, err := s.taskRepository.SelectAll()
	if err != nil {
		return nil, fmt.Errorf("%s %w", utils.Caller(), err)
	}

	return tasks, nil
}

func (s *TaskService) DeleteByID(id string) error {
	err := s.taskRepository.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("%s %w", utils.Caller(), err)
	}

	return nil
}
