package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	"api-rest/images"

	"github.com/gin-gonic/gin"
)

// images := make([]image,5)

func main() {
	fmt.Println(len(images.Images))
	fmt.Println(cap(images.Images))

	router := gin.Default()
	router.GET("/images", getImages)
	router.GET("/get-image", getImageByUser)
	router.GET("/image/:name", getImageByUser)
	router.POST("/images", postImages)
	router.PUT("/images", UpdateImage)

	router.Run("localhost:8080")

}

func postImages(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}

	CreateTable(db)

	defer db.Close()

	var newImage images.Image

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newImage); err != nil {
		return
	}

	_, err = InsertImageDB(db, &newImage)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
	}

	c.IndentedJSON(http.StatusCreated, newImage)
	// Add the new album to the slice.
	// images.Images = append(images.Images, newImage)
	// c.IndentedJSON(http.StatusCreated, newImage)
}

func getImages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, images.Images)
}

// getImageByUser locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getImageByUser(c *gin.Context) {
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var img images.Image
	sql := `SELECT id, user, name, size, format FROM images WHERE user = ?;`
	row := db.QueryRow(sql, req["username"])

	err = row.Scan(&img.Id, &img.Username, &img.ImageName, &img.Size, &img.Format)
	if err != nil {
		if err == row.Err() {
			c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, img)
}

func UpdateImage(c *gin.Context) {
	var req map[string]string

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	_, err = UpdateImageDB(db, req["Username"], req["Size"], req["Name"], req["Format"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully"})
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS images (
        id INTEGER PRIMARY KEY,
        user TEXT UNIQUE NOT NULL,
        name     TEXT NOT NULL,
        size INTEGER NOT NULL,
        format INTEGER NOT NULL
    );`

	return db.Exec(sql)
}

func InsertImageDB(db *sql.DB, c *images.Image) (int64, error) {
	sql := `INSERT INTO images (user, name, size, format) 
            VALUES (?,?, ?, ?);`
	result, err := db.Exec(sql, c.Username, c.ImageName, c.Size, c.Format)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateImageDB(db *sql.DB, username, size, name, format string) (int64, error) {
	sql := `UPDATE images SET name = ?, size = ?, format = ? WHERE username = ?;`
	result, err := db.Exec(sql, name, size, format, username)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
