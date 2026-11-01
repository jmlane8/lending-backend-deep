package config

// EnvPostgreSQL holds environment variables for the PostgreSQL connection pool.
type EnvPostgreSQL struct {
	PostgresHost        string `env:"POSTGRES_HOST"          envDefault:"localhost"`
	PostgresPort        int    `env:"POSTGRES_PORT"          envDefault:"5432"`
	PostgresDatabase    string `env:"POSTGRES_DATABASE"      envDefault:"lending"`
	PostgresSchema      string `env:"POSTGRES_SCHEMA"        envDefault:"public"`
	PostgresUsername    string `env:"POSTGRES_USERNAME"      envDefault:"postgres"`
	PostgresPassword    string `env:"POSTGRES_PASSWORD,unset"`
	PostgresSSLMode     string `env:"POSTGRES_SSL_MODE"      envDefault:"disable"`
	PostgresConnMaxOpen int    `env:"POSTGRES_CONN_MAX_OPEN" envDefault:"25"`
	PostgresConnMaxIdle int    `env:"POSTGRES_CONN_MAX_IDLE" envDefault:"5"`
	PostgresConnMinIdle int    `env:"POSTGRES_CONN_MIN_IDLE" envDefault:"1"`
}
