package models

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type Image struct {
	Id, Path, ImageName, Size, Format, UserId string
}

type UpdateImgReq struct {
	Name   string `json:"name" binding:"required"`
	Format string `json:"format" binding:"required"`
}

// Fills all the atributes of models.Image
func (i *Image) FillAttributes(path, imageName, size, userId string) error {
	var err error

	i.Path = path
	i.ImageName = imageName
	i.Size = size + " MB"
	i.Format, err = FindFormat(imageName)
	i.UserId = userId
	if err != nil {
		return err
	}
	return nil
}

func FindFormat(imageName string) (string, error) {
	index := -1
	for i, v := range imageName {
		if v == '.' {
			index = i
			//return imageName[i+1:], nil
		}
	}

	if index == -1 {
		return "error", fmt.Errorf("no dot in the file name: %s", imageName)
	}

	return imageName[index+1:], nil
}

func InsertImageDB(db *sql.DB, img *Image) (int64, error) {
	query := `INSERT INTO images (path, name, size, format,user_id) 
            VALUES (?, ?, ?, ?,?);`
	result, err := db.Exec(query, img.Path, img.ImageName, img.Size, img.Format, img.UserId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateImageDB(db *sql.DB, id int, path, size, name, format string) (int64, error) {
	query := `UPDATE images SET path = ?,name = ?, size = ?, format = ? WHERE id = ?;`
	result, err := db.Exec(query, path, name, size, format, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteImageDB(db *sql.DB, id int) (int64, error) {
	query := `DELETE FROM images WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
