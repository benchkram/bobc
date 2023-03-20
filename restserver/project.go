package restserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/benchkram/bobc/application"
	projectRepo "github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/bobc/restserver/generated"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	HeaderBobExists = "Bob-Exists"
)

// Returns a list of projects with name and ID, without hashes.
// (GET /api/projects)
func (s *S) GetProjects(ctx echo.Context) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	projects, err := s.app.Projects()
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	result := []generated.Project{}
	for _, project := range projects {
		p := project.ToProjectRestType()
		result = append(result, p)
	}

	return ctx.JSON(http.StatusOK, result)
}

// Create a new project by name.
// also adds the hashes after the creation of project
// (POST /api/project)
func (s *S) CreateProject(ctx echo.Context) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	projectCreate := generated.ProjectCreate{}
	err = ctx.Bind(&projectCreate)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	if projectCreate.Name == "" {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	p, err := s.app.ProjectCreate(projectCreate.Name, projectCreate.Description)
	if errors.Is(err, application.ErrProjectAlreadyExists) {
		return ctx.NoContent(http.StatusBadRequest)
	} else if errors.Is(err, application.ErrInvalidProjectName) {
		return ctx.NoContent(http.StatusBadRequest)
	} else if err != nil {
		errz.Log(err)
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, p.ToExtendedProjectRestType())
}

func (s *S) ProjectExists(ctx echo.Context, projectId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	// rest requires all parameters for a route to be named equally.
	projectName := projectId

	exists, err := s.app.ProjectExists(projectName)
	if err != nil {
		errz.Log(err)
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	ctx.Response().Header().Set(HeaderBobExists, strconv.FormatBool(exists))
	return ctx.JSON(http.StatusOK, nil)
}

// Delete a project by id.
// (DELETE /api/project/{project_id})
func (s *S) DeleteProject(ctx echo.Context, projectId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrInvalidProjectID.Error())
	}

	p, err := s.app.Project(pid)
	if err != nil {
		if errors.Is(err, projectRepo.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, nil)
		} else {
			errz.Log(err)
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	err = s.app.ProjectDelete(p.ID)
	if err != nil {
		errz.Log(err)
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, nil)
}

// Returns a single project by id.
// (GET /api/project/{project_id})
func (s *S) GetProject(ctx echo.Context, projectId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	pid, err := uuid.Parse(projectId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrInvalidProjectID)
	}

	project, err := s.app.Project(pid)
	if err != nil {
		if errors.Is(err, application.ErrProjectNotFound) {
			return ctx.NoContent(http.StatusNotFound)
		} else {
			errz.Log(err)
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	return ctx.JSON(http.StatusOK, project.ToExtendedProjectRestType())
}
