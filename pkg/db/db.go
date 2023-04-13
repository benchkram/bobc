package db

import (
	"errors"
	"log"

	"github.com/benchkram/errz"
	"gorm.io/gorm"
)

// Default errors
var (
	ErrDatabaseNil         = errors.New("Database is nil")
	ErrInvalidDatabaseType = errors.New("Invalid database type")
	ErrStateNotFound       = errors.New("state not found")
	ErrDuplicateUserMail   = errors.New("User with this email already in database")
)

const (
	host = "localhost"
	port = "5432"

	user     = "postgres"
	password = "postgres"

	name        = "postgres"
	nameTesting = "bob_testdb"

	connectStr = "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
)

type Database interface {
	Gorm() *gorm.DB
	Connect() error
}

type database struct {
	gorm *gorm.DB

	// dbType is the database type
	dbType DatabaseType

	// digital ocean postgres
	token               string
	clusterID           string
	privateDBConnection bool

	// postgres
	host         string
	port         string
	user         string
	password     string
	databaseName string
	useSSL       bool
}

func New(opts ...Option) Database {
	database := &database{}
	for _, opt := range opts {
		opt(database)
	}
	return database
}

func (db *database) Connect() (err error) {
	defer errz.Recover(&err)

	switch db.dbType {
	case DigitalOceanPostgres:
		gormDB, err := db.connectDigitalOceanPostgres()
		errz.Fatal(err)

		db.gorm = gormDB
	case Postgres:
		gormDB, err := db.connectPostgres()
		errz.Fatal(err)

		db.gorm = gormDB
	default:
		return ErrInvalidDatabaseType
	}

	log.Printf("Successfully connected to database [%s]\n", db.dbType)

	err = db.migrate()
	errz.Fatal(err)

	return nil
}

func (db *database) Gorm() *gorm.DB {
	return db.gorm
}
