package application

import (
	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/bobc/pkg/project"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateOrUpdate(project *project.P) error

	Project(id uuid.UUID) (*project.P, error)
	ProjectByName(name string) (*project.P, error)
	Projects() ([]*project.P, error)
	ProjectDelete(id uuid.UUID) error

	CreateArtifact(projectID uuid.UUID, artifactID string, filePath string, size int) error
	ProjectArtifact(projectID uuid.UUID, artifactID string) (*artifact.A, error)
	ProjectArtifactDelete(projectID uuid.UUID, artifactID string) error

	ProjectArtifactExists(projectID uuid.UUID, artifactID string) (bool, error)
}
