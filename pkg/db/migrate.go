package db

import (
	"log"
	"time"

	"github.com/benchkram/errz"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func (db *database) migrate() (err error) {
	defer errz.Recover(&err)

	log.Println("migrate db")

	if db.gorm == nil {
		return ErrDatabaseNil
	}

	m := gormigrate.New(db.gorm, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202303201338",
			Migrate: func(tx *gorm.DB) (err error) {
				defer errz.Recover(&err)

				if tx == nil {
					return ErrDatabaseNil
				}

				err = tx.AutoMigrate(&Artifact202303201338{})
				errz.Fatal(err)

				err = tx.AutoMigrate(&Project202303201338{})
				errz.Fatal(err)

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	return m.Migrate()
}

type Artifact202303201338 struct {
	ID        string `gorm:"primaryKey;"`
	ProjectID string `gorm:"primaryKey;column:project_id;not null" sql:"type:uuid;"`
	Size      int    `gorm:"column:size;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Artifact202303201338) TableName() string {
	return "artifacts"
}

type Project202303201338 struct {
	ID          string                  `gorm:"primaryKey;" sql:"type:uuid;"`
	Name        string                  `gorm:"column:name;not null"`
	Description string                  `gorm:"column:description;not null"`
	Artifacts   []*Artifact202303201338 `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Project202303201338) TableName() string {
	return "projects"
}
