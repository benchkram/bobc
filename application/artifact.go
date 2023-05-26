package application

import (
	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
)

// ProjectArtifactCreate creates a new artifact and copies src to the internal storage.
func (a *application) ProjectArtifactCreate(projectID uuid.UUID, artifactID string, src string, size int) (err error) {
	defer errz.Recover(&err)

	exists, err := a.ProjectArtifactExists(projectID, artifactID)
	errz.Fatal(err)

	if exists {
		return ErrArtifactAlreadyExists
	}

	err = a.projects.CreateArtifact(projectID, artifactID, src, size)
	errz.Fatal(err)

	return nil
}

// ProjectArtifactDelete deletes a artifact from database and s3 storage, does nothing if artifact does not exists
func (a *application) ProjectArtifactDelete(projectID uuid.UUID, artifactID string) (err error) {
	defer errz.Recover(&err)

	_, err = a.Project(projectID)
	errz.Fatal(err)

	err = a.projects.ProjectArtifactDelete(projectID, artifactID)
	errz.Fatal(err)

	return nil
}

func (a *application) ProjectArtifactExists(projectID uuid.UUID, artifactID string) (_ bool, err error) {
	defer errz.Recover(&err)

	_, err = a.Project(projectID)
	if err != nil {
		return false, err
	}

	return a.projects.ProjectArtifactExists(projectID, artifactID)
}

func (a *application) ProjectArtifact(projectID uuid.UUID, artifactID string) (_ *artifact.A, err error) {
	defer errz.Recover(&err)

	_, err = a.Project(projectID)
	errz.Fatal(err)

	artifact, err := a.projects.ProjectArtifact(projectID, artifactID)
	errz.Fatal(err)

	return artifact, nil
}
