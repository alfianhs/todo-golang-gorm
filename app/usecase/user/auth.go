package usecase_user

import (
	"context"
	"net/http"
	"time"

	postgresrepo "golang-gorm/app/repository/postgres"
	"golang-gorm/app/usecase"
	"golang-gorm/domain/model"
	"golang-gorm/domain/request"
	"golang-gorm/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository postgresrepo.UserRepository
	contextTimeout time.Duration
	validate       *validator.Validate
}

func NewAuthUsecase(d usecase.UsecaseDependency) AuthUsecase {
	return &authUsecase{
		userRepository: d.UserRepository,
		contextTimeout: d.Timeout,
		validate:       d.Validate,
	}
}

type AuthUsecase interface {
	Register(ctx context.Context, payload request.RegisterRequest) helpers.Response
	Login(ctx context.Context, payload request.LoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim model.JWTClaimUser) helpers.Response
}

func (u *authUsecase) Register(ctx context.Context, payload request.RegisterRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	validationResponse, err := helpers.ValidateBody(u.validate, payload)
	if err != nil {
		return validationResponse
	}

	// check user exist
	user, err := u.userRepository.FindOne(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if user != nil {
		return helpers.Response{
			Data:    nil,
			Message: "email already exist",
			Status:  http.StatusBadRequest,
		}
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// create user
	user = &model.User{
		ID:        uuid.New().String(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data:    user,
		Message: "user successfully registered",
		Status:  http.StatusCreated,
	}
}

func (u *authUsecase) Login(ctx context.Context, payload request.LoginRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	validationResponse, err := helpers.ValidateBody(u.validate, payload)
	if err != nil {
		return validationResponse
	}

	// check user exist
	user, err := u.userRepository.FindOne(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if user == nil {
		return helpers.Response{
			Data:    nil,
			Message: "email not found",
			Status:  http.StatusBadRequest,
		}
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: "wrong password",
			Status:  http.StatusBadRequest,
		}
	}

	// generate token
	token, err := helpers.GenerateJWTTokenUser(model.JWTClaimUser{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    "user",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(helpers.GetJWTTTL()))),
		},
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data: map[string]interface{}{
			"token": token,
			"user":  user,
		},
		Message: "login success",
		Status:  http.StatusOK,
	}
}

func (u *authUsecase) GetProfile(ctx context.Context, claim model.JWTClaimUser) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check user exist
	user, err := u.userRepository.FindOne(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if user == nil {
		return helpers.Response{
			Data:    nil,
			Message: "user not found",
			Status:  http.StatusBadRequest,
		}
	}

	return helpers.Response{
		Data:    user,
		Message: "success",
		Status:  http.StatusOK,
	}
}
