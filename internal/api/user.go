package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin197011/kk-ops/internal/model"
)

// UserRequest 用户请求结构
type UserRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Role     string `json:"role" form:"role" binding:"required"`
}

// ListUsersHandler 列出所有用户
func ListUsersHandler(c *gin.Context) {
	users, err := model.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		users = []*model.User{} // 使用空数组避免模板错误
	}

	// 获取用户信息用于导航栏显示
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":      "用户管理 - KK运维平台",
		"users":      users,
		"activeMenu": "users",
		"userID":     userID,
		"username":   username,
		"role":       role,
	})
}

// GetUserHandler 获取用户详情
func GetUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("GetUserHandler id: %d", id)
	if err != nil {
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"title":      "用户详情 - KK运维平台",
			"error":      "用户ID无效",
			"activeMenu": "users",
		})
		return
	}

	user, err := model.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusOK, "user_detail.html", gin.H{
			"title":      "用户详情 - KK运维平台",
			"error":      "获取用户详情失败: " + err.Error(),
			"activeMenu": "users",
		})
		return
	}

	// 获取用户信息用于导航栏显示
	currentUserID, _ := c.Get("userID")
	currentUsername, _ := c.Get("username")
	currentRole, _ := c.Get("role")

	c.HTML(http.StatusOK, "user_detail.html", gin.H{
		"title":      "用户详情 - KK运维平台",
		"user":       user,
		"activeMenu": "users",
		"userID":     currentUserID,
		"username":   currentUsername,
		"role":       currentRole,
	})
}

// CreateUserHandler 创建用户
func CreateUserHandler(c *gin.Context) {
	// 显示添加用户表单页面
	if c.Request.Method == "GET" {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "user_new.html", gin.H{
			"title":      "添加用户 - KK运维平台",
			"activeMenu": "users",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	var req UserRequest
	if err := c.ShouldBind(&req); err != nil {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "user_new.html", gin.H{
			"title":      "添加用户 - KK运维平台",
			"error":      "请填写所有必填字段",
			"activeMenu": "users",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	// 检查用户名是否已存在
	_, err := model.GetUserByUsername(req.Username)
	if err == nil {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "user_new.html", gin.H{
			"title":      "添加用户 - KK运维平台",
			"error":      "用户名已存在",
			"activeMenu": "users",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

	if err := model.CreateUser(user); err != nil {
		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "user_new.html", gin.H{
			"title":      "添加用户 - KK运维平台",
			"error":      "创建用户失败: " + err.Error(),
			"user":       user,
			"activeMenu": "users",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

// UpdateUserHandler 更新用户
func UpdateUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("UpdateUserHandler id: %d", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	// 先获取现有用户
	user, err := model.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var req UserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写所有必填字段"})
		return
	}

	// 检查是否修改了用户名并确保新用户名不与其他用户冲突
	if user.Username != req.Username {
		existingUser, err := model.GetUserByUsername(req.Username)
		if err == nil && existingUser.ID != id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已被占用"})
			return
		}
	}

	// 更新用户信息
	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role

	if err := model.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败: " + err.Error()})
		return
	}

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		if err := model.UpdatePassword(id, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功", "user": user})
}

// DeleteUserHandler 删除用户
func DeleteUserHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	log.Printf("DeleteUserHandler id: %d", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	// 检查是否是删除自己
	currentUserID, exists := c.Get("userID")
	if exists && currentUserID.(int) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录的用户"})
		return
	}

	if err := model.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
