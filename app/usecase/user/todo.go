package usecase_user

import (
	"context"
	postgresrepo "golang-gorm/app/repository/postgres"
	"golang-gorm/app/usecase"
	"golang-gorm/domain/model"
	"golang-gorm/domain/request"
	"golang-gorm/helpers"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type todoUsecase struct {
	todoRepository postgresrepo.TodoRepository
	contextTimeout time.Duration
	validate       *validator.Validate
}

func NewTodoUsecase(d usecase.UsecaseDependency) TodoUsecase {
	return &todoUsecase{
		todoRepository: d.TodoRepository,
		contextTimeout: d.Timeout,
		validate:       d.Validate,
	}
}

type TodoUsecase interface {
	GetAll(ctx context.Context, claim model.JWTClaimUser, query url.Values) helpers.PaginatedResponse
	GetOne(ctx context.Context, claim model.JWTClaimUser, todoID string) helpers.Response
	Create(ctx context.Context, claim model.JWTClaimUser, payload request.CreateTodoRequest) helpers.Response
	Update(ctx context.Context, claim model.JWTClaimUser, todoID string, payload request.UpdateTodoRequest) helpers.Response
	Delete(ctx context.Context, claim model.JWTClaimUser, todoID string) helpers.Response
}

func (u *todoUsecase) GetAll(ctx context.Context, claim model.JWTClaimUser, query url.Values) helpers.PaginatedResponse {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get offset & limit
	page, offset, limit := helpers.GetOffsetLimit(query)

	filters := map[string]interface{}{
		"user_id": claim.UserID,
	}

	// count first
	totalData, err := u.todoRepository.Count(ctx, filters)
	if err != nil {
		return helpers.PaginatedResponse{
			Status:  http.StatusInternalServerError,
			Message: "error count todo",
		}
	}

	if totalData == 0 {
		return helpers.PaginatedResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    []interface{}{},
			Meta: map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": totalData,
			},
		}
	}

	// fetch data
	todos, err := u.todoRepository.FetchList(ctx, offset, limit, filters)
	if err != nil {
		return helpers.PaginatedResponse{
			Status:  http.StatusInternalServerError,
			Message: "error fetch todo",
		}
	}

	return helpers.PaginatedResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    todos,
		Meta: map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": totalData,
		},
	}
}

func (u *todoUsecase) GetOne(ctx context.Context, claim model.JWTClaimUser, todoID string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check todo exist
	todo, err := u.todoRepository.FindOne(ctx, map[string]interface{}{
		"id":      todoID,
		"user_id": claim.UserID,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if todo == nil {
		return helpers.Response{
			Data:    nil,
			Message: "todo not found",
			Status:  http.StatusBadRequest,
		}
	}

	return helpers.Response{
		Data:    todo,
		Message: "success",
		Status:  http.StatusOK,
	}
}

func (u *todoUsecase) Create(ctx context.Context, claim model.JWTClaimUser, payload request.CreateTodoRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	logrus.Infof("payload: %v", payload.Name)

	// validate payload
	validationResponse, err := helpers.ValidateBody(u.validate, payload)
	if err != nil {
		return validationResponse
	}

	// create todo
	newTodo := model.Todo{
		ID:     uuid.New().String(),
		Name:   payload.Name,
		UserId: claim.UserID,
		Status: model.TodoStatusNotStarted,
	}

	// save todo
	err = u.todoRepository.Create(ctx, &newTodo)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data:    newTodo,
		Message: "success",
		Status:  http.StatusCreated,
	}
}

func (u *todoUsecase) Update(ctx context.Context, claim model.JWTClaimUser, todoID string, payload request.UpdateTodoRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check todo exist
	todo, err := u.todoRepository.FindOne(ctx, map[string]interface{}{
		"id":      todoID,
		"user_id": claim.UserID,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if todo == nil {
		return helpers.Response{
			Data:    nil,
			Message: "todo not found",
			Status:  http.StatusBadRequest,
		}
	}

	// validate payload
	validationResponse, err := helpers.ValidateBody(u.validate, payload)
	if err != nil {
		return validationResponse
	}

	// update todo
	todo.Name = payload.Name
	todo.Status = payload.Status

	// save todo
	err = u.todoRepository.Update(ctx, todo)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data:    todo,
		Message: "success",
		Status:  http.StatusOK,
	}
}

func (u *todoUsecase) Delete(ctx context.Context, claim model.JWTClaimUser, todoID string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check todo exist
	todo, err := u.todoRepository.FindOne(ctx, map[string]interface{}{
		"id":      todoID,
		"user_id": claim.UserID,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if todo == nil {
		return helpers.Response{
			Data:    nil,
			Message: "todo not found",
			Status:  http.StatusBadRequest,
		}
	}

	// delete todo
	err = u.todoRepository.Delete(ctx, todo)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data:    nil,
		Message: "todo successfully deleted",
		Status:  http.StatusOK,
	}
}
