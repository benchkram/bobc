package projectrepo

import (
	"fmt"
	"net/url"

	"github.com/benchkram/bobc/pkg/db"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type ArtifactStore interface {
	CreateArtifact(id, filePath string, size int) (err error)
	DeleteArtifact(id string) (err error)
	Artifact(id string) (addr *url.URL, err error)
}

type Repository struct {
	db            db.Database
	artifactStore ArtifactStore
}

func New(db db.Database, artifactStore ArtifactStore) *Repository {
	return &Repository{
		db:            db,
		artifactStore: artifactStore,
	}
}
