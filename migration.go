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
