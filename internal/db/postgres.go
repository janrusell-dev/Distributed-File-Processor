package db

import "database/sql"

func NewPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}
