package test

import (
	"os"
	"testing"

	"github.com/benchkram/bobc/test/environment"
	"github.com/benchkram/errz"
	_ "github.com/lib/pq"
)

var databasePostgres *environment.Postgres
var minioInstance *environment.MinIO

func TestMain(m *testing.M) {
	var err error

	// setup database
	databasePostgres, err = environment.NewPostgresStarted()
	if err != nil {
		errz.Log(err)
		os.Exit(1)
	}

	minioInstance, err = environment.NewMinioStarted()
	if err != nil {
		errz.Log(err)
		os.Exit(1)
	}

	code := m.Run()

	// shutdown database
	err = databasePostgres.Stop(true)
	if err != nil {
		if code == 0 {
			code = 1
		}
		errz.Log(err)
	}

	// shutdown minio
	err = minioInstance.Stop(true)
	if err != nil {
		if code == 0 {
			code = 1
		}
		errz.Log(err)
	}

	os.Exit(code)
}
