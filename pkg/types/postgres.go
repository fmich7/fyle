package types

// PostgresCredentials represents credentials for connecting to Postgres database
type PostgresCredentials struct {
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_HOST     string
	DB_PORT     string
}
