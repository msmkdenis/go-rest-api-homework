package mapper

import (
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/handlers/dto"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
)

func ToTaskResponse(task model.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:           task.ID,
		Description:  task.Description,
		Note:         task.Note,
		Applications: task.Applications,
	}
}

func ToTaskModel(task dto.TaskRequest) model.Task {
	return model.Task{
		ID:           task.ID,
		Description:  task.Description,
		Note:         task.Note,
		Applications: task.Applications,
	}
}
