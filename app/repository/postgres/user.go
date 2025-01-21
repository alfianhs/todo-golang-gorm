package postgresrepo

import (
	"context"

	"golang-gorm/domain/model"
	"golang-gorm/helpers"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type UserRepository interface {
	FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, error)
	FindOne(ctx context.Context, filters map[string]interface{}) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

func (r *userRepository) queryFilter(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	query = helpers.CommonFilter(query, filters)

	// filters
	if email, ok := filters["email"].(string); ok {
		query = query.Where("email = ?", email)
	}

	return query
}

func (r *userRepository) FetchList(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.User, error) {
	var users []*model.User

	err := r.queryFilter(r.db.WithContext(ctx), filters).
		Preload(string(model.UserRelationFile)).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) FindOne(ctx context.Context, filters map[string]interface{}) (*model.User, error) {
	var user model.User
	query := r.queryFilter(r.db.WithContext(ctx), filters)

	err := query.Preload(string(model.UserRelationFile)).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(user).Error
}
