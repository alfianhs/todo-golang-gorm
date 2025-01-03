package usecase

import (
	postgresrepo "golang-gorm/app/repository/postgres"
	"time"

	"github.com/go-playground/validator/v10"
)

type UsecaseDependency struct {
	Validate       *validator.Validate
	Timeout        time.Duration
	UserRepository postgresrepo.UserRepository
	TodoRepository postgresrepo.TodoRepository
}
