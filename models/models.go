package models

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

func CreateTables(db *sql.DB) (sql.Result, error) {

	sql := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
			password_hash TEXT NOT NULL,
			username VARCHAR(255) NOT NULL UNIQUE, 
			role VARCHAR(10) NOT NULL,
			refresh_token TEXT NOT NULL
			);
		CREATE TABLE IF NOT EXISTS images (
        id INTEGER PRIMARY KEY,
		path TEXT NOT NULL,
        name     TEXT NOT NULL,
        size INTEGER NOT NULL,
        format INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
    );`

	return db.Exec(sql)
}
