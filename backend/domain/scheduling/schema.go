package scheduling

import "time"

type Task struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	TaskType       string     `json:"task_type"`
	SchedulerID    int64      `json:"scheduler_id"`
	CronExpression string     `json:"cron_expression"`
	Enabled        bool       `json:"enabled"`
	NextRunAt      *time.Time `json:"next_run_at"`
	LastRunAt      *time.Time `json:"last_run_at"`
	LastRunStatus  string     `json:"last_run_status"`
	LastRunError   string     `json:"last_run_error"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type GetTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type SearchTasksReq struct {
	Name        string `json:"name"`
	TaskType    string `json:"task_type"`
	SchedulerID int64  `json:"scheduler_id"`
	Enabled     *bool  `json:"enabled"`
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
}

type AddTaskReq struct {
	Name           string `json:"name" binding:"required"`
	TaskType       string `json:"task_type" binding:"required"`
	SchedulerID    int64  `json:"scheduler_id" binding:"required"`
	CronExpression string `json:"cron_expression"`
	Enabled        *bool  `json:"enabled"`
}

type UpdateTaskReq struct {
	ID             int64   `json:"id" binding:"required"`
	Name           *string `json:"name"`
	TaskType       *string `json:"task_type"`
	SchedulerID    *int64  `json:"scheduler_id"`
	CronExpression *string `json:"cron_expression"`
	Enabled        *bool   `json:"enabled"`
}

type RunTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type StopTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type RestartTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetTaskRes struct {
	Task *Task `json:"task"`
}

type SearchTasksRes struct {
	Tasks []*Task `json:"tasks"`
	Total int64  `json:"total"`
}

type AddTaskRes struct {
	ID int64 `json:"id"`
}

type UpdateTaskRes struct {
	ID int64 `json:"id"`
}

type TaskActionRes struct {
	ID int64 `json:"id"`
}
