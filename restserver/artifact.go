package restserver

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benchkram/bobc/application"
	projectRepo "github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/errz"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UploadArtifact creates a new artifact inside a project
// (POST /api/project/{projectName}/artifacts
func (s *S) UploadArtifact(ctx echo.Context, projectName string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	p, err := s.app.ProjectByName(projectName)
	if errors.Is(err, application.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Println("Uploading artifact...", p.ID)

	return s.upload(ctx, p.ID)
}

func (s *S) upload(ctx echo.Context, projectID uuid.UUID) (err error) {
	defer errz.Recover(&err)
	println("upload")

	r := ctx.Request()
	err = r.ParseMultipartForm(0)
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	id := r.FormValue("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is empty")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	f, err := ioutil.TempFile("", "bob-server-upload-*")
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var count int64
	n, err := io.Copy(f, file)
	count += n

	fmt.Printf("received... %d", count)
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	fi, err := f.Stat()
	if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	f.Close()
	defer os.Remove(f.Name())

	fmt.Printf("Creating artifact: [projectId: %s, artifactId: %s]\n", projectID.String(), id)

	err = s.app.ProjectArtifactCreate(projectID, id, f.Name(), int(fi.Size()))
	if err != nil {
		if errors.Is(err, application.ErrProjectNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, application.ErrProjectNotFound)
		} else if errors.Is(err, application.ErrArtifactAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, application.ErrArtifactAlreadyExists)
		} else {
			errz.Log(err)
			return echo.NewHTTPError(http.StatusInternalServerError, nil)
		}
	}

	fmt.Println("Artifact created.")

	return nil
}

// ProjectArtifactExists returns http.StatusConflict if artifact exists under a project,
// else http.StatusOK.
// (HEAD /api/project/{projectName}/artifact/{artifactId})
func (s *S) ProjectArtifactExists(ctx echo.Context, projectName, artifactId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	p, err := s.app.ProjectByName(projectName)
	if errors.Is(err, application.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	exists, err := s.app.ProjectArtifactExists(p.ID, artifactId)
	if err != nil {
		if errors.Is(err, projectRepo.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, nil)
		} else {
			errz.Log(err)
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	ctx.Response().Header().Set(HeaderBobExists, strconv.FormatBool(exists))

	return ctx.JSON(http.StatusOK, nil)
}

// GetProjectArtifact returns specific project artifact
// (GET /api/project/{projectName}/artifact/{artifactId})
func (s *S) GetProjectArtifact(ctx echo.Context, projectName, artifactId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	p, err := s.app.ProjectByName(projectName)
	if errors.Is(err, application.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	h, err := s.app.ProjectArtifact(p.ID, artifactId)
	if err != nil {
		if errors.Is(err, projectRepo.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, nil)
		} else if errors.Is(err, projectRepo.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, nil)
		} else {
			errz.Log(err)
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	return ctx.JSON(http.StatusOK, h.ToRestType())
}

// GetProjectArtifacts returns a list of all the artifacts of a project
// (GET /api/project/{projectName}/artifacts)
func (s *S) GetProjectArtifacts(ctx echo.Context, projectName string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	p, err := s.app.ProjectByName(projectName)
	if errors.Is(err, application.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	ids := []string{}
	for _, a := range p.Artifacts {
		ids = append(ids, a.ID)
	}

	return ctx.JSON(http.StatusOK, ids)
}

// DeleteProjectArtifact deletes a project artifact
// (DELETE /api/project/{projectName}/artifact/{artifactId})
func (s *S) DeleteProjectArtifact(ctx echo.Context, projectName, artifactId string) (err error) {
	defer errz.Recover(&err)

	err = s.authenticator.Authenticate(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	p, err := s.app.ProjectByName(projectName)
	if errors.Is(err, application.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound)
	} else if err != nil {
		errz.Log(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = s.app.ProjectArtifactDelete(p.ID, artifactId)
	if err != nil {
		if errors.Is(err, projectRepo.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, nil)
		} else {
			errz.Log(err)
			return ctx.JSON(http.StatusInternalServerError, nil)
		}
	}

	return ctx.JSON(http.StatusOK, nil)
}
