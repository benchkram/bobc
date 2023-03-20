package application

import (
	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
)

// ProjectArtifactCreate creates a new artifact and copies src to the internal storage.
func (s *application) ProjectArtifactCreate(projectID uuid.UUID, artifactID string, src string, size int) (err error) {
	defer errz.Recover(&err)

	exists, err := s.ProjectArtifactExists(projectID, artifactID)
	errz.Fatal(err)

	if exists {
		return ErrArtifactAlreadyExists
	}

	err = s.projects.CreateArtifact(projectID, artifactID, src, size)
	errz.Fatal(err)

	return nil
}

// ProjectArtifactDelete deletes a artifact from database and s3 storage, does nothing if artifact does not exists
func (s *application) ProjectArtifactDelete(projectID uuid.UUID, artifactID string) (err error) {
	defer errz.Recover(&err)

	_, err = s.Project(projectID)
	errz.Fatal(err)

	err = s.projects.ProjectArtifactDelete(projectID, artifactID)
	errz.Fatal(err)

	return nil
}

func (s *application) ProjectArtifactExists(projectID uuid.UUID, artifactID string) (_ bool, err error) {
	defer errz.Recover(&err)

	_, err = s.Project(projectID)
	if err != nil {
		return false, err
	}

	return s.projects.ProjectArtifactExists(projectID, artifactID)
}

func (s *application) ProjectArtifact(projectID uuid.UUID, artifactID string) (_ *artifact.A, err error) {
	defer errz.Recover(&err)

	_, err = s.Project(projectID)
	errz.Fatal(err)

	artifact, err := s.projects.ProjectArtifact(projectID, artifactID)
	errz.Fatal(err)

	return artifact, nil
}
