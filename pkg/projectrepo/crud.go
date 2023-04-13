package projectrepo

import (
	"errors"

	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/bobc/pkg/db/model"
	"github.com/benchkram/bobc/pkg/project"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
)

func (r *Repository) CreateOrUpdate(project *project.P) (err error) {
	defer errz.Recover(&err)

	var projectExists bool

	_, err = r.Project(project.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			projectExists = false
		} else {
			errz.Fatal(err)
		}
	} else {
		projectExists = true
	}

	// update project
	if projectExists {
		err = r.db.Gorm().Save(project.ToProjectDatabaseType()).Error
		errz.Fatal(err)
		return nil
	}

	// create new project
	err = r.db.Gorm().Create(project.ToProjectDatabaseType()).Error
	errz.Fatal(err)

	return nil
}

func (r *Repository) Project(projectID uuid.UUID) (_ *project.P, err error) {
	defer errz.Recover(&err)

	projectGorm := model.Project{}
	result := r.db.Gorm().Where(&model.Project{
		ID: projectID.String(),
	}).Find(&projectGorm)
	errz.Fatal(result.Error)

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	// does not work
	// r.db.Gorm().Model(&projectGorm).Association("Artifact").Find(&projectGorm.Hashes)

	// does work, but creates warning unsupported relations for schema Artifact. Needs to be figured out
	r.db.Gorm().Where("project_id=?", projectGorm.ID).Find(&projectGorm.Artifacts)

	// does works fine. But Not safe to use raw query
	// r.db.Gorm().Raw("SELECT * from project_hashes WHERE project_id=?", projectGorm.ID).Find(&projectGorm.Hashes)

	return project.FromDBModel(&projectGorm)
}

func (r *Repository) ProjectsByName(name string) ([]*project.P, error) {
	var projects []*project.P

	ps := []model.Project{}

	err := r.db.Gorm().Where("name LIKE ?", "%"+name+"%").Find(&ps).Error
	errz.Fatal(err)
	for _, p := range ps {
		// does work, but creates warning unsupported relations for schema Artifact. Needs to be figured out
		err := r.db.Gorm().Where("project_id=?", p.ID).Find(&p.Artifacts).Error
		errz.Fatal(err)

		pOut, err := project.FromDBModel(&p)
		errz.Fatal(err)

		projects = append(projects, pOut)
	}

	return projects, nil
}

func (r *Repository) ProjectByName(projectName string) (*project.P, error) {
	var projectGorm model.Project

	result := r.db.Gorm().
		Where(&model.Project{Name: projectName}).
		Find(&projectGorm)
	errz.Fatal(result.Error)

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	// does work, but creates warning unsupported relations for schema Artifact. Needs to be figured out
	err := r.db.Gorm().Where("project_id=?", projectGorm.ID).Find(&projectGorm.Artifacts).Error
	errz.Fatal(err)

	return project.FromDBModel(&projectGorm)
}

func (r *Repository) Projects() (projects []*project.P, err error) {
	defer errz.Recover(&err)

	projects = []*project.P{}

	ps := []model.Project{}

	err = r.db.Gorm().Find(&ps).Error
	errz.Fatal(err)
	for _, p := range ps {
		// does work, but creates warning unsupported relations for schema Artifact. Needs to be figured out
		err := r.db.Gorm().Where("project_id=?", p.ID).Find(&p.Artifacts).Error
		errz.Fatal(err)

		pOut, err := project.FromDBModel(&p)
		errz.Fatal(err)
		projects = append(projects, pOut)
	}

	return projects, nil
}

func (r *Repository) ProjectDelete(projectID uuid.UUID) error {
	result := r.db.Gorm().Delete(&model.Project{
		ID: projectID.String(),
	})
	errz.Fatal(result.Error)

	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) CreateArtifact(projectID uuid.UUID, artifactID string, filePath string, size int) (err error) {
	defer errz.Recover(&err)

	var projectExists bool

	p, err := r.Project(projectID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			projectExists = false
		} else {
			errz.Fatal(err)
		}
	} else {
		projectExists = true
	}

	if !projectExists {
		return ErrNotFound
	}

	h := model.Artifact{
		ID:         uuid.New().String(),
		ArtifactID: artifactID,
		ProjectID:  p.ID.String(),
		Size:       size,
	}

	err = r.db.Gorm().Create(&h).Error
	errz.Fatal(err)

	err = r.artifactStore.CreateArtifact(h.ID, filePath, size)
	errz.Fatal(err)

	return nil
}

func (r *Repository) artifact(projectID uuid.UUID, artifactID string) (_ *model.Artifact, err error) {
	defer errz.Recover(&err)

	hashGorm := &model.Artifact{}
	result := r.db.Gorm().Where(&model.Artifact{
		ProjectID:  projectID.String(),
		ArtifactID: artifactID,
	}).Find(hashGorm)
	errz.Fatal(result.Error)

	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	return hashGorm, nil
}

func (r *Repository) ArtifactUpdate(projectID uuid.UUID, artifactID string) (err error) {
	defer errz.Recover(&err)

	hashGorm, err := r.artifact(projectID, artifactID)
	errz.Fatal(err)

	// hashGorm.StoragePath = hash.StoragePath

	err = r.db.Gorm().Save(hashGorm).Error
	errz.Fatal(err)

	return nil
}

func (r *Repository) ProjectArtifact(projectID uuid.UUID, artifactID string) (_ *artifact.A, err error) {
	defer errz.Recover(&err)

	artiGorm, err := r.artifact(projectID, artifactID)
	if err != nil {
		return nil, err
	}

	arti := artifact.FromDatabaseType(artiGorm)

	addr, err := r.artifactStore.Artifact(artiGorm.ID)
	errz.Fatal(err)

	arti.AccessLink = addr

	return arti, nil
}

func (r *Repository) ProjectArtifactDelete(projectID uuid.UUID, artifactID string) (err error) {
	defer errz.Recover(&err)

	hashGorm, err := r.artifact(projectID, artifactID)
	errz.Fatal(err)

	err = r.db.Gorm().Delete(hashGorm).Error
	errz.Fatal(err)

	err = r.artifactStore.DeleteArtifact(hashGorm.ID)
	errz.Fatal(err)

	return nil
}

func (r *Repository) ProjectArtifactExists(projectID uuid.UUID, artifactID string) (_ bool, err error) {
	defer errz.Recover(&err)

	_, err = r.artifact(projectID, artifactID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		} else {
			errz.Fatal(err)
		}
	}

	return true, nil

}
