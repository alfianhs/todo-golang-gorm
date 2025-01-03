package config

import (
	"golang-gorm/domain/model"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("todo_status", _todoStatusValidator)
	return validate
}

func _todoStatusValidator(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch model.TodoStatus(status) {
	case model.TodoStatusDone, model.TodoStatusNotStarted:
		return true
	default:
		return false
	}
}
