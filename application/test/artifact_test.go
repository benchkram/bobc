package test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/benchkram/errz"
	"github.com/stretchr/testify/assert"
)

func TestArtifactCreation(t *testing.T) {
	app, err := setup()
	assert.Nil(t, err)

	projectName := rnd.RandStringBytesMaskImprSrc(8)

	project, err := app.ProjectCreate(projectName, "a test project")
	assert.Nil(t, err)

	filepath := filepath.Join(os.TempDir(), "file.test")
	bigBuff := make([]byte, 750)
	err = ioutil.WriteFile(filepath, bigBuff, 0666)
	assert.Nil(t, err)
	defer os.RemoveAll(filepath)

	sha1Hash := rnd.RandSHA1(8)
	err = app.ProjectArtifactCreate(project.ID, sha1Hash, filepath, 750)
	errz.Log(err)
	assert.Nil(t, err)

	artifact, err := app.ProjectArtifact(project.ID, sha1Hash)
	assert.Nil(t, err)

	// get file from s3 storage through presigned link
	resp, err := http.Get(artifact.AccessLink.String())
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 750, len(body))
}

func TestArtifactDeletion(t *testing.T) {
	app, err := setup()
	assert.Nil(t, err)

	projectName := rnd.RandStringBytesMaskImprSrc(8)

	project, err := app.ProjectCreate(projectName, "a test project")
	assert.Nil(t, err)

	filepath := filepath.Join(os.TempDir(), "file.test")
	bigBuff := make([]byte, 750)
	err = ioutil.WriteFile(filepath, bigBuff, 0666)
	assert.Nil(t, err)
	defer os.RemoveAll(filepath)

	sha1Hash := rnd.RandSHA1(8)
	err = app.ProjectArtifactCreate(project.ID, sha1Hash, filepath, 750)
	assert.Nil(t, err)

	artifact, err := app.ProjectArtifact(project.ID, sha1Hash)
	assert.Nil(t, err)

	// get file from s3 storage through presigned link
	resp, err := http.Get(artifact.AccessLink.String())
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 750, len(body))

	err = app.ProjectArtifactDelete(project.ID, sha1Hash)
	assert.Nil(t, err)

	exists, err := app.ProjectArtifactExists(project.ID, sha1Hash)
	assert.Nil(t, err)
	assert.False(t, exists)

	// get file from s3 storage through presigned link should fail with "NoSuchKey" in response
	resp, err = http.Get(artifact.AccessLink.String())
	assert.Nil(t, err)
	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(body), "NoSuchKey"))
}
