package scheduling

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type SchedulingDomain struct {
	DB *dbdao.DB
}

func NewSchedulingDomain(db *dbdao.DB) *SchedulingDomain {
	return &SchedulingDomain{DB: db}
}

func (d *SchedulingDomain) GetTask(ctx context.Context, id int64) (*Task, error) {
	t, err := d.DB.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Task{
		ID:             t.ID,
		Name:           t.Name,
		Description:    t.Description,
		TaskType:       t.TaskType,
		SchedulerID:    t.SchedulerID,
		CronExpression: t.CronExpression,
		Enabled:        t.Enabled,
		NextRunAt:      t.NextRunAt,
		LastRunAt:      t.LastRunAt,
		LastRunStatus:  t.LastRunStatus,
		LastRunError:   t.LastRunError,
		CreatedAt:      &t.CreatedAt,
		UpdatedAt:      &t.UpdatedAt,
	}, nil
}

func (d *SchedulingDomain) SearchTasks(ctx context.Context, name, taskType string, schedulerID int64, enabled *bool, page, limit int64) ([]*Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	ts, total, err := d.DB.SearchTasks(ctx, name, taskType, schedulerID, enabled, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Task, len(ts))
	for i, t := range ts {
		result[i] = &Task{
			ID: t.ID, Name: t.Name, Description: t.Description, TaskType: t.TaskType, SchedulerID: t.SchedulerID,
			CronExpression: t.CronExpression, Enabled: t.Enabled, NextRunAt: t.NextRunAt, LastRunAt: t.LastRunAt,
			LastRunStatus: t.LastRunStatus, LastRunError: t.LastRunError, CreatedAt: &t.CreatedAt, UpdatedAt: &t.UpdatedAt,
		}
	}
	return result, total, nil
}

func (d *SchedulingDomain) AddTask(ctx context.Context, name, taskType string, schedulerID int64, cronExpression string, enabled bool) (*Task, error) {
	if name == "" || taskType == "" {
		return nil, errors.New("missing required fields")
	}
	t := &dbdao.Task{Name: name, TaskType: taskType, SchedulerID: schedulerID, CronExpression: cronExpression, Enabled: enabled}
	created, err := d.DB.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return &Task{ID: created.ID, Name: created.Name, TaskType: created.TaskType, SchedulerID: created.SchedulerID, CronExpression: created.CronExpression, Enabled: created.Enabled, CreatedAt: &created.CreatedAt, UpdatedAt: &created.UpdatedAt}, nil
}

func (d *SchedulingDomain) UpdateTask(ctx context.Context, id int64, name, taskType, cronExpression *string, schedulerID *int64, enabled *bool) (*Task, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	t := &dbdao.Task{ID: id}
	if name != nil {
		t.Name = *name
	}
	if taskType != nil {
		t.TaskType = *taskType
	}
	if schedulerID != nil {
		t.SchedulerID = *schedulerID
	}
	if cronExpression != nil {
		t.CronExpression = *cronExpression
	}
	if enabled != nil {
		t.Enabled = *enabled
	}
	updated, err := d.DB.UpdateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return &Task{ID: updated.ID, Name: updated.Name, TaskType: updated.TaskType, SchedulerID: updated.SchedulerID, CronExpression: updated.CronExpression, Enabled: updated.Enabled, CreatedAt: &updated.CreatedAt, UpdatedAt: &updated.UpdatedAt}, nil
}

func (d *SchedulingDomain) RunTask(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, d.DB.RunTask(ctx, id)
}

func (d *SchedulingDomain) StopTask(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, d.DB.StopTask(ctx, id)
}

func (d *SchedulingDomain) RestartTask(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, d.DB.RestartTask(ctx, id)
}
