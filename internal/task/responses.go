package task

import (
	"github.com/google/uuid"
)

type CreateResponse struct {
	Uuid uuid.UUID `json:"uuid"`
}

type GetAllResponse struct {
	Tasks []*Task `json:"tasks"`
}

type GetByUuidResponse struct {
	Task *Task `json:"task"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type UpdateResponse struct {
	Status Status `json:"status"`
}
