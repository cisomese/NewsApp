package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	mysqlCfg := mysql.NewConfig()

	dbport, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Print("incorrect port value")
	}
	mysqlCfg.User = os.Getenv("DB_USER")
	mysqlCfg.Passwd = os.Getenv("DB_PASSWORD")
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%d", os.Getenv("DB_HOST"), dbport)
	mysqlCfg.DBName = os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
