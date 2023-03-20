package db

type DatabaseType string

func (dbt *DatabaseType) String() string {
	return string(*dbt)
}

const (
	Postgres             DatabaseType = "postgres"
	DigitalOceanPostgres DatabaseType = "digital_ocean_postgres"
)
