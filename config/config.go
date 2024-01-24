package config

// The configuration values that the service uses to initialize and start.
type ServiceConfig struct {
	// The configuration of the external data store
	Datastore struct {
		// postgres-specific configuration details
		PostgreSQL PostgresConfig // PostgresConfig holds the data needed to create a connection to a postgres database
	}
}

type PostgresConfig struct {
	Username string `envconfig:"POSTGRES_USERNAME"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	Host     string `envconfig:"POSTGRES_HOST"`
	Database string `envconfig:"POSTGRES_DATABASE"`
}
