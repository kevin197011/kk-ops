package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevin197011/kk-ops/internal/model"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	ServerCount          int       `json:"server_count"`
	OnlineServerCount    int       `json:"online_server_count"`
	PipelineCount        int       `json:"pipeline_count"`
	PipelineSuccessCount int       `json:"pipeline_success_count"`
	PipelineFailedCount  int       `json:"pipeline_failed_count"`
	PipelinePendingCount int       `json:"pipeline_pending_count"`
	UserCount            int       `json:"user_count"`
	LastUpdated          time.Time `json:"last_updated"`
}

// PipelineRunInfo 流水线运行信息
type PipelineRunInfo struct {
	RunID        int       `json:"run_id"`
	PipelineID   int       `json:"pipeline_id"`
	PipelineName string    `json:"pipeline_name"`
	Status       string    `json:"status"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	UserName     string    `json:"user_name"`
}

// PipelineRun 表示一次流水线运行记录
type PipelineRun struct {
	PipelineID   int       `json:"pipeline_id"`
	PipelineName string    `json:"pipeline_name"`
	Status       string    `json:"status"`
	UserID       int       `json:"user_id"`
	UserName     string    `json:"user_name"`
	StartTime    time.Time `json:"start_time"`
}

// DashboardHandler 显示仪表盘页面
func DashboardHandler(c *gin.Context) {
	// 获取当前登录用户信息
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	log.Printf("Dashboard - User ID: %v, Username: %v, Role: %v", userID, username, role)

	// 处理可能的数据库连接问题
	var stats *DashboardStats
	var err error

	// 尝试获取统计数据，但处理可能的错误
	stats, err = getStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		// 如果无法获取真实数据，使用空数据
		stats = &DashboardStats{
			ServerCount:          0,
			OnlineServerCount:    0,
			PipelineCount:        0,
			PipelineSuccessCount: 0,
			PipelineFailedCount:  0,
			PipelinePendingCount: 0,
			UserCount:            0,
			LastUpdated:          time.Now(),
		}
	}

	// 获取最近的流水线运行，但处理可能的错误
	recentRuns, err := getRecentPipelineRuns(5)
	if err != nil {
		log.Printf("Error getting recent runs: %v", err)
		recentRuns = []PipelineRunInfo{} // 使用空数组
	}

	data := gin.H{
		"title":      "仪表盘 - KK运维平台",
		"userID":     userID,
		"username":   username,
		"role":       role,
		"stats":      stats,
		"recentRuns": recentRuns,
		"activeMenu": "dashboard", // 添加当前活动菜单
	}

	log.Printf("Dashboard data: %+v", data)

	c.HTML(http.StatusOK, "dashboard.html", data)
}

// getStats 获取统计数据
func getStats() (*DashboardStats, error) {
	// 使用同一个数据库连接
	db := model.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	ctx := context.Background()

	stats := &DashboardStats{
		ServerCount:          0,
		OnlineServerCount:    0,
		PipelineCount:        0,
		PipelineSuccessCount: 0,
		PipelineFailedCount:  0,
		PipelinePendingCount: 0,
		UserCount:            0,
		LastUpdated:          time.Now(),
	}

	// 获取服务器总数
	serverQuery := `SELECT COUNT(*) FROM servers`
	err := db.QueryRow(ctx, serverQuery).Scan(&stats.ServerCount)
	if err != nil {
		log.Printf("Error counting servers: %v", err)
		// 继续执行，不要完全中断
	}

	// 获取在线服务器数量
	onlineQuery := `SELECT COUNT(*) FROM servers WHERE status = 'online'`
	err = db.QueryRow(ctx, onlineQuery).Scan(&stats.OnlineServerCount)
	if err != nil {
		log.Printf("Error counting online servers: %v", err)
	}

	// 获取流水线总数
	pipelineQuery := `SELECT COUNT(*) FROM pipelines`
	err = db.QueryRow(ctx, pipelineQuery).Scan(&stats.PipelineCount)
	if err != nil {
		log.Printf("Error counting pipelines: %v", err)
	}

	// 获取成功的流水线数量
	successQuery := `SELECT COUNT(*) FROM pipelines WHERE status = 'success'`
	err = db.QueryRow(ctx, successQuery).Scan(&stats.PipelineSuccessCount)
	if err != nil {
		log.Printf("Error counting success pipelines: %v", err)
	}

	// 获取失败的流水线数量
	failedQuery := `SELECT COUNT(*) FROM pipelines WHERE status = 'failed'`
	err = db.QueryRow(ctx, failedQuery).Scan(&stats.PipelineFailedCount)
	if err != nil {
		log.Printf("Error counting failed pipelines: %v", err)
	}

	// 获取待处理/正在执行的流水线数量
	pendingQuery := `SELECT COUNT(*) FROM pipelines WHERE status IN ('pending', 'running')`
	err = db.QueryRow(ctx, pendingQuery).Scan(&stats.PipelinePendingCount)
	if err != nil {
		log.Printf("Error counting pending pipelines: %v", err)
	}

	// 获取用户总数
	userQuery := `SELECT COUNT(*) FROM users`
	err = db.QueryRow(ctx, userQuery).Scan(&stats.UserCount)
	if err != nil {
		log.Printf("Error counting users: %v", err)
	}

	return stats, nil
}

// getRecentPipelineRuns 获取最近的流水线运行记录
func getRecentPipelineRuns(limit int) ([]PipelineRunInfo, error) {
	db := model.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	ctx := context.Background()

	// 处理pipeline_runs表可能不存在的情况
	// 首先检查表是否存在
	checkTableQuery := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'pipeline_runs'
		)
	`
	var tableExists bool
	err := db.QueryRow(ctx, checkTableQuery).Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking if pipeline_runs table exists: %v", err)
		return []PipelineRunInfo{}, nil
	}

	if !tableExists {
		log.Println("pipeline_runs table does not exist yet")
		return []PipelineRunInfo{}, nil
	}

	query := `
		SELECT
			pr.id, pr.pipeline_id, p.name, pr.status, pr.start_time, pr.end_time, u.username
		FROM
			pipeline_runs pr
		JOIN
			pipelines p ON pr.pipeline_id = p.id
		JOIN
			users u ON pr.trigger_user = u.id
		ORDER BY
			pr.start_time DESC
		LIMIT $1
	`

	rows, err := db.Query(ctx, query, limit)
	if err != nil {
		log.Printf("Error querying pipeline runs: %v", err)
		return []PipelineRunInfo{}, nil
	}
	defer rows.Close()

	runs := make([]PipelineRunInfo, 0)
	for rows.Next() {
		var run PipelineRunInfo
		var endTime interface{} // 使用interface{}来处理可能为NULL的end_time

		if err := rows.Scan(&run.RunID, &run.PipelineID, &run.PipelineName, &run.Status, &run.StartTime, &endTime, &run.UserName); err != nil {
			log.Printf("Error scanning pipeline run row: %v", err)
			continue
		}

		// 处理可能为NULL的end_time
		if t, ok := endTime.(time.Time); ok {
			run.EndTime = t
		}

		runs = append(runs, run)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating pipeline run rows: %v", err)
	}

	return runs, nil
}
