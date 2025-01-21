package usecase

import (
	postgresrepo "golang-gorm/app/repository/postgres"
	s3repo "golang-gorm/app/repository/s3"
	"time"

	"github.com/go-playground/validator/v10"
)

type UsecaseDependency struct {
	Validate       *validator.Validate
	Timeout        time.Duration
	S3Repository   s3repo.S3Repo
	UserRepository postgresrepo.UserRepository
	TodoRepository postgresrepo.TodoRepository
	FileRepository postgresrepo.FileRepository
}
