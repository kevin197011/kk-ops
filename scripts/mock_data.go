package scripts

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kevin197011/kk-ops/internal/model"
)

// MockData 生成测试数据
func MockData() {
	log.Println("开始生成测试数据...")

	// 初始化数据库连接
	if err := model.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer model.CloseDB()

	db := model.GetDB()
	if db == nil {
		log.Fatalf("无法获取数据库连接")
	}

	// 创建用户数据
	createMockUsers()

	// 创建服务器数据
	createMockServers()

	// 创建流水线数据
	createMockPipelines()

	log.Println("测试数据生成完成!")
}

// createMockUsers 创建测试用户
func createMockUsers() {
	log.Println("创建测试用户...")

	// 检查是否已有管理员用户
	adminExists, _ := model.GetUserByUsername("admin")
	if adminExists == nil {
		// 创建管理员用户
		admin := &model.User{
			Username:  "admin",
			Password:  "admin123", // 会在CreateUser方法中被哈希
			Email:     "admin@example.com",
			Role:      "admin",
			CreatedAt: time.Now(),
		}

		if err := model.CreateUser(admin); err != nil {
			log.Printf("创建管理员用户失败: %v", err)
		} else {
			log.Println("创建管理员用户成功")
		}
	} else {
		log.Println("管理员用户已存在，跳过创建")
	}

	// 创建普通用户
	usernames := []string{"user1", "user2", "devops1", "tester1"}
	for _, username := range usernames {
		userExists, _ := model.GetUserByUsername(username)
		if userExists == nil {
			user := &model.User{
				Username:  username,
				Password:  fmt.Sprintf("%s123", username), // 会在CreateUser方法中被哈希
				Email:     fmt.Sprintf("%s@example.com", username),
				Role:      "user",
				CreatedAt: time.Now(),
			}

			if err := model.CreateUser(user); err != nil {
				log.Printf("创建用户 %s 失败: %v", username, err)
			} else {
				log.Printf("创建用户 %s 成功", username)
			}
		} else {
			log.Printf("用户 %s 已存在，跳过创建", username)
		}
	}
}

// createMockServers 创建测试服务器
func createMockServers() {
	log.Println("创建测试服务器...")

	serverTypes := []string{"web", "app", "db", "cache", "api"}
	envs := []string{"dev", "test", "stage", "prod"}

	for i := 1; i <= 10; i++ {
		// 生成服务器名称 sv-{环境}-{类型}{序号}
		serverType := serverTypes[i%len(serverTypes)]
		env := envs[i%len(envs)]
		serverName := fmt.Sprintf("sv-%s-%s%d", env, serverType, i)

		// 检查是否有同名服务器
		var count int
		err := model.GetDB().QueryRow(context.Background(),
			"SELECT COUNT(*) FROM servers WHERE name = $1", serverName).Scan(&count)

		if err != nil || count > 0 {
			log.Printf("服务器 %s 已存在或查询失败，跳过创建", serverName)
			continue
		}

		// 创建服务器
		server := &model.Server{
			Name:        serverName,
			IP:          fmt.Sprintf("192.168.1.%d", 100+i),
			Port:        22,
			Username:    "root",
			Password:    "password123",
			Description: fmt.Sprintf("%s环境%s服务器", env, serverType),
			Status:      "running",
			AuthType:    "password",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := model.CreateServer(server); err != nil {
			log.Printf("创建服务器 %s 失败: %v", serverName, err)
		} else {
			log.Printf("创建服务器 %s 成功", serverName)
		}
	}
}

// createMockPipelines 创建测试流水线
func createMockPipelines() {
	log.Println("创建测试流水线...")

	pipelineNames := []string{
		"前端构建部署",
		"后端服务部署",
		"数据库脚本执行",
		"自动化测试",
		"安全扫描",
	}

	for i, name := range pipelineNames {
		// 查询数据库中是否已存在该流水线
		ctx := context.Background()
		var count int
		err := model.GetDB().QueryRow(ctx,
			"SELECT COUNT(*) FROM pipelines WHERE name = $1", name).Scan(&count)
		if err != nil {
			log.Printf("查询流水线 %s 失败: %v", name, err)
			continue
		}

		if count == 0 {
			// 创建流水线
			pipeline := &model.Pipeline{
				Name:         name,
				Description:  fmt.Sprintf("%s流水线描述", name),
				GitRepo:      fmt.Sprintf("https://github.com/example/%s.git", name),
				Branch:       "main",
				BuildScript:  fmt.Sprintf("echo '构建%s'\nnpm install\nnpm run build", name),
				DeployScript: fmt.Sprintf("echo '部署%s'\nscp -r dist/ server:/opt/app/", name),
				ServerID:     i + 1, // 假设服务器ID从1开始
				Status:       "ready",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := model.CreatePipeline(pipeline); err != nil {
				log.Printf("创建流水线 %s 失败: %v", name, err)
			} else {
				log.Printf("创建流水线 %s 成功", name)

				// 为前两个流水线添加一些运行记录
				if i < 2 {
					for j := 0; j < 3; j++ {
						runTime := time.Now().Add(time.Duration(-j*24) * time.Hour)
						status := "success"
						if j == 1 {
							status = "failed"
						}

						_, err := model.GetDB().Exec(ctx,
							"INSERT INTO pipeline_runs (pipeline_id, status, logs, start_time, end_time, trigger_user) VALUES ($1, $2, $3, $4, $5, $6)",
							pipeline.ID, status, fmt.Sprintf("流水线运行输出示例-%d", j), runTime, runTime.Add(5*time.Minute), 1)

						if err != nil {
							log.Printf("为流水线 %s 创建运行记录失败: %v", name, err)
						} else {
							log.Printf("为流水线 %s 创建运行记录成功", name)
						}
					}
				}
			}
		} else {
			log.Printf("流水线 %s 已存在，跳过创建", name)
		}
	}
}
