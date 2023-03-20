package environment

import (
	"context"
	"strconv"

	"github.com/benchkram/bobc/application"
	"github.com/benchkram/bobc/pkg/artifactstore"
	"github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/bobc/restserver/authenticator"

	database "github.com/benchkram/bobc/pkg/db"
	"github.com/benchkram/bobc/restserver"
	"github.com/benchkram/errz"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/phayes/freeport"
	"github.com/rs/zerolog/log"
)

// initialize a database connection using the config
// establish connection with the database
func InitDatabase(config *database.Config) (db database.Database, err error) {
	defer errz.Recover(&err)

	db = database.New(
		database.WithConfig(config),
	)

	err = db.Connect()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// initialize the application using the config
// provided by newly created postgres
func InitApp(conf *database.Config, minioConfig MinIOConfig) (_ application.Application, err error) {
	defer errz.Recover(&err)

	db, err := InitDatabase(conf)
	errz.Fatal(err)

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

	app := application.New(
		application.WithProjectRepository(projectRepo),
	)

	return app, nil
}

// initialize the rest server from getting random
// port from the free ports
func InitRestServer(app application.Application, apiKey []byte) (*restserver.S, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	// create rest-server
	authn := authenticator.New(apiKey)

	// create rest-server
	s, err := restserver.New(
		restserver.WithArtifactService(app),
		restserver.WithAuthenticator(authn),
		restserver.WithHost(defaultRestServerHostname, strconv.Itoa(port)),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}
