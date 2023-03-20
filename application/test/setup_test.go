package test

import (
	"context"
	"log"

	"github.com/benchkram/bobc/application"
	"github.com/benchkram/bobc/pkg/artifactstore"
	"github.com/benchkram/bobc/pkg/db"
	"github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/errz"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func setup() (application.Application, error) {
	db := db.New(
		db.WithPostgres(
			databasePostgres.Config().Host,
			databasePostgres.Config().Port,
			databasePostgres.Config().User,
			databasePostgres.Config().Password,
			databasePostgres.Config().Name,
			false,
		),
	)
	err := db.Connect()
	errz.Fatal(err)

	minioConfig := minioInstance.Config()
	minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
		Secure: minioConfig.UseSSL,
	})
	errz.Fatal(err)

	// Make a new bucket.
	bucketName := "bob-server-test-bucket"
	location := "us-east-1"

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			errz.Fatal(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	artifactStore := artifactstore.New(
		minioClient,
		artifactstore.WithBucketName(bucketName),
	)

	projectRepo := projectrepo.New(db, artifactStore)

	return application.New(
		application.WithProjectRepository(projectRepo),
	), nil
}
