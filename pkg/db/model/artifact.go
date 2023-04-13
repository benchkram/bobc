package model

import "time"

type Artifact struct {
	ID         string `gorm:"primaryKey"`
	ProjectID  string `gorm:"column:project_id;not null;index" sql:"type:uuid"`
	ArtifactID string `gorm:"column:artifact_id;not null;index"`
	Size       int    `gorm:"column:size;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Artifact) TableName() string {
	return "artifacts"
}
