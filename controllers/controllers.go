package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"

	"api-rest/models"
	"database/sql"
	"strconv"
)

func PostImages(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}

	models.CreateTable(db)

	defer db.Close()

	var newImage models.Image

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newImage); err != nil {
		return
	}

	id, errDB := models.InsertImageDB(db, &newImage)

	if errDB != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
	}

	newImage.Id = strconv.Itoa(int(id))
	c.IndentedJSON(http.StatusCreated, newImage)
	// Add the new album to the slice.
	// images.Images = append(images.Images, newImage)
	// c.IndentedJSON(http.StatusCreated, newImage)
}

func GetImages(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, size, format FROM images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(&img.Id, &img.ImageName, &img.Size, &img.Format)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		images = append(images, img)
	}

	c.JSON(http.StatusOK, images)
}

// GetImageById locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func GetImageById(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var img models.Image
	sql := `SELECT id, name, size, format FROM images WHERE id = ?;`
	row := db.QueryRow(sql, id)

	err = row.Scan(&img.Id, &img.ImageName, &img.Size, &img.Format)

	if err != nil {
		if err == row.Err() {
			c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		} else if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, img)
}

func UpdateImage(c *gin.Context) {
	id := c.Param("id")
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

	req["id"] = id
	idDB, _ := strconv.Atoi(id)
	rowsAffected, errDB := models.UpdateImageDB(db, idDB, req["Size"], req["Name"], req["Format"])

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	} else if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully", "result": req})
}

func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	idDB, _ := strconv.Atoi(id)
	rowsAffected, errDB := models.DeleteImageDB(db, idDB)

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	} else if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
