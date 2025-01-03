package config

import (
	"golang-gorm/app/delivery/http/middleware"
	http_user "golang-gorm/app/delivery/http/user"
	postgresrepo "golang-gorm/app/repository/postgres"
	"golang-gorm/app/usecase"
	usecase_user "golang-gorm/app/usecase/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	GinEngine *gin.Engine
	DB        *gorm.DB
	Validator *validator.Validate
	Timeout   time.Duration
}

func Bootstrap(config BootstrapConfig) {
	// init repository
	userRepository := postgresrepo.NewUserRepository(config.DB)
	todoRepository := postgresrepo.NewTodoRepository(config.DB)

	// init usecase
	userAuthUsecase := usecase_user.NewAuthUsecase(usecase.UsecaseDependency{
		UserRepository: userRepository,
		Validate:       config.Validator,
		Timeout:        config.Timeout,
	})
	userTodoUsecase := usecase_user.NewTodoUsecase(usecase.UsecaseDependency{
		TodoRepository: todoRepository,
		Validate:       config.Validator,
		Timeout:        config.Timeout,
	})

	// init auth middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// init http delivery
	http_user.NewAuthHandler(config.GinEngine, authMiddleware, userAuthUsecase)
	http_user.NewTodoHandler(config.GinEngine, authMiddleware, userTodoUsecase)
}
