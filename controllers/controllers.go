package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"

	"api-rest/converts"
	"api-rest/models"
	"database/sql"
	"strconv"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Save file
	dst := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	models.CreateTable(db)

	defer db.Close()

	var newImage models.Image

	newImage.CompleteImage(dst, file.Filename, strconv.FormatInt(file.Size/1024/1024, 10))

	_, errDB := models.InsertImageDB(db, &newImage)
	if errDB != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded", "filename": file.Filename})
}

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

	rows, err := db.Query("SELECT id, path, name, size, format FROM images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(&img.Id, &img.Path, &img.ImageName, &img.Size, &img.Format)
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
	sql := `SELECT id,path, name, size, format FROM images WHERE id = ?;`
	row := db.QueryRow(sql, id)

	err = row.Scan(&img.Id, &img.Path, &img.ImageName, &img.Size, &img.Format)

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

	image := make([]string, 5, 5)

	sql := `SELECT id,path,name,size,format FROM images WHERE id = ?;`
	row := db.QueryRow(sql, id)
	row.Scan(&image[0], &image[1], &image[2], &image[3], &image[4])

	fmt.Println(image[1])
	fmt.Println(image[1:2])

	if converts.DeleteImages(image[1:2]) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

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

func ConvertImage(c *gin.Context) {
	file, err := c.FormFile("image")
	format := c.Request.FormValue("format")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file"})
		return
	}

	// Save original file
	inputPath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create output filename with new extension
	outputFilename := file.Filename[:len(file.Filename)-4] + "." + format
	outputPath := "./uploads/" + outputFilename

	convertErr := converts.ConvertImage(format, inputPath, outputPath)

	if convertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": convertErr.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+outputFilename)
	c.File(outputPath)

}
