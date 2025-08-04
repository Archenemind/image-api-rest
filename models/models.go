package models

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type Image struct {
	Id, Path, ImageName, Size, Format string
}

// Fills all the atributes of models.Image
func (i *Image) CompleteImage(path, imageName, size string) error {
	var err error

	i.Path = path
	i.ImageName = imageName
	i.Size = size + " MB"
	i.Format, err = FindFormat(imageName)

	if err != nil {
		return err
	}
	return nil
}

func FindFormat(imageName string) (string, error) {
	for i, v := range imageName {
		if v == '.' {
			return imageName[i+1:], nil
		}
	}
	return "error", fmt.Errorf("no dot in the file name")
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS images (
        id INTEGER PRIMARY KEY,
		path TEXT NOT NULL,
        name     TEXT NOT NULL,
        size INTEGER NOT NULL,
        format INTEGER NOT NULL
    );`

	return db.Exec(sql)
}

func InsertImageDB(db *sql.DB, c *Image) (int64, error) {
	sql := `INSERT INTO images (path, name, size, format) 
            VALUES (?, ?, ?, ?);`
	result, err := db.Exec(sql, c.Path, c.ImageName, c.Size, c.Format)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateImageDB(db *sql.DB, id int, path, size, name, format string) (int64, error) {
	sql := `UPDATE images SET path = ?,name = ?, size = ?, format = ? WHERE id = ?;`
	result, err := db.Exec(sql, path, name, size, format, id)
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
