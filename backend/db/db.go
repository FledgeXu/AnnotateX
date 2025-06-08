package db

import (
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDB(dsn string) *sqlx.DB {
	var err error
	DB, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxIdleTime(5 * time.Minute)
	DB.SetConnMaxLifetime(30 * time.Minute)
	return DB
}
