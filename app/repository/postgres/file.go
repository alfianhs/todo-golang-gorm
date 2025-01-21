package postgresrepo

import (
	"context"
	"golang-gorm/domain/model"
	"golang-gorm/helpers"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

type FileRepository interface {
	FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.File, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	FindOne(ctx context.Context, filters map[string]interface{}) (*model.File, error)
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, file *model.File) error
}

func (r *fileRepository) queryFilter(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	query = helpers.CommonFilter(query, filters)

	// filters

	return query
}

func (r *fileRepository) FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.File, error) {
	var files []*model.File

	err := r.queryFilter(r.db.WithContext(ctx), filters).
		Offset(offset).
		Limit(limit).
		Find(&files).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return files, nil
}

func (r *fileRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64

	err := r.queryFilter(r.db.WithContext(ctx), filters).Model(&model.File{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *fileRepository) FindOne(ctx context.Context, filters map[string]interface{}) (*model.File, error) {
	var file model.File

	err := r.queryFilter(r.db.WithContext(ctx), filters).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &file, nil
}

func (r *fileRepository) Create(ctx context.Context, file *model.File) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) Update(ctx context.Context, file *model.File) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(file).Error
}
