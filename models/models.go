package models

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type Image struct {
	Id, Path, ImageName, Size, Format string
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS images (
        id INTEGER PRIMARY KEY,
        name     TEXT NOT NULL,
        size INTEGER NOT NULL,
        format INTEGER NOT NULL
    );`

	return db.Exec(sql)
}

func InsertImageDB(db *sql.DB, c *Image) (int64, error) {
	sql := `INSERT INTO images (path, name, size, format) 
            VALUES (?, ?, ?);`
	result, err := db.Exec(sql, c.Path, c.ImageName, c.Size, c.Format)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateImageDB(db *sql.DB, id int, size, name, format string) (int64, error) {
	sql := `UPDATE images SET name = ?, size = ?, format = ? WHERE id = ?;`
	result, err := db.Exec(sql, name, size, format, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteImageDB(db *sql.DB, id int) (int64, error) {
	sql := `DELETE FROM images WHERE id = ?`
	result, err := db.Exec(sql, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
