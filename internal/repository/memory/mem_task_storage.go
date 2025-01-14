package memory

import (
	"sync"

	"go.uber.org/zap"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/apperrors"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/utils"
)

type TaskStorage struct {
	mu     sync.RWMutex
	logger *zap.Logger
	Tasks  map[string]model.Task
}

func NewTaskStorage(logger *zap.Logger) *TaskStorage {
	return &TaskStorage{
		mu:     sync.RWMutex{},
		logger: logger,
		Tasks: map[string]model.Task{
			"1": {
				ID:          "1",
				Description: "Сделать финальное задание темы REST API",
				Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
				Applications: []string{
					"VS Code",
					"Terminal",
					"git",
				},
			},
			"2": {
				ID:          "2",
				Description: "Протестировать финальное задание с помощью Postmen",
				Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
				Applications: []string{
					"VS Code",
					"Terminal",
					"git",
					"Postman",
				},
			},
		},
	}
}

func (r *TaskStorage) Insert(task model.Task) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	savedTask := task
	r.Tasks[savedTask.ID] = savedTask

	return &savedTask, nil
}

func (r *TaskStorage) SelectByID(id string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.Tasks[id]
	if !ok {
		return nil, apperrors.NewValueError(utils.Caller(), apperrors.ErrTaskNotFound)
	}

	return &task, nil
}

func (r *TaskStorage) SelectAll() ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]model.Task, 0, len(r.Tasks))
	for _, task := range r.Tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskStorage) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.Tasks[id]
	if !ok {
		return apperrors.NewValueError(utils.Caller(), apperrors.ErrTaskNotFound)
	}

	delete(r.Tasks, id)

	return nil
}
