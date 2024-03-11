package database

import (
	"cobaMetrics/app/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func ConnectDB(config config.IConfig) *sql.DB {
	dbConfig := config.Config().Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("cant connect to db : %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	log.Info("success connect database")

	return db
}
