package utils

import (
	"net/http"
	"strings"
	"time"

	"os"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var JwtSecret = []byte(GetJwtSecret())
var RefreshSecret = []byte(GetRefreshSecret()) // Should match the one in controllers
var AdminSecret = []byte(GetAdminSecret())

func GetJwtSecret() (JwtSecret string) {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	return jwtSecret
}
func GetRefreshSecret() (RefreshSecret string) {
	godotenv.Load()
	refreshSecret := os.Getenv("REFRESH_SECRET")

	return refreshSecret
}
func GetAdminSecret() string {
	godotenv.Load()
	adminSecret := os.Getenv("ADMIN_REGISTER_PASSWORD")

	return adminSecret
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
		}

		c.Next()

	}
}

func GenerateJWT(userID int64, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func GenerateRefreshToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString(RefreshSecret)
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
