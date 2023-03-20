package model

import "time"

type Project struct {
	ID string `gorm:"primaryKey;" sql:"type:uuid;"`

	Name        string      `gorm:"column:name;not null"`
	Description string      `gorm:"column:description;not null"`
	Artifacts   []*Artifact `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Project) TableName() string {
	return "projects"
}
