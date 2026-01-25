package config

type DBConfig struct {
	ConnectionString string `yaml:"connectionString"`
	// TODO: Add configs for connection pool, timeouts, retries, etc.
}
