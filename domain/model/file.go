package model

import "time"

type File struct {
	ID        string     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	MimeType  string     `gorm:"column:mime_type;type:varchar(255);not null" json:"mime_type"`
	Size      int64      `gorm:"column:size;type:bigint;not null" json:"size"`
	Url       string     `gorm:"column:url;type:varchar(255);not null" json:"url"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"-"`

	User []User `gorm:"foreignKey:avatar_id;references:id;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

func (m *File) TableName() string {
	return "files"
}
