package postgresrepo

import (
	"context"
	"golang-gorm/domain/model"
	"golang-gorm/helpers"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

type TodoRepository interface {
	FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.Todo, error)
	FindOne(ctx context.Context, filters map[string]interface{}) (*model.Todo, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Create(ctx context.Context, todo *model.Todo) error
	UpdateOne(ctx context.Context, todo *model.Todo) error
	DeleteOne(ctx context.Context, todo *model.Todo) error
}

func (r *todoRepository) queryFilter(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	query = helpers.CommonFilter(query, filters)

	// filters
	if userID, ok := filters["user_id"].(string); ok {
		query = query.Where("user_id = ?", userID)
	}

	return query
}

func (r *todoRepository) FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.Todo, error) {
	var todos []*model.Todo

	err := r.queryFilter(r.db.WithContext(ctx), filters).
		Offset(offset).
		Limit(limit).
		Find(&todos).Error
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64

	err := r.queryFilter(r.db.WithContext(ctx), filters).Model(&model.Todo{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *todoRepository) FindOne(ctx context.Context, filters map[string]interface{}) (*model.Todo, error) {
	var todo model.Todo

	err := r.queryFilter(r.db.WithContext(ctx), filters).First(&todo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) UpdateOne(ctx context.Context, todo *model.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) DeleteOne(ctx context.Context, todo *model.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(todo).UpdateColumn("deleted_at", time.Now()).Error
}
