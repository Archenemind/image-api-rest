package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/Archenemind/image-api-rest/models"
	"github.com/Archenemind/image-api-rest/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	JwtSecret     = utils.JwtSecret
	RefreshSecret = utils.RefreshSecret
	AdminSecret   = utils.AdminSecret
)

type RegisterOrLoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	_, err = models.CreateTables(db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating DB"})
		return
	}

	var registerRequest RegisterOrLoginRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	refreshToken, _ := utils.GenerateRefreshToken(registerRequest.UserName)

	user := &models.User{}
	user.FillAttributes("", registerRequest.UserName, string(passwordHash), refreshToken, "client")

	id, err := models.InsertUserDB(db, user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(id, registerRequest.UserName, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "username": registerRequest.UserName,
		"access token": token, "refresh token": refreshToken, "role": user.Role})
}

func LoginUser(c *gin.Context) {
	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var loginReq RegisterOrLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := `SELECT id, password_hash, username,refresh_token FROM users WHERE username = ?;`
	err = db.QueryRow(query, loginReq.UserName).Scan(&user.Id, &user.PasswordHash,
		&user.UserName, &user.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	id, _ := strconv.Atoi(user.Id)
	token, err := utils.GenerateJWT(int64(id), user.UserName, user.Role) // You'll need to convert user.Id to int64
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access token": token, "username": user.UserName,
		"refresh token": user.RefreshToken})
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	var req models.UpdateUserRequest

	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if req.Role != "client" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "role is incorrect"})
		return
	}

	db, err := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	var user models.User
	row := db.QueryRow(`SELECT id,username,password_hash,role FROM users WHERE id =?;`, userId)

	err = row.Scan(&user.Id, &user.UserName, &user.PasswordHash, &user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if req.Role == "admin" && req.AdminPassword != string(AdminSecret) {
		c.JSON(http.StatusUnauthorized, "wrong admin password")
	}

	models.UpdateUserDB(db, &req, userId, string(hash))

	c.JSON(http.StatusOK, gin.H{"username": req.UserName, "role": req.Role})
}

func RefreshToken(c *gin.Context) {

	db, _ := sql.Open("sqlite", "./images.db?_pragma=foreign_keys(1)")
	defer db.Close()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.RefreshSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	var user models.User

	query := `SELECT id, username,
	refresh_token FROM users WHERE username = ?;`
	db.QueryRow(query, username).Scan(&user.Id,
		&user.UserName, &user.RefreshToken)

	if user.RefreshToken != tokenString {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token mismatch"})
		return
	}

	userId, _ := strconv.Atoi(user.Id)

	newToken, _ := utils.GenerateJWT(int64(userId), username, user.Role)
	c.JSON(http.StatusOK, gin.H{"access token": newToken})
}
