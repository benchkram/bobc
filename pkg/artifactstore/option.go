package artifactstore

import "github.com/minio/minio-go/v7"

type Option func(r *Repository)

func WithMinioClient(c *minio.Client) Option {
	return func(r *Repository) {
		r.minio = c
	}
}

func WithBucketName(bn string) Option {
	return func(r *Repository) {
		r.bucketName = bn
	}
}
