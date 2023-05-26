package application

import (
	"errors"
	"regexp"

	"github.com/benchkram/bobc/pkg/project"
	"github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
)

func (a *application) ProjectCreate(name, description string) (_ *project.P, err error) {
	defer errz.Recover(&err)

	if !a.projectNameValid(name) {
		return nil, ErrInvalidProjectName
	}

	//_, err = a.projects.ProjectByName(name)
	//if err == nil {
	//	return nil, ErrProjectAlreadyExists
	//} else if !errors.Is(err, projectrepo.ErrNotFound) {
	//	errz.Fatal(err)
	//}

	p := project.New(name, description)
	err = a.projects.CreateOrUpdate(p)
	errz.Fatal(err)

	return p, nil
}

// projectNameValid reports if the project name passed is valid. Valid names
// are considered those that only contain alphanumerics, hyphens, periods and
// underscores, and have a length of at least 1 and at most 100 unicode
// codepoints. The names `.` and `..` are not allowed.
//
//	Ref: https://github.com/dead-claudia/github-limits#repository-names
func (a *application) projectNameValid(name string) bool {
	// alphanumerics, hyphens, periods, underscores. length 1-100 codepoints
	rex := regexp.MustCompile(`^[A-Za-z0-9-_.]{1,100}$`)

	return rex.MatchString(name) && name != "." && name != ".."
}

func (a *application) Project(id uuid.UUID) (_ *project.P, err error) {
	defer errz.Recover(&err)

	p, err := a.projects.Project(id)
	if errors.Is(err, projectrepo.ErrNotFound) {
		return nil, ErrProjectNotFound
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (a *application) Projects() (_ []*project.P, err error) {
	defer errz.Recover(&err)

	return a.projects.Projects()
}

func (a *application) ProjectByName(name string) (_ *project.P, err error) {
	defer errz.Recover(&err)

	p, err := a.projects.ProjectByName(name)
	if errors.Is(err, projectrepo.ErrNotFound) {
		return nil, ErrProjectNotFound
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

func (a *application) ProjectExists(name string) (exists bool, err error) {
	defer errz.Recover(&err)

	_, err = a.projects.ProjectByName(name)
	if err == nil {
		exists = true
	} else {
		if errors.Is(err, projectrepo.ErrNotFound) {
			exists = false
		} else {
			errz.Fatal(err)
		}
	}

	return exists, nil
}

func (a *application) ProjectDelete(projectID uuid.UUID) (err error) {
	defer errz.Recover(&err)

	// Delete from database
	err = a.projects.ProjectDelete(projectID)
	if errors.Is(err, projectrepo.ErrNotFound) {
		return ErrProjectNotFound
	}
	errz.Fatal(err)

	return nil
}
