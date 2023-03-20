package application

import (
	"sync"

	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/bobc/pkg/project"
	"github.com/google/uuid"
)

type Application interface {
	Projects() (_ []*project.P, err error)
	Project(id uuid.UUID) (*project.P, error)
	ProjectByName(name string) (_ *project.P, err error)
	//ProjectsByName(name string) ([]*project.P, error)
	ProjectExists(name string) (bool, error)
	ProjectCreate(name, description string) (*project.P, error)
	ProjectDelete(id uuid.UUID) error

	ProjectArtifact(projectID uuid.UUID, artifactID string) (*artifact.A, error)
	ProjectArtifactExists(projectID uuid.UUID, artifactID string) (bool, error)
	ProjectArtifactCreate(projectID uuid.UUID, artifactID string, src string, size int) error
	ProjectArtifactDelete(projectID uuid.UUID, artifactID string) error
}

type application struct {
	// projects is the storage abstraction for projects
	projects ProjectRepository

	// mux is used to not allow specific operations to be called in parallel
	mux sync.Mutex
}

func New(opts ...Option) Application {
	// intialize defaults here
	app := &application{}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

	return app
}
