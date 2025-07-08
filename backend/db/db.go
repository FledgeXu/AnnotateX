package db

import (
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // 注册 "pgx"
	"github.com/jmoiron/sqlx"
)

var (
	dbInstance *sqlx.DB
	once       sync.Once
	initErr    error
)

func InitDB(dsn string) *sqlx.DB {
	once.Do(func() {
		dbInstance, initErr = sqlx.Connect("pgx", dsn)
		if initErr != nil {
			log.Fatalf("connect db failed: %v", initErr)
		}
		dbInstance.SetMaxOpenConns(10)
		dbInstance.SetMaxIdleConns(5)
		dbInstance.SetConnMaxIdleTime(5 * time.Minute)
		dbInstance.SetConnMaxLifetime(30 * time.Minute)
	})
	return dbInstance
}
