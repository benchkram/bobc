package model

import "time"

type Artifact struct {
	ID        string `gorm:"primaryKey;"`
	ProjectID string `gorm:"primaryKey;column:project_id;not null" sql:"type:uuid;"`
	Size      int    `gorm:"column:size;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Artifact) TableName() string {
	return "artifacts"
}
