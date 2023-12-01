package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/handlers/dto"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/mapper"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/repository/memory"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/service"
)

type TaskHandlersSuite struct {
	suite.Suite
	h *TaskHandlers
	r *memory.TaskStorage
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlersSuite))
}

func (s *TaskHandlersSuite) SetupTest() {
	router := chi.NewRouter()
	logger, loggerErr := zap.NewProduction()
	if loggerErr != nil {
		logger.Fatal("Failed to create zap logger", zap.Error(loggerErr))
	}
	s.r = memory.NewTaskStorage(logger)
	taskService := service.NewTaskService(s.r, logger)
	s.h = NewTaskHandlers(taskService, logger, router)
}

func (s *TaskHandlersSuite) TearDownTest() {
	clear(s.r.Tasks)
}

func (s *TaskHandlersSuite) TestAddTask() {
	successTask := model.Task{
		ID:           "3",
		Description:  "Description",
		Note:         "Note",
		Applications: []string{"app1", "app2"},
	}

	successTaskDto := mapper.ToTaskResponse(successTask)
	successTaskJSON, err := json.Marshal(successTaskDto)
	require.NoError(s.T(), err)

	invalidTask := model.Task{
		ID:           "3",
		Note:         "Note",
		Applications: []string{"app1", "app2"},
	}

	invalidTaskJSON, err := json.Marshal(invalidTask)
	require.NoError(s.T(), err)

	testCases := []struct {
		name         string
		method       string
		header       http.Header
		body         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "BadRequest - empty body",
			method:       http.MethodPost,
			header:       map[string][]string{"Content-Type": {""}},
			body:         "",
			path:         "http://localhost:8080/tasks",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Wrong content type\n",
		},
		{
			name:         "Success",
			method:       http.MethodPost,
			body:         string(successTaskJSON),
			header:       map[string][]string{"Content-Type": {"application/json"}},
			path:         "http://localhost:8080/tasks",
			expectedCode: http.StatusCreated,
			expectedBody: string(successTaskJSON),
		},
		{
			name:         "Bad request - invalid task",
			method:       http.MethodPost,
			body:         string(invalidTaskJSON),
			header:       map[string][]string{"Content-Type": {"application/json"}},
			path:         "http://localhost:8080/tasks",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request\n",
		},
		{
			name:         "Bad request - unable to unmarshall",
			method:       http.MethodPost,
			body:         "unable to unmarshall",
			header:       map[string][]string{"Content-Type": {"application/json"}},
			path:         "http://localhost:8080/tasks",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Unable to unmarshal request\n",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.header.Get("Content-Type"))
			w := httptest.NewRecorder()
			s.h.router.ServeHTTP(w, request)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func (s *TaskHandlersSuite) TestGetByID() {
	successTask := model.Task{
		ID:           "3",
		Description:  "Description",
		Note:         "Note",
		Applications: []string{"app1", "app2"},
	}
	s.r.Tasks[successTask.ID] = successTask
	successTaskDto := mapper.ToTaskResponse(successTask)
	successTaskJSON, err := json.Marshal(successTaskDto)
	require.NoError(s.T(), err)

	testCases := []struct {
		name         string
		method       string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodGet,
			path:         "http://localhost:8080/tasks/3",
			expectedCode: http.StatusOK,
			expectedBody: string(successTaskJSON),
		},
		{
			name:         "Bad request - unable to find task",
			method:       http.MethodGet,
			path:         "http://localhost:8080/tasks/4",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Task with id 4 not found\n",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			s.h.router.ServeHTTP(w, request)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func (s *TaskHandlersSuite) TestGetAll() {
	tasks := make([]dto.TaskResponse, 0, len(s.r.Tasks))
	for _, v := range s.r.Tasks {
		tasks = append(tasks, mapper.ToTaskResponse(v))
	}

	testCases := []struct {
		name         string
		method       string
		path         string
		expectedCode int
		expectedBody []dto.TaskResponse
	}{
		{
			name:         "Success",
			method:       http.MethodGet,
			path:         "http://localhost:8080/tasks",
			expectedCode: http.StatusOK,
			expectedBody: tasks,
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			fmt.Println(s.r.Tasks)
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			s.h.router.ServeHTTP(w, request)

			assert.Equal(t, test.expectedCode, w.Code)

			var answer []dto.TaskResponse
			err := json.Unmarshal(w.Body.Bytes(), &answer)
			require.NoError(t, err)
			assert.Equal(t, len(test.expectedBody), len(answer))
		})
	}
}

func (s *TaskHandlersSuite) TestDeleteByID() {
	testCases := []struct {
		name         string
		method       string
		path         string
		elements     int
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success",
			method:       http.MethodDelete,
			path:         "http://localhost:8080/tasks/1",
			elements:     1,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Bad request - unable to find task",
			method:       http.MethodDelete,
			path:         "http://localhost:8080/tasks/4",
			elements:     1,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Task with id 4 not found\n",
		},
	}

	for _, test := range testCases {
		s.T().Run(test.name, func(t *testing.T) {
			fmt.Println(s.r.Tasks)
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			s.h.router.ServeHTTP(w, request)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
			assert.Equal(t, test.elements, len(s.r.Tasks))
		})
	}
}
