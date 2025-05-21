package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kevin197011/kk-ops/internal/api"
	"github.com/kevin197011/kk-ops/internal/config"
	"github.com/kevin197011/kk-ops/internal/middleware"
	"github.com/kevin197011/kk-ops/internal/model"
)

func main() {
	// 加载环境配置
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	if err := model.InitDB(); err != nil {
		log.Printf("Database initialization failed: %v", err)
		log.Println("Continuing with limited functionality - database features will not work")
	} else {
		// 只有在成功初始化数据库的情况下才注册关闭函数
		defer model.CloseDB()
	}

	// 创建Gin实例
	r := gin.Default()

	// 设置模板函数
	r.SetFuncMap(template.FuncMap{
		"subtract": func(a, b int) int { return a - b },
	})

	// 加载静态文件
	r.Static("/static/css", "./ui/static/css")
	r.Static("/static/js", "./ui/static/js")
	r.LoadHTMLGlob("ui/templates/*")

	// 注册中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// 不需要认证的路由
	r.GET("/", api.IndexHandler)
	r.GET("/login", api.ShowLoginPage)
	r.POST("/login", api.LoginHandler)
	r.GET("/register", api.ShowRegisterPage)
	r.POST("/register", api.RegisterHandler)
	r.GET("/logout", api.LogoutHandler)

	// 需要认证的路由
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 仪表盘
		authorized.GET("/dashboard", api.DashboardHandler)

		// 服务器管理
		servers := authorized.Group("/servers")
		{
			servers.GET("", api.ListServersHandler)
			servers.GET("/new", api.CreateServerHandler)
			servers.GET("/:id", api.GetServerHandler)
			servers.POST("", api.CreateServerHandler)
			servers.PUT("/:id", api.UpdateServerHandler)
			servers.DELETE("/:id", api.DeleteServerHandler)
			servers.GET("/:id/ssh", api.SSHConnectHandler)
		}

		// CICD管理
		cicd := authorized.Group("/cicd")
		{
			cicd.GET("/pipelines", api.ListPipelinesHandler)
			cicd.GET("/pipelines/new", api.CreatePipelineHandler) // Must come before the :id route
			cicd.GET("/pipelines/:id", api.GetPipelineHandler)
			cicd.POST("/pipelines", api.CreatePipelineHandler)
			cicd.PUT("/pipelines/:id", api.UpdatePipelineHandler)
			cicd.DELETE("/pipelines/:id", api.DeletePipelineHandler)
			cicd.POST("/pipelines/:id/run", api.RunPipelineHandler)
		}

		// 用户管理
		users := authorized.Group("/users")
		{
			users.GET("", api.ListUsersHandler)
			users.GET("/new", api.CreateUserHandler) // Must come before the :id route
			users.GET("/:id", api.GetUserHandler)
			users.POST("", api.CreateUserHandler)
			users.PUT("/:id", api.UpdateUserHandler)
			users.DELETE("/:id", api.DeleteUserHandler)
		}
	}

	// 优雅停机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		model.CloseDB()
		os.Exit(0)
	}()

	// 启动服务器
	port := config.GetConfig().ServerPort
	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
