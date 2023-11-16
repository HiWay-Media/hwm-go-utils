package db


import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(log *zap.SugaredLogger, dbUsername string, dbPassword string, dbHost string, dbPort int, dbName string, dbIdleConn, dbMaxConn string ) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connection on %s, err: %s", dsn, err.Error())
	}
	db, _ := conn.DB()
	idle, _ := strconv.Atoi(dbIdleConn)
	max, _ := strconv.Atoi(dbMaxConn)

	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(max)

	return conn
}
