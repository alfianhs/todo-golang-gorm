package model

import (
	"time"
)

type User struct {
	ID        string     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	AvatarID  *string    `gorm:"column:avatar_id;type:uuid" json:"avatar_id"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email     string     `gorm:"column:email;type:varchar(255);not null;unique" json:"email"`
	Password  string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"-"`

	Todos  []Todo `gorm:"foreignKey:user_id;references:id" json:"todo,omitempty"`
	Avatar *File  `gorm:"foreignKey:avatar_id;references:id;constraint:OnDelete:SET NULL" json:"avatar,omitempty"`
}

func (m *User) TableName() string {
	return "users"
}

type UserRelation string

const (
	UserRelationTodo UserRelation = "Todo"
	UserRelationFile UserRelation = "Avatar"
)
