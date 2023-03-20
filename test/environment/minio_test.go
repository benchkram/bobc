package environment

import (
	"testing"

	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

func TestMinioSetup(t *testing.T) {
	minioInstance, err := NewMinioStarted()
	assert.Nil(t, err)

	litter.Dump(minioInstance.Config())

	err = minioInstance.Stop(true)
	assert.Nil(t, err)
}
