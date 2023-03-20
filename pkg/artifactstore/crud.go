package artifactstore

import (
	"context"
	"net/url"
	"time"

	"github.com/benchkram/errz"
	"github.com/minio/minio-go/v7"
)

func (r *Repository) CreateArtifact(id string, filePath string, size int) (err error) {
	defer errz.Recover(&err)

	_, err = r.minio.FPutObject(context.Background(),
		r.bucketName,
		id,
		filePath,
		minio.PutObjectOptions{
			ContentType: "application/tar+gzip",
		},
	)
	errz.Fatal(err)

	return nil
}

func (r *Repository) DeleteArtifact(id string) (err error) {
	defer errz.Recover(&err)

	err = r.minio.RemoveObject(
		context.Background(),
		r.bucketName,
		id,
		minio.RemoveObjectOptions{},
	)
	errz.Fatal(err)

	return nil
}

func (r *Repository) Artifact(id string) (addr *url.URL, err error) {
	defer errz.Recover(&err)

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+id+"\"")

	addr, err = r.minio.PresignedGetObject(
		context.Background(),
		r.bucketName,
		id,
		10*time.Minute,
		reqParams,
	)
	errz.Fatal(err)

	return addr, nil
}
