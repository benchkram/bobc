// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package generated

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetHealth request
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteProject request
	DeleteProject(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProject request
	GetProject(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ProjectExists request
	ProjectExists(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteProjectArtifact request
	DeleteProjectArtifact(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProjectArtifact request
	GetProjectArtifact(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ProjectArtifactExists request
	ProjectArtifactExists(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProjectArtifacts request
	GetProjectArtifacts(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UploadArtifact request  with any body
	UploadArtifactWithBody(ctx context.Context, projectName string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProjects request
	GetProjects(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateProject request  with any body
	CreateProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateProject(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteProject(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteProjectRequest(c.Server, projectName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProject(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProjectRequest(c.Server, projectName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ProjectExists(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewProjectExistsRequest(c.Server, projectName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteProjectArtifact(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteProjectArtifactRequest(c.Server, projectName, artifactId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProjectArtifact(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProjectArtifactRequest(c.Server, projectName, artifactId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ProjectArtifactExists(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewProjectArtifactExistsRequest(c.Server, projectName, artifactId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProjectArtifacts(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProjectArtifactsRequest(c.Server, projectName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UploadArtifactWithBody(ctx context.Context, projectName string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUploadArtifactRequestWithBody(c.Server, projectName, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProjects(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProjectsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateProjectRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateProject(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateProjectRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetHealthRequest generates requests for GetHealth
func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/health")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDeleteProjectRequest generates requests for DeleteProject
func NewDeleteProjectRequest(server string, projectName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProjectRequest generates requests for GetProject
func NewGetProjectRequest(server string, projectName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewProjectExistsRequest generates requests for ProjectExists
func NewProjectExistsRequest(server string, projectName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("HEAD", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewDeleteProjectArtifactRequest generates requests for DeleteProjectArtifact
func NewDeleteProjectArtifactRequest(server string, projectName string, artifactId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "artifactId", runtime.ParamLocationPath, artifactId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s/artifact/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProjectArtifactRequest generates requests for GetProjectArtifact
func NewGetProjectArtifactRequest(server string, projectName string, artifactId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "artifactId", runtime.ParamLocationPath, artifactId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s/artifact/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewProjectArtifactExistsRequest generates requests for ProjectArtifactExists
func NewProjectArtifactExistsRequest(server string, projectName string, artifactId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "artifactId", runtime.ParamLocationPath, artifactId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s/artifact/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("HEAD", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetProjectArtifactsRequest generates requests for GetProjectArtifacts
func NewGetProjectArtifactsRequest(server string, projectName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s/artifacts", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUploadArtifactRequestWithBody generates requests for UploadArtifact with any type of body
func NewUploadArtifactRequestWithBody(server string, projectName string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "projectName", runtime.ParamLocationPath, projectName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/project/%s/artifacts", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetProjectsRequest generates requests for GetProjects
func NewGetProjectsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/projects")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateProjectRequest calls the generic CreateProject builder with application/json body
func NewCreateProjectRequest(server string, body CreateProjectJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateProjectRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateProjectRequestWithBody generates requests for CreateProject with any type of body
func NewCreateProjectRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/projects")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetHealth request
	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)

	// DeleteProject request
	DeleteProjectWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*DeleteProjectResponse, error)

	// GetProject request
	GetProjectWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*GetProjectResponse, error)

	// ProjectExists request
	ProjectExistsWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*ProjectExistsResponse, error)

	// DeleteProjectArtifact request
	DeleteProjectArtifactWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*DeleteProjectArtifactResponse, error)

	// GetProjectArtifact request
	GetProjectArtifactWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*GetProjectArtifactResponse, error)

	// ProjectArtifactExists request
	ProjectArtifactExistsWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*ProjectArtifactExistsResponse, error)

	// GetProjectArtifacts request
	GetProjectArtifactsWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*GetProjectArtifactsResponse, error)

	// UploadArtifact request  with any body
	UploadArtifactWithBodyWithResponse(ctx context.Context, projectName string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UploadArtifactResponse, error)

	// GetProjects request
	GetProjectsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetProjectsResponse, error)

	// CreateProject request  with any body
	CreateProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error)

	CreateProjectWithResponse(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error)
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Success
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ExtendedProject
}

// Status returns HTTPResponse.Status
func (r GetProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ProjectExistsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ProjectExistsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ProjectExistsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteProjectArtifactResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteProjectArtifactResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteProjectArtifactResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProjectArtifactResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Artifact
}

// Status returns HTTPResponse.Status
func (r GetProjectArtifactResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProjectArtifactResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ProjectArtifactExistsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ProjectArtifactExistsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ProjectArtifactExistsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProjectArtifactsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ArtifactIds
}

// Status returns HTTPResponse.Status
func (r GetProjectArtifactsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProjectArtifactsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UploadArtifactResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r UploadArtifactResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UploadArtifactResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProjectsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Project
}

// Status returns HTTPResponse.Status
func (r GetProjectsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProjectsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ExtendedProject
}

// Status returns HTTPResponse.Status
func (r CreateProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetHealthWithResponse request returning *GetHealthResponse
func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

// DeleteProjectWithResponse request returning *DeleteProjectResponse
func (c *ClientWithResponses) DeleteProjectWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*DeleteProjectResponse, error) {
	rsp, err := c.DeleteProject(ctx, projectName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteProjectResponse(rsp)
}

// GetProjectWithResponse request returning *GetProjectResponse
func (c *ClientWithResponses) GetProjectWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*GetProjectResponse, error) {
	rsp, err := c.GetProject(ctx, projectName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProjectResponse(rsp)
}

// ProjectExistsWithResponse request returning *ProjectExistsResponse
func (c *ClientWithResponses) ProjectExistsWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*ProjectExistsResponse, error) {
	rsp, err := c.ProjectExists(ctx, projectName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseProjectExistsResponse(rsp)
}

// DeleteProjectArtifactWithResponse request returning *DeleteProjectArtifactResponse
func (c *ClientWithResponses) DeleteProjectArtifactWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*DeleteProjectArtifactResponse, error) {
	rsp, err := c.DeleteProjectArtifact(ctx, projectName, artifactId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteProjectArtifactResponse(rsp)
}

// GetProjectArtifactWithResponse request returning *GetProjectArtifactResponse
func (c *ClientWithResponses) GetProjectArtifactWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*GetProjectArtifactResponse, error) {
	rsp, err := c.GetProjectArtifact(ctx, projectName, artifactId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProjectArtifactResponse(rsp)
}

// ProjectArtifactExistsWithResponse request returning *ProjectArtifactExistsResponse
func (c *ClientWithResponses) ProjectArtifactExistsWithResponse(ctx context.Context, projectName string, artifactId string, reqEditors ...RequestEditorFn) (*ProjectArtifactExistsResponse, error) {
	rsp, err := c.ProjectArtifactExists(ctx, projectName, artifactId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseProjectArtifactExistsResponse(rsp)
}

// GetProjectArtifactsWithResponse request returning *GetProjectArtifactsResponse
func (c *ClientWithResponses) GetProjectArtifactsWithResponse(ctx context.Context, projectName string, reqEditors ...RequestEditorFn) (*GetProjectArtifactsResponse, error) {
	rsp, err := c.GetProjectArtifacts(ctx, projectName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProjectArtifactsResponse(rsp)
}

// UploadArtifactWithBodyWithResponse request with arbitrary body returning *UploadArtifactResponse
func (c *ClientWithResponses) UploadArtifactWithBodyWithResponse(ctx context.Context, projectName string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UploadArtifactResponse, error) {
	rsp, err := c.UploadArtifactWithBody(ctx, projectName, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUploadArtifactResponse(rsp)
}

// GetProjectsWithResponse request returning *GetProjectsResponse
func (c *ClientWithResponses) GetProjectsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetProjectsResponse, error) {
	rsp, err := c.GetProjects(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProjectsResponse(rsp)
}

// CreateProjectWithBodyWithResponse request with arbitrary body returning *CreateProjectResponse
func (c *ClientWithResponses) CreateProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error) {
	rsp, err := c.CreateProjectWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateProjectResponse(rsp)
}

func (c *ClientWithResponses) CreateProjectWithResponse(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error) {
	rsp, err := c.CreateProject(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateProjectResponse(rsp)
}

// ParseGetHealthResponse parses an HTTP response from a GetHealthWithResponse call
func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Success
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseDeleteProjectResponse parses an HTTP response from a DeleteProjectWithResponse call
func ParseDeleteProjectResponse(rsp *http.Response) (*DeleteProjectResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeleteProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetProjectResponse parses an HTTP response from a GetProjectWithResponse call
func ParseGetProjectResponse(rsp *http.Response) (*GetProjectResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ExtendedProject
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseProjectExistsResponse parses an HTTP response from a ProjectExistsWithResponse call
func ParseProjectExistsResponse(rsp *http.Response) (*ProjectExistsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ProjectExistsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseDeleteProjectArtifactResponse parses an HTTP response from a DeleteProjectArtifactWithResponse call
func ParseDeleteProjectArtifactResponse(rsp *http.Response) (*DeleteProjectArtifactResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeleteProjectArtifactResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetProjectArtifactResponse parses an HTTP response from a GetProjectArtifactWithResponse call
func ParseGetProjectArtifactResponse(rsp *http.Response) (*GetProjectArtifactResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetProjectArtifactResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Artifact
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseProjectArtifactExistsResponse parses an HTTP response from a ProjectArtifactExistsWithResponse call
func ParseProjectArtifactExistsResponse(rsp *http.Response) (*ProjectArtifactExistsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ProjectArtifactExistsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetProjectArtifactsResponse parses an HTTP response from a GetProjectArtifactsWithResponse call
func ParseGetProjectArtifactsResponse(rsp *http.Response) (*GetProjectArtifactsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetProjectArtifactsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ArtifactIds
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUploadArtifactResponse parses an HTTP response from a UploadArtifactWithResponse call
func ParseUploadArtifactResponse(rsp *http.Response) (*UploadArtifactResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &UploadArtifactResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetProjectsResponse parses an HTTP response from a GetProjectsWithResponse call
func ParseGetProjectsResponse(rsp *http.Response) (*GetProjectsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetProjectsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Project
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateProjectResponse parses an HTTP response from a CreateProjectWithResponse call
func ParseCreateProjectResponse(rsp *http.Response) (*CreateProjectResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CreateProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ExtendedProject
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

