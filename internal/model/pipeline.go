package model

import (
	"context"
	"time"
)

// Pipeline CICD流水线模型
type Pipeline struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	GitRepo      string     `json:"git_repo"`
	Branch       string     `json:"branch"`
	BuildScript  string     `json:"build_script"`
	DeployScript string     `json:"deploy_script"`
	ServerID     int        `json:"server_id"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastRunAt    *time.Time `json:"last_run_at"`
}

// PipelineRun 流水线运行记录
type PipelineRun struct {
	ID          int       `json:"id"`
	PipelineID  int       `json:"pipeline_id"`
	Status      string    `json:"status"` // pending, running, success, failed
	Logs        string    `json:"logs"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TriggerUser int       `json:"trigger_user"`
}

// CreatePipeline 创建流水线
func CreatePipeline(pipeline *Pipeline) error {
	pipeline.CreatedAt = time.Now()
	pipeline.UpdatedAt = time.Now()
	pipeline.Status = "inactive"

	query := `
		INSERT INTO pipelines (name, description, git_repo, branch, build_script, deploy_script, server_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	return GetDB().QueryRow(context.Background(), query,
		pipeline.Name, pipeline.Description, pipeline.GitRepo, pipeline.Branch,
		pipeline.BuildScript, pipeline.DeployScript, pipeline.ServerID,
		pipeline.Status, pipeline.CreatedAt, pipeline.UpdatedAt).Scan(&pipeline.ID)
}

// GetPipelineByID 根据ID获取流水线
func GetPipelineByID(id int) (*Pipeline, error) {
	query := `
		SELECT id, name, description, git_repo, branch, build_script, deploy_script, server_id, status, created_at, updated_at, last_run_at
		FROM pipelines
		WHERE id = $1
	`

	pipeline := &Pipeline{}
	err := GetDB().QueryRow(context.Background(), query, id).Scan(
		&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.GitRepo, &pipeline.Branch,
		&pipeline.BuildScript, &pipeline.DeployScript, &pipeline.ServerID, &pipeline.Status,
		&pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.LastRunAt)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

// GetAllPipelines 获取所有流水线
func GetAllPipelines() ([]*Pipeline, error) {
	query := `
		SELECT id, name, description, git_repo, branch, server_id, status, created_at, updated_at, last_run_at
		FROM pipelines
		ORDER BY id
	`

	rows, err := GetDB().Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pipelines := make([]*Pipeline, 0)
	for rows.Next() {
		pipeline := &Pipeline{}
		err := rows.Scan(
			&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.GitRepo, &pipeline.Branch,
			&pipeline.ServerID, &pipeline.Status, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.LastRunAt)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, pipeline)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pipelines, nil
}

// UpdatePipeline 更新流水线
func UpdatePipeline(pipeline *Pipeline) error {
	pipeline.UpdatedAt = time.Now()

	query := `
		UPDATE pipelines
		SET name = $1, description = $2, git_repo = $3, branch = $4, build_script = $5, deploy_script = $6, server_id = $7, status = $8, updated_at = $9
		WHERE id = $10
	`

	_, err := GetDB().Exec(context.Background(), query,
		pipeline.Name, pipeline.Description, pipeline.GitRepo, pipeline.Branch,
		pipeline.BuildScript, pipeline.DeployScript, pipeline.ServerID,
		pipeline.Status, pipeline.UpdatedAt, pipeline.ID)
	return err
}

// DeletePipeline 删除流水线
func DeletePipeline(id int) error {
	query := `DELETE FROM pipelines WHERE id = $1`
	_, err := GetDB().Exec(context.Background(), query, id)
	return err
}

// CreatePipelineRun 创建流水线运行记录
func CreatePipelineRun(run *PipelineRun) error {
	run.StartTime = time.Now()
	run.Status = "pending"

	query := `
		INSERT INTO pipeline_runs (pipeline_id, status, logs, start_time, trigger_user)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := GetDB().QueryRow(context.Background(), query,
		run.PipelineID, run.Status, run.Logs, run.StartTime, run.TriggerUser).Scan(&run.ID)
	if err != nil {
		return err
	}

	// 更新流水线最后运行时间
	updateQuery := `
		UPDATE pipelines
		SET last_run_at = $1, status = 'running'
		WHERE id = $2
	`
	_, err = GetDB().Exec(context.Background(), updateQuery, run.StartTime, run.PipelineID)
	return err
}

// UpdatePipelineRunStatus 更新流水线运行状态
func UpdatePipelineRunStatus(runID int, status string, logs string) error {
	query := `
		UPDATE pipeline_runs
		SET status = $1, logs = $2
	`

	args := []interface{}{status, logs}

	if status == "success" || status == "failed" {
		query += `, end_time = $3 WHERE id = $4`
		args = append(args, time.Now(), runID)
	} else {
		query += ` WHERE id = $3`
		args = append(args, runID)
	}

	_, err := GetDB().Exec(context.Background(), query, args...)
	return err
}

// GetPipelineRuns 获取流水线运行记录
func GetPipelineRuns(pipelineID int, limit int) ([]*PipelineRun, error) {
	query := `
		SELECT id, pipeline_id, status, logs, start_time, end_time, trigger_user
		FROM pipeline_runs
		WHERE pipeline_id = $1
		ORDER BY start_time DESC
		LIMIT $2
	`

	rows, err := GetDB().Query(context.Background(), query, pipelineID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runs := make([]*PipelineRun, 0)
	for rows.Next() {
		run := &PipelineRun{}
		err := rows.Scan(&run.ID, &run.PipelineID, &run.Status, &run.Logs, &run.StartTime, &run.EndTime, &run.TriggerUser)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}
