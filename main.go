package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/benchkram/bobc/application"
	"github.com/benchkram/bobc/pkg/artifactstore"
	"github.com/benchkram/bobc/pkg/projectrepo"
	"github.com/benchkram/bobc/restserver"

	database "github.com/benchkram/bobc/pkg/db"
	"github.com/benchkram/bobc/pkg/wait"
	"github.com/benchkram/bobc/restserver/authenticator"
	"github.com/benchkram/errz"
	"github.com/logrusorgru/aurora"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

func init() {
	configInit()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		errz.Log(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bob-server",
	Short: "cli to start server & client for bob server",
	Long:  `TODO`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readGlobalConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	fmt.Printf("\n  %s\n\n", aurora.Green("Starting bob-server"))

	minioClient, err := minio.New(GlobalConfig.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(GlobalConfig.S3AccessKeyID, GlobalConfig.S3SecretAccessKey, ""),
		Secure: GlobalConfig.S3UseSSL,
	})
	errz.Fatal(err)

	err = minioClient.MakeBucket(context.Background(), GlobalConfig.S3BucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), GlobalConfig.S3BucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own bucket %s\n", GlobalConfig.S3BucketName)
		} else {
			errz.Fatal(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", GlobalConfig.S3BucketName)
	}

	artifactStore := artifactstore.New(
		minioClient,
		artifactstore.WithBucketName(GlobalConfig.S3BucketName),
	)

	// create database
	db := database.New(
		database.WithPostgres(
			GlobalConfig.PostgresHost,
			GlobalConfig.PostgresPort,
			GlobalConfig.PostgresUser,
			GlobalConfig.PostgresPass,
			GlobalConfig.PostgresDBName,
			GlobalConfig.PostgresUseSSL,
		),
	)

	err = db.Connect()
	errz.Fatal(err)

	projectRepo := projectrepo.New(db, artifactStore)

	app := application.New(
		application.WithProjectRepository(projectRepo),
	)

	restOpts := []restserver.Option{
		restserver.WithArtifactService(app),
		restserver.WithHost(GlobalConfig.Hostname, GlobalConfig.Port),
		restserver.WithUploadDir(GlobalConfig.UploadDir),
	}

	restOpts = append(restOpts, restserver.WithHost(GlobalConfig.Hostname, GlobalConfig.Port))
	restOpts = append(restOpts, restserver.WithUploadDir(GlobalConfig.UploadDir))

	// create rest-server
	authn := authenticator.New([]byte(GlobalConfig.ApiKey))
	restOpts = append(restOpts, restserver.WithAuthenticator(authn))

	server, err := restserver.New(
		restOpts...,
	)
	errz.Fatal(err)

	// Run rest-server
	err = server.Start()
	errz.Fatal(err)

	wait.ForCtrlC()

	err = server.Stop()
	errz.Fatal(err)
}
