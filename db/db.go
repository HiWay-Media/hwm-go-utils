package db

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Open connects to MySQL and applies the pool limits. dbIdleConn and dbMaxConn
// are ignored when empty or not positive (database/sql defaults apply).
func Open(dbUsername string, dbPassword string, dbHost string, dbPort int, dbName string, dbIdleConn, dbMaxConn string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// never include dsn: it contains the password
		return nil, fmt.Errorf("connect to %s@%s:%d/%s: %w", dbUsername, dbHost, dbPort, dbName, err)
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	if idle, err := strconv.Atoi(dbIdleConn); err == nil && idle > 0 {
		sqlDB.SetMaxIdleConns(idle)
	}
	if max, err := strconv.Atoi(dbMaxConn); err == nil && max > 0 {
		sqlDB.SetMaxOpenConns(max)
	}
	return conn, nil
}

// InitDB is Open for callers that want the process to stop on failure.
// Prefer Open in new code.
func InitDB(log *zap.SugaredLogger, dbUsername string, dbPassword string, dbHost string, dbPort int, dbName string, dbIdleConn, dbMaxConn string) *gorm.DB {
	conn, err := Open(dbUsername, dbPassword, dbHost, dbPort, dbName, dbIdleConn, dbMaxConn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	return conn
}
