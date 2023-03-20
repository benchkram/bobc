package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	"github.com/benchkram/errz"
	"github.com/digitalocean/godo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (db *database) connectDigitalOceanPostgres() (gormDB *gorm.DB, err error) {
	defer errz.Recover(&err)

	conn, err := getDatabaseConnection(db.token, db.clusterID, "", db.privateDBConnection)
	errz.Log(err)
	errz.Fatal(err)

	return connectGorm(conn)
}

func getDatabaseConnection(token, clusterID, databaseID string, privateDBConn bool) (conn *godo.DatabaseConnection, err error) {
	defer errz.Recover(&err)

	fmt.Println("get database connection")

	var cluster *godo.Database
	attempt := 1
	err = retry.Do(
		func() error {
			log.Printf("Trying to get database connection: [attempt: %d]\n", attempt)
			attempt++

			client := godo.NewFromToken(token)
			ctx := context.Background()

			var err error
			//cluster, _, err = client.Databases.GetPool(ctx, clusterID, "conpool")
			cluster, _, err = client.Databases.Get(ctx, clusterID)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(100),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.FixedDelay),
	)
	errz.Fatal(err)

	// litter.Dump(cluster)
	fmt.Printf("using private connection: %v\n", privateDBConn)

	if privateDBConn {
		return cluster.PrivateConnection, nil
	}

	return cluster.Connection, nil
}

func connectGorm(conn *godo.DatabaseConnection) (db *gorm.DB, err error) {
	defer errz.Recover(&err)

	fmt.Println("get gorm db")

	SSLMode := "disable"

	// TODO: enable ssl
	if conn.SSL {
		SSLMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", conn.Host, conn.User, conn.Password, conn.Database, conn.Port, SSLMode)
	fmt.Printf("dsn: %s\n", dsn)
	fmt.Println("try to open gorm")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	errz.Fatal(err)

	return db, nil
}
