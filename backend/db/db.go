package db

import (
	"annotate-x/models"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	dbInstance *sqlx.DB
	once       sync.Once
	initErr    error
)

func InitDB(dsn models.DataSourceName) *sqlx.DB {
	once.Do(func() {
		dbInstance, initErr = sqlx.Connect("pgx", string(dsn))
		if initErr != nil {
			panic(fmt.Sprintf("connect db failed: %v", initErr))
		}
		dbInstance.SetMaxOpenConns(10)
		dbInstance.SetMaxIdleConns(5)
		dbInstance.SetConnMaxIdleTime(5 * time.Minute)
		dbInstance.SetConnMaxLifetime(30 * time.Minute)
	})
	return dbInstance
}
