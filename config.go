package main

import (
	"fmt"

	"github.com/benchkram/bobc/pkg/db"
	"github.com/benchkram/bobc/restserver"
	"github.com/fatih/structs"
	"github.com/logrusorgru/aurora"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"
)

var GlobalConfig *config

func readGlobalConfig() {
	// Priority of configuration options
	// 1: CLI Parameters
	// 2: environment
	// 2: config.yaml
	// 3: defaults
	config, err := readConfig(defaultConfig.AsMap())
	if err != nil {
		panic(err.Error())
	}
	config.Print()

	// Set config object for main package
	GlobalConfig = config
}

// get tje default postgres config from the db package
var postgresConfig = db.NewConfig()

var defaultConfig = &config{
	Hostname: "0.0.0.0",
	Port:     "8100",

	DisablePostgres: false,
	PostgresHost:    postgresConfig.Host,
	PostgresPort:    postgresConfig.Port,
	PostgresUser:    postgresConfig.User,
	PostgresPass:    postgresConfig.Password,
	PostgresDBName:  postgresConfig.Name,
	PostgresUseSSL:  postgresConfig.UseSSL,

	DisableS3:         false,
	S3Endpoint:        "localhost:9000",
	S3AccessKeyID:     "minioadmin",
	S3SecretAccessKey: "minioadmin",
	S3UseSSL:          false,
	S3BucketName:      "artifacts",

	UploadDir: restserver.DefaultUploadDir,

	ApiKey: "",
}

func configInit() {
	// Keep cli parameters in sync with the config struct
	rootCmd.PersistentFlags().String("hostname", "", "hostname to listen to")
	rootCmd.PersistentFlags().String("port", "", "port to listen to")

	rootCmd.PersistentFlags().Bool("disable-pg", false, "disable postgres usage only for landing page view")
	rootCmd.PersistentFlags().String("pg-host", defaultConfig.PostgresHost, "hostname to connect with postgres database")
	rootCmd.PersistentFlags().String("pg-port", defaultConfig.PostgresPort, "port to connect with postgres database")
	rootCmd.PersistentFlags().String("pg-user", defaultConfig.PostgresUser, "username for the postgres database")
	rootCmd.PersistentFlags().String("pg-pass", defaultConfig.PostgresPass, "password for the postgres database")
	rootCmd.PersistentFlags().String("pg-db-name", defaultConfig.PostgresDBName, "database name on postgres database")
	rootCmd.PersistentFlags().Bool("pg-use-ssl", defaultConfig.PostgresUseSSL, "whether to use SSL when connecting to postgres")

	rootCmd.PersistentFlags().Bool("disable-s3", defaultConfig.DisableS3, "disable s3 client")
	rootCmd.PersistentFlags().String("s3-endpoint", defaultConfig.S3Endpoint, "s3 endpoint")
	rootCmd.PersistentFlags().String("s3-access-key-id", defaultConfig.S3AccessKeyID, "s3 access key id")
	rootCmd.PersistentFlags().String("s3-secret-access-key", defaultConfig.S3SecretAccessKey, "s3 secret access key")
	rootCmd.PersistentFlags().Bool("s3-use-ssl", defaultConfig.S3UseSSL, "s3 use ssl")
	rootCmd.PersistentFlags().String("s3-bucket-name", defaultConfig.S3BucketName, "s3 bucket name")

	rootCmd.PersistentFlags().String("upload-dir", defaultConfig.UploadDir, "Upload directory on system to upload hash files")

	rootCmd.PersistentFlags().String("api-key", defaultConfig.ApiKey, "API key to check against when authenticating against the http server")

	// CLI PARAMETERS
	_ = viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	_ = viper.BindPFlag("disable-pg", rootCmd.PersistentFlags().Lookup("disable-pg"))
	_ = viper.BindPFlag("pg-host", rootCmd.PersistentFlags().Lookup("pg-host"))
	_ = viper.BindPFlag("pg-port", rootCmd.PersistentFlags().Lookup("pg-port"))
	_ = viper.BindPFlag("pg-user", rootCmd.PersistentFlags().Lookup("pg-user"))
	_ = viper.BindPFlag("pg-pass", rootCmd.PersistentFlags().Lookup("pg-pass"))
	_ = viper.BindPFlag("pg-db-name", rootCmd.PersistentFlags().Lookup("pg-db-name"))
	_ = viper.BindPFlag("pg-use-ssl", rootCmd.PersistentFlags().Lookup("pg-use-ssl"))

	_ = viper.BindPFlag("disable-s3", rootCmd.PersistentFlags().Lookup("disable-s3"))
	_ = viper.BindPFlag("s3-endpoint", rootCmd.PersistentFlags().Lookup("s3-endpoint"))
	_ = viper.BindPFlag("s3-access-key-id", rootCmd.PersistentFlags().Lookup("s3-access-key-id"))
	_ = viper.BindPFlag("s3-secret-access-key", rootCmd.PersistentFlags().Lookup("s3-secret-access-key"))
	_ = viper.BindPFlag("s3-use-ssl", rootCmd.PersistentFlags().Lookup("s3-use-ssl"))
	_ = viper.BindPFlag("s3-bucket-name", rootCmd.PersistentFlags().Lookup("s3-bucket-name"))

	_ = viper.BindPFlag("keto-read-endpoint", rootCmd.PersistentFlags().Lookup("keto-read-endpoint"))
	_ = viper.BindPFlag("keto-write-endpoint", rootCmd.PersistentFlags().Lookup("keto-write-endpoint"))
	_ = viper.BindPFlag("keto-default-namespace", rootCmd.PersistentFlags().Lookup("keto-default-namespace"))

	_ = viper.BindPFlag("upload-dir", rootCmd.PersistentFlags().Lookup("upload-dir"))

	_ = viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))

	// ENVIRONMENT VARS
	_ = viper.BindEnv("hostname", "HOSTNAME")
	_ = viper.BindEnv("port", "PORT")

	_ = viper.BindEnv("disable-pg", "DISABLE_POSTGRES")
	_ = viper.BindEnv("pg-host", "POSTGRES_HOST")
	_ = viper.BindEnv("pg-port", "POSTGRES_PORT")
	_ = viper.BindEnv("pg-user", "POSTGRES_USER")
	_ = viper.BindEnv("pg-pass", "POSTGRES_PASSWORD")
	_ = viper.BindEnv("pg-db-name", "POSTGRES_DB_NAME")
	_ = viper.BindEnv("pg-use-ssl", "POSTGRES_USE_SSL")

	_ = viper.BindEnv("disable-s3", "DISABLE_S3")
	_ = viper.BindEnv("s3-endpoint", "S3_ENDPOINT")
	_ = viper.BindEnv("s3-access-key-id", "S3_ACCESS_KEY_ID")
	_ = viper.BindEnv("s3-secret-access-key", "S3_SECRET_ACCESS_KEY")
	_ = viper.BindEnv("s3-use-ssl", "S3_USE_SSL")
	_ = viper.BindEnv("s3-bucket-name", "S3_BUCKET_NAME")

	_ = viper.BindEnv("keto-read-endpoint", "KETO_READ_ENDPOINT")
	_ = viper.BindEnv("keto-write-endpoint", "KETO_WRITE_ENDPOINT")
	_ = viper.BindEnv("keto-default-namespace", "KETO_DEFAULT_NAMESPACE")

	_ = viper.BindEnv("upload-dir", "UPLOAD_DIRECTORY")

	_ = viper.BindEnv("api-key", "API_KEY")
}

// Create private data struct to hold config options.
// `mapstructure` => viper tags
// `struct` => fatih structs tag
type config struct {
	// REST Server
	Hostname string `mapstructure:"hostname" structs:"hostname"`
	Port     string `mapstructure:"port" structs:"port"`

	// Postgres
	DisablePostgres bool   `mapstructure:"disable-pg" structs:"disable-pg"`
	PostgresHost    string `mapstructure:"pg-host" structs:"pg-host"`
	PostgresPort    string `mapstructure:"pg-port" structs:"pg-port"`
	PostgresUser    string `mapstructure:"pg-user" structs:"pg-user"`
	PostgresPass    string `mapstructure:"pg-pass" structs:"pg-pass"`
	PostgresDBName  string `mapstructure:"pg-db-name" structs:"pg-db-name"`
	PostgresUseSSL  bool   `mapstructure:"pg-use-ssl" structs:"pg-use-ssl"`

	// Object store
	DisableS3         bool   `mapstructure:"disable-s3" structs:"disable-s3"`
	S3Endpoint        string `mapstructure:"s3-endpoint" structs:"s3-endpoint"`
	S3AccessKeyID     string `mapstructure:"s3-access-key-id" structs:"s3-access-key-id"`
	S3SecretAccessKey string `mapstructure:"s3-secret-access-key" structs:"s3-secret-access-key"`
	S3UseSSL          bool   `mapstructure:"s3-use-ssl" structs:"s3-use-ssl"`
	S3BucketName      string `mapstructure:"s3-bucket-name" structs:"s3-bucket-name"`

	// Upload
	UploadDir string `mapstructure:"upload-dir" structs:"upload-dir"`

	// authentication
	ApiKey string `mapstructure:"api-key" structs:"api-key"`
}

func (c *config) AsMap() map[string]interface{} {
	return structs.Map(c)
}

func (c *config) Print() {
	litter.Dump(c)
}

// readConfig a helper to read default from a default config object.
func readConfig(defaults map[string]interface{}) (*config, error) {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	// Read config from file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Printf("%s\n", aurora.Yellow("Could not find a config file"))
		default:
			return nil, fmt.Errorf("config file invalid: %s \n", err)
		}
	}

	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}
	return c, nil

}
