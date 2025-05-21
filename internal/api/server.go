package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin197011/kk-ops/internal/model"
)

// ServerRequest 服务器请求结构
type ServerRequest struct {
	Name        string `json:"name" form:"name" binding:"required"`
	IP          string `json:"ip" form:"ip" binding:"required,ip"`
	Port        int    `json:"port" form:"port" binding:"required,numeric"`
	Username    string `json:"username" form:"username" binding:"required"`
	Password    string `json:"password" form:"password"`
	Description string `json:"description" form:"description"`
	Status      string `json:"status" form:"status" binding:"required"`
	SSHKeyPath  string `json:"ssh_key_path" form:"ssh_key_path"`
	AuthType    string `json:"auth_type" form:"auth_type" binding:"required"`
}

// ListServersHandler 列出所有服务器
func ListServersHandler(c *gin.Context) {
	servers, err := model.GetAllServers()
	if err != nil {
		log.Printf("Error getting servers: %v", err)
		servers = []*model.Server{} // 使用空数组避免模板错误
	}

	// 获取用户信息用于导航栏显示
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.HTML(http.StatusOK, "servers.html", gin.H{
		"title":      "服务器管理 - KK-OPS",
		"servers":    servers,
		"activeMenu": "servers",
		"userID":     userID,
		"username":   username,
		"role":       role,
	})
}

// GetServerHandler 获取服务器详情
func GetServerHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("GetServerHandler id: %d", id)
	if err != nil {
		c.HTML(http.StatusOK, "server_detail.html", gin.H{
			"title":      "服务器详情 - KK-OPS",
			"error":      "服务器ID无效",
			"activeMenu": "servers",
		})
		return
	}

	server, err := model.GetServerByID(id)
	if err != nil {
		c.HTML(http.StatusOK, "server_detail.html", gin.H{
			"title":      "服务器详情 - KK-OPS",
			"error":      "获取服务器详情失败: " + err.Error(),
			"activeMenu": "servers",
		})
		return
	}

	// 获取用户信息用于导航栏显示
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.HTML(http.StatusOK, "server_detail.html", gin.H{
		"title":      "服务器详情 - KK-OPS",
		"server":     server,
		"activeMenu": "servers",
		"userID":     userID,
		"username":   username,
		"role":       role,
	})
}

// CreateServerHandler 创建服务器
func CreateServerHandler(c *gin.Context) {
	// 显示添加服务器表单页面
	if c.Request.Method == "GET" {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "server_new.html", gin.H{
			"title":      "添加服务器 - KK-OPS",
			"activeMenu": "servers",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	var req ServerRequest
	if err := c.ShouldBind(&req); err != nil {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "server_new.html", gin.H{
			"title":      "添加服务器 - KK-OPS",
			"error":      "请填写所有必填字段",
			"activeMenu": "servers",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	server := &model.Server{
		Name:        req.Name,
		IP:          req.IP,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Description: req.Description,
		Status:      req.Status,
		SSHKeyPath:  req.SSHKeyPath,
		AuthType:    req.AuthType,
	}

	if err := model.CreateServer(server); err != nil {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "server_new.html", gin.H{
			"title":      "添加服务器 - KK-OPS",
			"error":      "创建服务器失败: " + err.Error(),
			"server":     server,
			"activeMenu": "servers",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	c.Redirect(http.StatusFound, "/servers")
}

// UpdateServerHandler 更新服务器
func UpdateServerHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("UpdateServerHandler id: %d", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器ID无效"})
		return
	}

	// 先获取现有服务器
	server, err := model.GetServerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	var req ServerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写所有必填字段"})
		return
	}

	// 更新服务器信息
	server.Name = req.Name
	server.IP = req.IP
	server.Port = req.Port
	server.Username = req.Username
	server.Description = req.Description
	server.Status = req.Status
	server.SSHKeyPath = req.SSHKeyPath
	server.AuthType = req.AuthType

	if err := model.UpdateServer(server); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器失败: " + err.Error()})
		return
	}

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		if err := model.UpdateServerPassword(id, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务器密码失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务器更新成功", "server": server})
}

// DeleteServerHandler 删除服务器
func DeleteServerHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器ID无效"})
		return
	}

	if err := model.DeleteServer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除服务器失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务器删除成功"})
}

// SSHConnectHandler 处理SSH连接
func SSHConnectHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "服务器ID无效"})
		return
	}

	// 获取服务器详情
	server, err := model.GetServerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务器不存在"})
		return
	}

	// 构建SSH连接信息
	sshInfo := gin.H{
		"server_id":    server.ID,
		"host":         server.IP,
		"port":         server.Port,
		"username":     server.Username,
		"auth_type":    server.AuthType,
		"ssh_key_path": server.SSHKeyPath,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "SSH连接信息准备就绪",
		"ssh_info": sshInfo,
	})
}
