package restserverclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/benchkram/bobc/restserver/generated"
	"github.com/benchkram/errz"
)

var (
	ErrCantMakeRequest     = fmt.Errorf("failed to make the request")
	ErrInvalidStatusCode   = fmt.Errorf("failed request, invalid status code")
	ErrCantReadResponse    = fmt.Errorf("can not read response from server")
	ErrInvalidJsonResponse = fmt.Errorf("invalid response from the server")
	ErrItemNotFound        = fmt.Errorf("desired Item not found")
)

const (
	HeaderBobExists = "Bob-Exists"
)

type C struct {
	client *generated.ClientWithResponses
}

// Creates New Client from the address without the protocol
func New(address string, apiKey []byte) (*C, error) {
	c, err := generated.NewClientWithResponses(
		"http://"+address,
		generated.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+string(apiKey))
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return &C{
		client: c,
	}, nil
}

func (c *C) Health() bool {
	response, err := c.client.GetHealth(context.Background())
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (c *C) Project(projectId string) (*generated.ExtendedProject, error) {

	response, err := c.client.GetProject(context.Background(), projectId)
	if err != nil {
		return nil, ErrCantMakeRequest
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, ErrCantReadResponse
	}

	var project generated.ExtendedProject
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, ErrInvalidJsonResponse
	}
	return &project, nil
}

func (c *C) Projects() ([]generated.Project, error) {
	response, err := c.client.GetProjects(context.Background())
	if err != nil {
		return nil, ErrCantMakeRequest
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, ErrCantReadResponse
	}

	var projects []generated.Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, ErrInvalidJsonResponse
	}
	return projects, nil
}

func (c *C) ProjectCreate(project generated.ProjectCreate) (*generated.ExtendedProject, error) {

	response, err := c.client.CreateProjectWithResponse(
		context.Background(),
		generated.CreateProjectJSONRequestBody(project),
	)
	if err != nil {
		return nil, ErrCantMakeRequest
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("create project failed %w, http status %d, body: %s", ErrInvalidStatusCode, response.StatusCode(), string(response.Body))
	}

	return response.JSON200, nil
}

func (c *C) ProjectExists(name string) (bool, error) {
	response, err := c.client.ProjectExistsWithResponse(
		context.Background(),
		name,
	)
	if err != nil {
		return false, fmt.Errorf("project exists failed %w", err)
	}

	if response.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("project exists failed %w", ErrInvalidStatusCode)
	}

	exists := response.HTTPResponse.Header.Get(HeaderBobExists)
	e, err := strconv.ParseBool(exists)
	if err != nil {
		return false, fmt.Errorf("project exists failed %w", err)
	}

	return e, nil
}

func (c *C) ProjectDelete(projectId string) error {
	response, err := c.client.DeleteProject(context.Background(), projectId)
	if err != nil {
		return ErrCantMakeRequest
	}

	if response.StatusCode == http.StatusNotFound {
		return ErrItemNotFound
	}

	if response.StatusCode != http.StatusOK {
		return ErrInvalidStatusCode
	}

	return nil
}

func (c *C) Artifact(projectId string, hash string) (*generated.Artifact, error) {
	response, err := c.client.GetProjectArtifact(context.Background(), projectId, hash)
	if err != nil {
		return nil, ErrCantMakeRequest
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[code: %d], %w", response.StatusCode, ErrInvalidStatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, ErrCantReadResponse
	}

	var artifact generated.Artifact
	if err := json.Unmarshal(body, &artifact); err != nil {
		return nil, ErrInvalidJsonResponse
	}
	return &artifact, nil
}

func (c *C) ArtifactCreate(projectId string, hash string, src string) (err error) {
	defer errz.Recover(&err)

	f, err := os.Open(src)
	errz.Fatal(err)
	defer f.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err = w.WriteField("id", hash)
	errz.Fatal(err)

	fieldWriter, err := w.CreateFormFile("file", hash)
	_, err = io.Copy(fieldWriter, f)
	errz.Fatal(err)

	w.Close()

	response, err := c.client.UploadArtifactWithBodyWithResponse(
		context.Background(),
		projectId,
		w.FormDataContentType(),
		&b,
	)
	if err != nil {
		return ErrCantMakeRequest
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("[code: %d], %w", response.StatusCode(), ErrInvalidStatusCode)
	}

	return nil
}

// func (c *C) UpdateHash(projectId string, hash string, update generated.ProjectArtifactUpdate) (*generated.ProjectArtifact, error) {
// 	body := generated.UpdateProjectArtifactJSONRequestBody(update)
// 	response, err := c.client.UpdateProjectArtifact(context.Background(), projectId, hash, body)
// 	if err != nil {
// 		return nil, ErrCantMakeRequest
// 	}

// 	if response.StatusCode != http.StatusOK {
// 		return nil, ErrInvalidStatusCode
// 	}

// 	responseBody, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, ErrCantReadResponse
// 	}

// 	var projectHash generated.ProjectArtifact
// 	if err := json.Unmarshal(responseBody, &projectHash); err != nil {
// 		return nil, ErrInvalidJsonResponse
// 	}
// 	return &projectHash, nil
// }

func (c *C) ArtifactDelete(projectId string, hash string) error {

	response, err := c.client.DeleteProjectArtifact(context.Background(), projectId, hash)
	if err != nil {
		return ErrCantMakeRequest
	}

	fmt.Println(response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return ErrInvalidStatusCode
	}

	return nil
}

func (c *C) ArtifactExists(projectId string, hash string) (bool, error) {
	response, err := c.client.ProjectArtifactExistsWithResponse(context.Background(), projectId, hash)

	if err != nil {
		return false, fmt.Errorf("project hash exists failed %w", err)
	}

	if response.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("project hash exists failed %w", ErrInvalidStatusCode)
	}

	exists := response.HTTPResponse.Header.Get(HeaderBobExists)
	e, err := strconv.ParseBool(exists)
	if err != nil {
		return false, fmt.Errorf("project hash exists failed %w", err)
	}

	return e, nil
}
