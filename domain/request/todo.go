package request

import "golang-gorm/domain/model"

type CreateTodoRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateTodoRequest struct {
	Name   string           `json:"name" validate:"required"`
	Status model.TodoStatus `json:"status" validate:"required,todo_status"`
}
