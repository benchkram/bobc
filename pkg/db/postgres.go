package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/benchkram/errz"
	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (db *database) connectPostgres() (gormDB *gorm.DB, err error) {
	err = db.prepareDatabase()
	errz.Fatal(err)

	return db.connectGorm()
}

func (db *database) connectString() string {
	sslMode := "disable"
	if db.useSSL {
		sslMode = "require"
	}

	const connectStr = "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"
	return fmt.Sprintf(connectStr,
		db.host,
		db.port,
		db.user,
		db.databaseName,
		db.password,
		sslMode,
	)
}

func (db *database) connectGorm() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(db.connectString()), &gorm.Config{})
}

// prepareDatabase for first time use.
func (db *database) prepareDatabase() error {
	err := db.withRawDB(func(conn *pgx.Conn) error {
		createDatabase := fmt.Sprintf("CREATE DATABASE %s;", db.databaseName)
		return conn.PgConn().Exec(context.Background(), createDatabase).Close()
	})

	// TODO: Remove string search
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 42P04") {
		// Ignore "SQLSTATE 42P04" error which is returned when the DB already exists.
		err = nil
	}
	errz.Fatal(err)

	return nil
}

func (db *database) withRawDB(fn func(*pgx.Conn) error) (err error) {
	defer errz.Recover(&err)

	var conn *pgx.Conn
	attempt := 1
	err = retry.Do(
		func() error {
			log.Printf("Trying to connect to database: [attempt: %d]\n", attempt)
			attempt++

			var err error
			conn, err = pgx.Connect(context.Background(), db.connectString())
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(20),
		retry.Delay(5*time.Second),
	)
	errz.Fatal(err)

	defer conn.Close(context.Background())
	return fn(conn)
}
