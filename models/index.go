package models

import "database/sql"

var db *sql.DB

func New(dbPool *sql.DB) {
	db = dbPool
}
