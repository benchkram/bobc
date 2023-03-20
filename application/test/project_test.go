package test

import (
	"testing"

	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/stretchr/testify/assert"
)

func TestProjectCreation(t *testing.T) {
	app, err := setup()
	assert.Nil(t, err)

	projectName := rnd.RandStringBytesMaskImprSrc(8)

	project, err := app.ProjectCreate(projectName, "")
	assert.Nil(t, err)

	exists, err := app.ProjectExists(projectName)
	assert.Nil(t, err)
	assert.True(t, exists)

	_, err = app.Project(project.ID)
	assert.Nil(t, err)
}
