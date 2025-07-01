package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	var (
		dbUser = os.Getenv("MYSQL_USERNAME")
		dbPass = os.Getenv("MYSQL_PASSWORD")
		dbHost = os.Getenv("MYSQL_HOST")
		dbPort = os.Getenv("MYSQL_PORT")
		dbName = os.Getenv("MYSQL_DBNAME")
	)

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
