package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TaskModel struct {
	ID             int64      `gorm:"primarykey;column:id"`
	Name           string     `gorm:"column:name;type:varchar(255);not null"`
	Description    string     `gorm:"column:description;type:varchar(500)"`
	TaskType       string     `gorm:"column:task_type;type:varchar(100);not null"`
	SchedulerID    int64      `gorm:"column:scheduler_id;index"`
	CronExpression string     `gorm:"column:cron_expression;type:varchar(100)"`
	Enabled        bool       `gorm:"column:enabled;default:true"`
	NextRunAt      *time.Time `gorm:"column:next_run_at"`
	LastRunAt      *time.Time `gorm:"column:last_run_at"`
	LastRunStatus  string     `gorm:"column:last_run_status;type:varchar(50)"`
	LastRunError   string     `gorm:"column:last_run_error;type:text"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (TaskModel) TableName() string {
	return "scheduling_tasks"
}

type Task struct {
	ID             int64
	Name           string
	Description    string
	TaskType       string
	SchedulerID    int64
	CronExpression string
	Enabled        bool
	NextRunAt      *time.Time
	LastRunAt      *time.Time
	LastRunStatus  string
	LastRunError   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (d *DB) GetTask(ctx context.Context, id int64) (*Task, error) {
	var model TaskModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, result.Error
	}

	return &Task{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description,
		TaskType:       model.TaskType,
		SchedulerID:    model.SchedulerID,
		CronExpression: model.CronExpression,
		Enabled:        model.Enabled,
		NextRunAt:      model.NextRunAt,
		LastRunAt:      model.LastRunAt,
		LastRunStatus:  model.LastRunStatus,
		LastRunError:   model.LastRunError,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}, nil
}

func (d *DB) SearchTasks(
	ctx context.Context,
	name, taskType string,
	schedulerID int64,
	enabled *bool,
	page, limit int64,
) ([]*Task, int64, error) {
	var models []TaskModel
	q := d.DB().WithContext(ctx)

	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if taskType != "" {
		q = q.Where("task_type = ?", taskType)
	}
	if schedulerID > 0 {
		q = q.Where("scheduler_id = ?", schedulerID)
	}
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}

	var total int64
	q.Model(&TaskModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	tasks := make([]*Task, len(models))
	for i, m := range models {
		tasks[i] = &Task{
			ID:             m.ID,
			Name:           m.Name,
			Description:    m.Description,
			TaskType:       m.TaskType,
			SchedulerID:    m.SchedulerID,
			CronExpression: m.CronExpression,
			Enabled:        m.Enabled,
			NextRunAt:      m.NextRunAt,
			LastRunAt:      m.LastRunAt,
			LastRunStatus:  m.LastRunStatus,
			LastRunError:   m.LastRunError,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		}
	}

	return tasks, total, nil
}

func (d *DB) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	model := &TaskModel{
		Name:           task.Name,
		Description:    task.Description,
		TaskType:       task.TaskType,
		SchedulerID:    task.SchedulerID,
		CronExpression: task.CronExpression,
		Enabled:        task.Enabled,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &Task{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description,
		TaskType:       model.TaskType,
		SchedulerID:    model.SchedulerID,
		CronExpression: model.CronExpression,
		Enabled:        model.Enabled,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}, nil
}

func (d *DB) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	updates := map[string]interface{}{}

	if task.Name != "" {
		updates["name"] = task.Name
	}
	if task.Description != "" {
		updates["description"] = task.Description
	}
	if task.TaskType != "" {
		updates["task_type"] = task.TaskType
	}
	if task.SchedulerID > 0 {
		updates["scheduler_id"] = task.SchedulerID
	}
	if task.CronExpression != "" {
		updates["cron_expression"] = task.CronExpression
	}

	result := d.DB().WithContext(ctx).Model(&TaskModel{}).Where("id = ?", task.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("task not found")
	}

	return d.GetTask(ctx, task.ID)
}

func (d *DB) RunTask(ctx context.Context, id int64) error {
	now := time.Now()
	result := d.DB().WithContext(ctx).Model(&TaskModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_run_status": "running",
		"last_run_at":     now,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (d *DB) StopTask(ctx context.Context, id int64) error {
	result := d.DB().WithContext(ctx).Model(&TaskModel{}).Where("id = ?", id).Update("last_run_status", "stopped")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (d *DB) RestartTask(ctx context.Context, id int64) error {
	now := time.Now()
	result := d.DB().WithContext(ctx).Model(&TaskModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_run_status": "pending",
		"last_run_at":     now,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}
