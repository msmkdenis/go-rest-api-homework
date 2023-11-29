package dto

type TaskRequest struct {
	ID           string   `json:"id" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Note         string   `json:"note" validate:"required"`
	Applications []string `json:"applications" validate:"required"`
}

type TaskResponse struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}
