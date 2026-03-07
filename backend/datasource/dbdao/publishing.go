package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PublisherModel struct {
	ID         int64     `gorm:"primarykey;column:id"`
	Name       string    `gorm:"column:name;type:varchar(255);not null"`
	Platform   string    `gorm:"column:platform;type:varchar(100);not null;index"`
	PlatformID string    `gorm:"column:platform_id;type:varchar(255)"`
	Enabled    bool      `gorm:"column:enabled;default:true"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (PublisherModel) TableName() string {
	return "publishers"
}

type Publisher struct {
	ID         int64
	Name       string
	Platform   string
	PlatformID string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (d *DB) GetPublisher(ctx context.Context, id int64) (*Publisher, error) {
	var m PublisherModel
	r := d.DB().WithContext(ctx).Where("id = ?", id).First(&m)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("publisher not found")
		}
		return nil, r.Error
	}
	return &Publisher{ID: m.ID, Name: m.Name, Platform: m.Platform, PlatformID: m.PlatformID, Enabled: m.Enabled, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (d *DB) SearchPublishers(ctx context.Context, platform string, enabled *bool, page, limit int64) ([]*Publisher, int64, error) {
	var ms []PublisherModel
	q := d.DB().WithContext(ctx)
	if platform != "" {
		q = q.Where("platform = ?", platform)
	}
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}
	var total int64
	q.Model(&PublisherModel{}).Count(&total)
	r := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&ms)
	if r.Error != nil {
		return nil, 0, r.Error
	}
	ps := make([]*Publisher, len(ms))
	for i, m := range ms {
		ps[i] = &Publisher{ID: m.ID, Name: m.Name, Platform: m.Platform, PlatformID: m.PlatformID, Enabled: m.Enabled, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
	}
	return ps, total, nil
}

func (d *DB) CreatePublisher(ctx context.Context, p *Publisher) (*Publisher, error) {
	m := &PublisherModel{Name: p.Name, Platform: p.Platform, PlatformID: p.PlatformID, Enabled: p.Enabled}
	r := d.DB().WithContext(ctx).Create(m)
	if r.Error != nil {
		return nil, r.Error
	}
	return &Publisher{ID: m.ID, Name: m.Name, Platform: m.Platform, PlatformID: m.PlatformID, Enabled: m.Enabled, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (d *DB) UpdatePublisher(ctx context.Context, p *Publisher) (*Publisher, error) {
	up := map[string]interface{}{}
	if p.Name != "" {
		up["name"] = p.Name
	}
	if p.Platform != "" {
		up["platform"] = p.Platform
	}
	if p.PlatformID != "" {
		up["platform_id"] = p.PlatformID
	}
	up["enabled"] = p.Enabled
	r := d.DB().WithContext(ctx).Model(&PublisherModel{}).Where("id = ?", p.ID).Updates(up)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected == 0 {
		return nil, errors.New("publisher not found")
	}
	return d.GetPublisher(ctx, p.ID)
}

type PublishTaskModel struct {
	ID          int64      `gorm:"primarykey;column:id"`
	PublisherID int64      `gorm:"column:publisher_id;index"`
	ContentID   int64      `gorm:"column:content_id;index"`
	Content     string     `gorm:"column:content;type:longtext"`
	Status      string     `gorm:"column:status;type:varchar(50)"`
	Error       string     `gorm:"column:error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime:milli"`
	StartedAt   *time.Time `gorm:"column:started_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

func (PublishTaskModel) TableName() string {
	return "publish_tasks"
}

type PublishTask struct {
	ID          int64
	PublisherID int64
	ContentID   int64
	Content     string
	Status      string
	Error       string
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

func (d *DB) CreatePublishTask(ctx context.Context, t *PublishTask) (*PublishTask, error) {
	m := &PublishTaskModel{PublisherID: t.PublisherID, ContentID: t.ContentID, Content: t.Content, Status: t.Status}
	r := d.DB().WithContext(ctx).Create(m)
	if r.Error != nil {
		return nil, r.Error
	}
	return &PublishTask{ID: m.ID, PublisherID: m.PublisherID, ContentID: m.ContentID, Content: m.Content, Status: m.Status, CreatedAt: m.CreatedAt}, nil
}

func (d *DB) GetPublishTask(ctx context.Context, id int64) (*PublishTask, error) {
	var m PublishTaskModel
	r := d.DB().WithContext(ctx).Where("id = ?", id).First(&m)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, r.Error
	}
	return &PublishTask{ID: m.ID, PublisherID: m.PublisherID, ContentID: m.ContentID, Content: m.Content, Status: m.Status, Error: m.Error, CreatedAt: m.CreatedAt, StartedAt: m.StartedAt, CompletedAt: m.CompletedAt}, nil
}

func (d *DB) SearchPublishTasks(ctx context.Context, status, platform, createdAt string, page, limit int64) ([]*PublishTask, int64, error) {
	var ms []PublishTaskModel
	q := d.DB().WithContext(ctx)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	q.Model(&PublishTaskModel{}).Count(&total)
	r := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&ms)
	if r.Error != nil {
		return nil, 0, r.Error
	}
	ts := make([]*PublishTask, len(ms))
	for i, m := range ms {
		ts[i] = &PublishTask{ID: m.ID, PublisherID: m.PublisherID, ContentID: m.ContentID, Content: m.Content, Status: m.Status, Error: m.Error, CreatedAt: m.CreatedAt, StartedAt: m.StartedAt, CompletedAt: m.CompletedAt}
	}
	return ts, total, nil
}

type PublishLogModel struct {
	ID            int64     `gorm:"primarykey;column:id"`
	PublishTaskID int64     `gorm:"column:publish_task_id;index"`
	LogType       string    `gorm:"column:log_type;type:varchar(50)"`
	Message       string    `gorm:"column:message;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (PublishLogModel) TableName() string {
	return "publish_logs"
}

type PublishLog struct {
	ID            int64
	PublishTaskID int64
	LogType       string
	Message       string
	CreatedAt     time.Time
}

func (d *DB) GetPublishLog(ctx context.Context, id int64) (*PublishLog, error) {
	var m PublishLogModel
	r := d.DB().WithContext(ctx).Where("id = ?", id).First(&m)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("log not found")
		}
		return nil, r.Error
	}
	return &PublishLog{ID: m.ID, PublishTaskID: m.PublishTaskID, LogType: m.LogType, Message: m.Message, CreatedAt: m.CreatedAt}, nil
}

func (d *DB) SearchPublishLogs(ctx context.Context, publishTaskID int64, logType string, page, limit int64) ([]*PublishLog, int64, error) {
	var ms []PublishLogModel
	q := d.DB().WithContext(ctx).Where("publish_task_id = ?", publishTaskID)
	if logType != "" {
		q = q.Where("log_type = ?", logType)
	}
	var total int64
	q.Model(&PublishLogModel{}).Count(&total)
	r := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&ms)
	if r.Error != nil {
		return nil, 0, r.Error
	}
	ls := make([]*PublishLog, len(ms))
	for i, m := range ms {
		ls[i] = &PublishLog{ID: m.ID, PublishTaskID: m.PublishTaskID, LogType: m.LogType, Message: m.Message, CreatedAt: m.CreatedAt}
	}
	return ls, total, nil
}
