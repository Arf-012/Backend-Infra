package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbSocket := os.Getenv("DB_SOCKET")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	var dsn string
	var connectionMode string

	// 1. Check if DB_SOCKET is provided (For Linux/Production UNIX socket)
	if dbSocket != "" {
		connectionMode = fmt.Sprintf("unix socket (%s)", dbSocket)
		dsn = fmt.Sprintf(
			"%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"+
				"&timeout=5s&readTimeout=5s&writeTimeout=5s&interpolateParams=true",
			dbUser, dbPass, dbSocket, dbName,
		)
	} else {
		// 2. Fallback to TCP/IP if DB_SOCKET is empty (For Windows/Development)
		if dbHost == "" {
			dbHost = "127.0.0.1"
		}
		if dbPort == "" {
			dbPort = "3306"
		}
		connectionMode = fmt.Sprintf("tcp://%s:%s", dbHost, dbPort)
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"+
				"&timeout=5s&readTimeout=5s&writeTimeout=5s&interpolateParams=true",
			dbUser, dbPass, dbHost, dbPort, dbName,
		)
	}

	// 3. Open connection via GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	DB = db
	log.Printf("Database connected with pooling via %s", connectionMode)
}
