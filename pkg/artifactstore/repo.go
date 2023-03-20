package artifactstore

import (
	"github.com/minio/minio-go/v7"
)

type Repository struct {
	// minio is used to connect to s3 compatible storage
	minio *minio.Client

	// bucketName used to store artifacts
	bucketName string
}

func New(minioClient *minio.Client, opts ...Option) *Repository {
	r := &Repository{
		minio:      minioClient,
		bucketName: "artifacts",
	}

	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}

	return r
}
