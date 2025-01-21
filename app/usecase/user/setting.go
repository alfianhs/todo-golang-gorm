package usecase_user

import (
	"context"
	"encoding/base64"
	"fmt"
	postgresrepo "golang-gorm/app/repository/postgres"
	s3repo "golang-gorm/app/repository/s3"
	"golang-gorm/app/usecase"
	"golang-gorm/domain/model"
	"golang-gorm/domain/request"
	"golang-gorm/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type settingUsecase struct {
	userRepository postgresrepo.UserRepository
	fileRepository postgresrepo.FileRepository
	s3Repository   s3repo.S3Repo
	contextTimeout time.Duration
	validate       *validator.Validate
}

func NewSettingUsecase(d usecase.UsecaseDependency) SettingUsecase {
	return &settingUsecase{
		userRepository: d.UserRepository,
		fileRepository: d.FileRepository,
		s3Repository:   d.S3Repository,
		contextTimeout: d.Timeout,
		validate:       d.Validate,
	}
}

type SettingUsecase interface {
	UpdateProfile(ctx context.Context, claim model.JWTClaimUser, payload request.UserUpdateProfileRequest) helpers.Response
}

func (u *settingUsecase) UpdateProfile(ctx context.Context, claim model.JWTClaimUser, payload request.UserUpdateProfileRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	validationResponse, err := helpers.ValidateBody(u.validate, payload)
	if err != nil {
		return validationResponse
	}

	// check user exist
	user, err := u.userRepository.FindOne(ctx, map[string]interface{}{
		"user_id": claim.UserID,
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

	// upload photo profile
	file, err := u.uploadProfilePicture(ctx, payload.ProfilePicture, user.Name)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// update user
	user.Name = payload.Name
	user.AvatarID = &file.ID

	// save user
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return helpers.Response{
			Data:    nil,
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return helpers.Response{
		Data:    user,
		Message: "success",
		Status:  http.StatusOK,
	}
}

func (u *settingUsecase) uploadProfilePicture(ctx context.Context, base64Data string, name string) (*model.File, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	logrus.Info("upload profile picture")

	// split metadata and data base64
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid base64 data")
	}

	// get mimetype
	meta := parts[0]
	mimeType := strings.TrimPrefix(strings.Split(meta, ";")[0], "data:")
	if mimeType == "" {
		return nil, fmt.Errorf("unable to parse MIME type")
	}

	// validate mimetype
	allowed := helpers.IsMimeTypeAllowed(mimeType, "image")
	if !allowed {
		return nil, fmt.Errorf("mimetype not allowed")
	}

	// get extension
	ext, err := helpers.GetExtensionFromMimeType(mimeType)
	if err != nil {
		return nil, err
	}

	// decode base64
	file, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// get file size
	fileSize := int64(len(file))

	// validate size
	maxSize := int64(10 * 1024 * 1024)
	if fileSize > maxSize {
		return nil, fmt.Errorf("file size must be less than 10MB")
	}

	// set objectName
	fileName := fmt.Sprintf("%s_%s.%s", time.Now().Format("20060102"), name, ext)
	year, month, _ := time.Now().Date()
	objectName := fmt.Sprintf("profile_pictures/%d/%s/%s", year, month, fileName)

	logrus.Info(objectName)

	// upload to s3
	url, err := u.s3Repository.UploadFile(ctx, objectName, file, mimeType)
	if err != nil {
		return nil, err
	}

	logrus.Info(url)

	// save to database
	newFile := model.File{
		ID:       uuid.New().String(),
		Name:     fileName,
		MimeType: mimeType,
		Size:     fileSize,
		Url:      url,
	}

	err = u.fileRepository.Create(ctx, &newFile)
	if err != nil {
		return nil, err
	}

	return &newFile, nil
}
