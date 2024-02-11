package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ngfenglong/food-randomizer-BE/pkg/config"
)

func OpenDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
