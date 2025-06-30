package main

import (
	"database/sql"
	"log"
)

func CreateUserTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT,
			email VARCHAR(255) UNIQUE NOT NULL,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			createdAt datetime,
			updatedAt datetime,
			PRIMARY KEY (id)
		);
	`

	if _, err := db.Query(query); err != nil {
		log.Fatal(err)
	}
}

func CreateTodoTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT,
			user_id INT NOT NULL,
			title VARCHAR(256) NOT NULL,
			content VARCHAR(256) NOT NULL,
			createdAt DATETIME DEFAULT NOW(),
			updatedAt DATETIME DEFAULT NOW(),
			done TINYINT(1) NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		);
	`
	if _, err := db.Query(query); err != nil {
		log.Fatal(err)
	}
}
