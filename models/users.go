package models

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type User struct {
	Id, UserName, PasswordHash, RefreshToken, Role string
}

type UpdateUserRequest struct {
	UserName      string `json:"username" binding:"required"`
	Role          string `json:"role" binding:"required"`
	Password      string `json:"password" binding:"required"`
	AdminPassword string `json:"admin_password"`
}

func (u *User) FillAttributes(id, userName, passwordHash, refreshToken, role string) {
	u.Id = id
	u.UserName = userName
	u.PasswordHash = passwordHash
	u.RefreshToken = refreshToken
	u.Role = role
}

func InsertUserDB(db *sql.DB, u *User) (int64, error) {
	query := `INSERT INTO users (username, password_hash, refresh_token, role) VALUES (?,?,?,?)`

	result, err := db.Exec(query, u.UserName, u.PasswordHash, u.RefreshToken, u.Role)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateUserDB(db *sql.DB, u *UpdateUserRequest, userId, hash string) (int64, error) {
	query := `UPDATE users SET username = ?, password_hash = ?, role = ? WHERE id = ?;`

	result, err := db.Exec(query, u.UserName, hash, u.Role, userId)

	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteUserDB(db *sql.DB, id string) (int64, error) {
	query := `DELETE FROM users WHERE id = ?;`

	result, err := db.Exec(query, id)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
