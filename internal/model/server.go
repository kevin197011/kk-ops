package model

import (
	"context"
	"database/sql"
	"time"
)

// Server 服务器模型
type Server struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	SSHKeyPath  string    `json:"ssh_key_path"` // SSH密钥路径
	AuthType    string    `json:"auth_type"`    // 认证类型: password 或 key
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateServer 创建服务器
func CreateServer(server *Server) error {
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	query := `
		INSERT INTO servers (name, ip, port, username, password, description, status, ssh_key_path, auth_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	return GetDB().QueryRow(context.Background(), query,
		server.Name, server.IP, server.Port, server.Username, server.Password,
		server.Description, server.Status, server.SSHKeyPath, server.AuthType, server.CreatedAt, server.UpdatedAt).Scan(&server.ID)
}

// GetServerByID 根据ID获取服务器
func GetServerByID(id int) (*Server, error) {
	query := `
		SELECT id, name, ip, port, username, password, description, status, ssh_key_path, auth_type, created_at, updated_at
		FROM servers
		WHERE id = $1
	`

	server := &Server{}
	var sshKeyPath, authType sql.NullString

	err := GetDB().QueryRow(context.Background(), query, id).Scan(
		&server.ID, &server.Name, &server.IP, &server.Port, &server.Username, &server.Password,
		&server.Description, &server.Status, &sshKeyPath, &authType, &server.CreatedAt, &server.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle NULL values
	if sshKeyPath.Valid {
		server.SSHKeyPath = sshKeyPath.String
	}
	if authType.Valid {
		server.AuthType = authType.String
	} else {
		server.AuthType = "password" // Default to password auth if NULL
	}

	return server, nil
}

// GetAllServers 获取所有服务器
func GetAllServers() ([]*Server, error) {
	query := `
		SELECT id, name, ip, port, username, description, status, ssh_key_path, auth_type, created_at, updated_at
		FROM servers
		ORDER BY id
	`

	rows, err := GetDB().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servers := make([]*Server, 0)
	for rows.Next() {
		server := &Server{}
		var sshKeyPath, authType sql.NullString

		err := rows.Scan(
			&server.ID, &server.Name, &server.IP, &server.Port, &server.Username,
			&server.Description, &server.Status, &sshKeyPath, &authType, &server.CreatedAt, &server.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		if sshKeyPath.Valid {
			server.SSHKeyPath = sshKeyPath.String
		}
		if authType.Valid {
			server.AuthType = authType.String
		} else {
			server.AuthType = "password" // Default to password auth if NULL
		}

		servers = append(servers, server)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return servers, nil
}

// UpdateServer 更新服务器
func UpdateServer(server *Server) error {
	server.UpdatedAt = time.Now()

	query := `
		UPDATE servers
		SET name = $1, ip = $2, port = $3, username = $4, description = $5, status = $6, ssh_key_path = $7, auth_type = $8, updated_at = $9
		WHERE id = $10
	`

	_, err := GetDB().Exec(context.Background(), query,
		server.Name, server.IP, server.Port, server.Username,
		server.Description, server.Status, server.SSHKeyPath, server.AuthType, server.UpdatedAt, server.ID)
	return err
}

// UpdateServerPassword 更新服务器密码
func UpdateServerPassword(serverID int, newPassword string) error {
	query := `
		UPDATE servers
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := GetDB().Exec(context.Background(), query, newPassword, time.Now(), serverID)
	return err
}

// DeleteServer 删除服务器
func DeleteServer(id int) error {
	query := `DELETE FROM servers WHERE id = $1`
	_, err := GetDB().Exec(context.Background(), query, id)
	return err
}
