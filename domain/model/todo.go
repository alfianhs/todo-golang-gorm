package model

import "time"

type Todo struct {
	ID        string     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	UserId    string     `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Status    TodoStatus `gorm:"column:status;type:todo_status;not null" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"-"`

	User User `gorm:"foreignKey:user_id;references:id" json:"-"`
}

type TodoStatus string

const (
	TodoStatusDone       TodoStatus = "Done"
	TodoStatusNotStarted TodoStatus = "NotStarted"
)
