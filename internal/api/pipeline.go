package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevin197011/kk-ops/internal/model"
)

// PipelineRequest 流水线请求结构
type PipelineRequest struct {
	Name         string `json:"name" form:"name" binding:"required"`
	Description  string `json:"description" form:"description"`
	GitRepo      string `json:"git_repo" form:"git_repo" binding:"required"`
	Branch       string `json:"branch" form:"branch" binding:"required"`
	BuildScript  string `json:"build_script" form:"build_script" binding:"required"`
	DeployScript string `json:"deploy_script" form:"deploy_script" binding:"required"`
	ServerID     int    `json:"server_id" form:"server_id" binding:"required,numeric"`
	Status       string `json:"status" form:"status"`
}

// ListPipelinesHandler 列出所有流水线
func ListPipelinesHandler(c *gin.Context) {
	pipelines, err := model.GetAllPipelines()
	if err != nil {
		log.Printf("Error getting pipelines: %v", err)
		pipelines = []*model.Pipeline{} // 使用空数组避免模板错误
	}

	// 获取用户信息用于导航栏显示
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.HTML(http.StatusOK, "pipelines.html", gin.H{
		"title":      "CICD管理 - KK-OPS",
		"pipelines":  pipelines,
		"activeMenu": "cicd",
		"userID":     userID,
		"username":   username,
		"role":       role,
	})
}

// GetPipelineHandler 获取流水线详情
func GetPipelineHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting idStr to int: %v", err)
		c.HTML(http.StatusOK, "pipeline_detail.html", gin.H{
			"title":      "流水线详情 - KK-OPS",
			"error":      "流水线ID无效",
			"activeMenu": "cicd",
		})
		return
	}

	pipeline, err := model.GetPipelineByID(id)
	if err != nil {
		log.Printf("Error getting pipeline by id: %v", err)
		c.HTML(http.StatusOK, "pipeline_detail.html", gin.H{
			"title":      "流水线详情 - KK-OPS",
			"error":      "获取流水线详情失败: " + err.Error(),
			"activeMenu": "cicd",
		})
		return
	}

	// 获取该流水线的最近5次运行记录
	runs, err := model.GetPipelineRuns(id, 5)
	if err != nil {
		log.Printf("Error getting pipeline runs: %v", err)
		runs = []*model.PipelineRun{} // 使用空数组避免模板错误
	}

	// 获取服务器详情
	server, err := model.GetServerByID(pipeline.ServerID)
	if err != nil {
		log.Printf("Error getting server: %v", err)
	}

	// 获取用户信息用于导航栏显示
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.HTML(http.StatusOK, "pipeline_detail.html", gin.H{
		"title":      "流水线详情 - KK-OPS",
		"pipeline":   pipeline,
		"runs":       runs,
		"server":     server,
		"activeMenu": "cicd",
		"userID":     userID,
		"username":   username,
		"role":       role,
	})
}

// CreatePipelineHandler 创建流水线
func CreatePipelineHandler(c *gin.Context) {
	// 显示添加流水线表单页面
	if c.Request.Method == "GET" {
		// 获取所有服务器列表，用于下拉选择
		servers, err := model.GetAllServers()
		if err != nil {
			log.Printf("Error getting servers: %v", err)
			servers = []*model.Server{} // 使用空数组避免模板错误
		}

		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "pipeline_new.html", gin.H{
			"title":      "添加流水线 - KK-OPS",
			"servers":    servers,
			"activeMenu": "cicd",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	var req PipelineRequest
	if err := c.ShouldBind(&req); err != nil {
		// 获取所有服务器列表，用于下拉选择
		servers, err := model.GetAllServers()
		if err != nil {
			log.Printf("Error getting servers: %v", err)
			servers = []*model.Server{} // 使用空数组避免模板错误
		}

		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "pipeline_new.html", gin.H{
			"title":      "添加流水线 - KK-OPS",
			"error":      "请填写所有必填字段",
			"servers":    servers,
			"activeMenu": "cicd",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	pipeline := &model.Pipeline{
		Name:         req.Name,
		Description:  req.Description,
		GitRepo:      req.GitRepo,
		Branch:       req.Branch,
		BuildScript:  req.BuildScript,
		DeployScript: req.DeployScript,
		ServerID:     req.ServerID,
		Status:       "inactive", // 默认为未激活状态
	}

	if err := model.CreatePipeline(pipeline); err != nil {
		// 获取所有服务器列表，用于下拉选择
		servers, err := model.GetAllServers()
		if err != nil {
			log.Printf("Error getting servers: %v", err)
			servers = []*model.Server{} // 使用空数组避免模板错误
		}

		// 获取用户信息用于导航栏显示
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.HTML(http.StatusOK, "pipeline_new.html", gin.H{
			"title":      "添加流水线 - KK-OPS",
			"error":      "创建流水线失败: " + err.Error(),
			"pipeline":   pipeline,
			"servers":    servers,
			"activeMenu": "cicd",
			"userID":     userID,
			"username":   username,
			"role":       role,
		})
		return
	}

	c.Redirect(http.StatusFound, "/cicd/pipelines")
}

// UpdatePipelineHandler 更新流水线
func UpdatePipelineHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "流水线ID无效"})
		return
	}

	// 先获取现有流水线
	pipeline, err := model.GetPipelineByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "流水线不存在"})
		return
	}

	var req PipelineRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写所有必填字段"})
		return
	}

	// 更新流水线信息
	pipeline.Name = req.Name
	pipeline.Description = req.Description
	pipeline.GitRepo = req.GitRepo
	pipeline.Branch = req.Branch
	pipeline.BuildScript = req.BuildScript
	pipeline.DeployScript = req.DeployScript
	pipeline.ServerID = req.ServerID

	if req.Status != "" {
		pipeline.Status = req.Status
	}

	if err := model.UpdatePipeline(pipeline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新流水线失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "流水线更新成功", "pipeline": pipeline})
}

// DeletePipelineHandler 删除流水线
func DeletePipelineHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "流水线ID无效"})
		return
	}

	if err := model.DeletePipeline(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除流水线失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "流水线删除成功"})
}

// RunPipelineHandler 运行流水线
func RunPipelineHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "流水线ID无效"})
		return
	}

	// 获取流水线详情
	_, err = model.GetPipelineByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "流水线不存在"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的操作"})
		return
	}

	// 创建流水线运行记录
	pipelineRun := &model.PipelineRun{
		PipelineID:  id,
		Status:      "pending",
		Logs:        "流水线启动中...\n",
		TriggerUser: userID.(int),
	}

	if err := model.CreatePipelineRun(pipelineRun); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动流水线失败: " + err.Error()})
		return
	}

	// TODO: 在后台协程中实际执行流水线任务
	// 这里可以启动一个goroutine来执行实际的构建和部署任务

	c.JSON(http.StatusOK, gin.H{
		"message": "流水线启动成功",
		"run_id":  pipelineRun.ID,
	})
}
