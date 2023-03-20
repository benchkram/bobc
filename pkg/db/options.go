package db

type Option func(*database)

func WithDigitalOceanPostgres(token, clusterID string, privateDBConn bool) Option {
	return func(db *database) {
		db.token = token
		db.clusterID = clusterID
		db.privateDBConnection = privateDBConn

		db.dbType = DigitalOceanPostgres
	}
}

func WithPostgres(host, port, user, password, databaseName string, useSSL bool) Option {
	return func(db *database) {
		db.host = host
		db.port = port
		db.user = user
		db.password = password
		db.databaseName = databaseName
		db.useSSL = useSSL

		db.dbType = Postgres
	}
}

func WithConfig(config *Config) Option {
	return func(db *database) {
		db.host = config.Host
		db.port = config.Port
		db.user = config.User
		db.password = config.Password
		db.databaseName = config.Name
		db.useSSL = config.UseSSL

		db.dbType = Postgres
	}
}
