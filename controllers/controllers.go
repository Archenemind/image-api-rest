package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"

	"database/sql"
	"strconv"

	"github.com/Archenemind/image-api-rest/models"
	"github.com/Archenemind/image-api-rest/utils"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	format, _ := models.FindFormat(file.Filename)

	if format != "png" && format != "jpg" && format != "webp" &&
		format != "avif" {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "the file is not an allowed image format"})
		return
	}
	utils.CreateDirectoryIfNotExists("/uploads")
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

	_, err = models.CreateTables(db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating DB"})
		return
	}

	defer db.Close()

	var newImage models.Image

	newImage.CompleteImage(dst, file.Filename, strconv.FormatInt(file.Size/1024/1024, 10))

	lastId, errDB := models.InsertImageDB(db, &newImage)
	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDB.Error()})
		return
	}

	dst = dst[0:10] + strconv.FormatInt(lastId, 10) + "&" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.Exec(`UPDATE images SET path = ? WHERE id = ?;`, dst, lastId)

	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded", "filename": file.Filename, "path": dst, "id": lastId})
}

func GetImages(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, path, name, size, format FROM images`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var img models.Image
		err := rows.Scan(&img.Id, &img.Path,
			&img.ImageName, &img.Size, &img.Format)
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
	c.JSON(http.StatusOK, gin.H{
		"base64 image": utils.ConvertImageToBase64(img.Path, img.Format)})
}

// Changes the format of the existing images
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

	var image models.Image
	query := `SELECT path, name, size, format FROM images WHERE id = ?;`
	row := db.QueryRow(query, id)

	if err := row.Scan(&image.Path, &image.ImageName, &image.Size, &image.Format); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
	}

	if req["Format"] != image.Format &&
		image.Format != "webp" && image.Format != "avif" {
		err = utils.ConvertImage(image.Path[len(image.Path)-3:], req["Format"],
			image.Path, image.Path[:len(image.Path)-3]+req["Format"])

		if err != nil {
			image.Path = image.Path[:len(image.Path)-3] + req["Format"]
			utils.DeleteImages([]string{image.Path})

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.DeleteImages([]string{image.Path})
		image.Path = image.Path[:len(image.Path)-3] + req["Format"]

	} else if req["Format"] != image.Format &&
		(image.Format == "webp" || image.Format == "avif") {

		err = utils.ConvertImage(image.Format, req["Format"],
			image.Path, image.Path[:len(image.Path)-4]+req["Format"])

		if err != nil {
			image.Path = image.Path[:len(image.Path)-4] + req["Format"]
			utils.DeleteImages([]string{image.Path})

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		utils.DeleteImages([]string{image.Path})
		image.Path = image.Path[:len(image.Path)-4] + req["Format"]

	}
	newPath := image.Path[:10] + id + "&" + req["Name"] + "." + req["Format"]

	utils.ChangeFileName(image.Path, newPath)

	fileSize := utils.GetFileSize(newPath)

	req["id"] = id
	idDB, _ := strconv.Atoi(id)
	rowsAffected, errDB := models.UpdateImageDB(db, idDB,
		newPath, strconv.FormatFloat(float64(fileSize), 'f', 1, 32)+" MB",
		req["Name"], req["Format"])

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

	var path string

	row := db.QueryRow(`SELECT path FROM images WHERE id = ?;`, id)
	row.Scan(&path)

	if utils.DeleteImages([]string{path}) != nil {
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

	c.JSON(http.StatusOK, gin.H{"message": "image deleted successfully"})
}

func ConvertAndDeleteImage(c *gin.Context) {
	file, err := c.FormFile("image")
	requestedFormat := c.Request.FormValue("format")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file"})
		return
	}

	fileFormat, _ := models.FindFormat(file.Filename)

	if fileFormat != "png" && fileFormat != "jpg" && fileFormat != "webp" &&
		fileFormat != "avif" {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "the file is not an allowed image format"})
		return
	}

	utils.CreateDirectoryIfNotExists("/temporal")
	// Save original file
	inputPath := "./temporal/" + file.Filename
	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create output filename with new extension
	outputFilename := file.Filename[:len(file.Filename)-4] + "." + requestedFormat
	outputPath := "./temporal/" + outputFilename

	convertErr := utils.ConvertImage(file.Filename[len(file.Filename)-3:],
		requestedFormat, inputPath, outputPath)

	if convertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": convertErr.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+outputFilename)
	c.File(outputPath)

	utils.DeleteImages([]string{inputPath, outputPath})
}
