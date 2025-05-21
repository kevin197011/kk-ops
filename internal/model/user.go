package model

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (username, password, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	return GetDB().QueryRow(context.Background(), query,
		user.Username, user.Password, user.Email, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, username, password, email, role, last_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	var lastLogin sql.NullTime

	err := GetDB().QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Role,
		&lastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle NULL last_login
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, password, email, role, last_login, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &User{}
	var lastLogin sql.NullTime

	err := GetDB().QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Role,
		&lastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle NULL last_login
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetAllUsers 获取所有用户
func GetAllUsers() ([]*User, error) {
	query := `
		SELECT id, username, email, role, last_login, created_at, updated_at
		FROM users
		ORDER BY id
	`

	rows, err := GetDB().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		var lastLogin sql.NullTime

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role,
			&lastLogin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Handle NULL last_login
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser 更新用户
func UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = $1, email = $2, role = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := GetDB().Exec(context.Background(), query,
		user.Username, user.Email, user.Role, user.UpdatedAt, user.ID)
	return err
}

// UpdatePassword 更新用户密码
func UpdatePassword(userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err = GetDB().Exec(context.Background(), query, string(hashedPassword), time.Now(), userID)
	return err
}

// UpdateLastLogin 更新用户最后登录时间
func UpdateLastLogin(userID int) error {
	now := time.Now()
	query := `
		UPDATE users
		SET last_login = $1
		WHERE id = $2
	`

	_, err := GetDB().Exec(context.Background(), query, now, userID)
	return err
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := GetDB().Exec(context.Background(), query, id)
	return err
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
