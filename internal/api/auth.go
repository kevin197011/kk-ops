package api

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevin197011/kk-ops/internal/middleware"
	"github.com/kevin197011/kk-ops/internal/model"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
}

// IndexHandler 首页处理器
func IndexHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/dashboard")
}

// ShowLoginPage 显示登录页面
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录 - KK-OPS运维管理平台",
	})
}

// LoginHandler 登录处理器
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "登录 - KK-OPS运维管理平台",
			"error": "用户名和密码不能为空",
		})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil || !user.CheckPassword(req.Password) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "登录 - KK-OPS运维管理平台",
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "登录 - KK-OPS运维管理平台",
			"error": "系统错误，请稍后再试",
		})
		return
	}

	// 更新最后登录时间
	if err := model.UpdateLastLogin(user.ID); err != nil {
		log.Printf("Failed to update last login time for user %d: %v", user.ID, err)
		// 不阻止登录过程继续
	}

	// 设置Cookie
	c.SetCookie("token", token, 86400, "/", "", false, true)

	c.Redirect(http.StatusFound, "/dashboard")
}

// ShowRegisterPage 显示注册页面
func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "注册 - KK-OPS运维管理平台",
	})
}

// RegisterHandler 注册处理器
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "注册 - KK-OPS运维管理平台",
			"error": "请填写所有必填字段",
		})
		return
	}

	// 检查用户名是否已存在
	_, err := model.GetUserByUsername(req.Username)
	if err == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "注册 - KK-OPS运维管理平台",
			"error": "用户名已存在",
		})
		return
	}

	// 创建新用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "user", // 默认普通用户角色
	}

	if err := model.CreateUser(user); err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "注册 - KK-OPS运维管理平台",
			"error": "注册失败，请稍后再试",
		})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "注册 - KK-OPS运维管理平台",
			"error": "系统错误，请稍后再试",
		})
		return
	}

	// 设置Cookie
	c.SetCookie("token", token, 86400, "/", "", false, true)

	c.Redirect(http.StatusFound, "/dashboard")
}

// LogoutHandler 登出处理器
func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
