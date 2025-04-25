package db

import (
	"database/sql"

	"github.com/sajitha-tj/go-sts/internal/configs"
)

func New(conf *configs.Database) (*sql.DB, error) {
	db, err := sql.Open("postgres", "user="+conf.Username+" password="+conf.Password+" dbname="+conf.Name+" sslmode="+conf.SSLMode)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}