package database

type DBType string

const (
	PostgresSQL DBType = "postgres"
	MySQL       DBType = "mysql"
	SQLite3     DBType = "sqlite3"
)
