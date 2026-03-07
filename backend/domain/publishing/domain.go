package publishing

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type PublishingDomain struct {
	DB *dbdao.DB
}

func NewPublishingDomain(db *dbdao.DB) *PublishingDomain {
	return &PublishingDomain{DB: db}
}

func (d *PublishingDomain) GetPublisher(ctx context.Context, id int64) (*Publisher, error) {
	p, err := d.DB.GetPublisher(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		ID:         p.ID,
		Name:       p.Name,
		Platform:   p.Platform,
		PlatformID: p.PlatformID,
		Enabled:    p.Enabled,
		CreatedAt:  &p.CreatedAt,
		UpdatedAt:  &p.UpdatedAt,
	}, nil
}

func (d *PublishingDomain) SearchPublishers(ctx context.Context, platform string, enabled *bool, page, limit int64) ([]*Publisher, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	ps, total, err := d.DB.SearchPublishers(ctx, platform, enabled, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Publisher, len(ps))
	for i, p := range ps {
		result[i] = &Publisher{
			ID:         p.ID,
			Name:       p.Name,
			Platform:   p.Platform,
			PlatformID: p.PlatformID,
			Enabled:    p.Enabled,
			CreatedAt:  &p.CreatedAt,
			UpdatedAt:  &p.UpdatedAt,
		}
	}
	return result, total, nil
}

func (d *PublishingDomain) AddPublisher(ctx context.Context, name, platform, platformID string, enabled bool) (*Publisher, error) {
	if name == "" || platform == "" {
		return nil, errors.New("missing required fields")
	}
	p := &dbdao.Publisher{
		Name:       name,
		Platform:   platform,
		PlatformID: platformID,
		Enabled:    enabled,
	}
	created, err := d.DB.CreatePublisher(ctx, p)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		ID:         created.ID,
		Name:       created.Name,
		Platform:   created.Platform,
		PlatformID: created.PlatformID,
		Enabled:    created.Enabled,
		CreatedAt:  &created.CreatedAt,
		UpdatedAt:  &created.UpdatedAt,
	}, nil
}

func (d *PublishingDomain) UpdatePublisher(ctx context.Context, id int64, name, platform, platformID *string, enabled *bool) (*Publisher, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	p := &dbdao.Publisher{ID: id}
	if name != nil {
		p.Name = *name
	}
	if platform != nil {
		p.Platform = *platform
	}
	if platformID != nil {
		p.PlatformID = *platformID
	}
	if enabled != nil {
		p.Enabled = *enabled
	}
	updated, err := d.DB.UpdatePublisher(ctx, p)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		ID:         updated.ID,
		Name:       updated.Name,
		Platform:   updated.Platform,
		PlatformID: updated.PlatformID,
		Enabled:    updated.Enabled,
		CreatedAt:  &updated.CreatedAt,
		UpdatedAt:  &updated.UpdatedAt,
	}, nil
}

func (d *PublishingDomain) PublishContent(ctx context.Context, publisherID, contentID int64) (int64, error) {
	if publisherID <= 0 || contentID <= 0 {
		return 0, errors.New("invalid publisher_id or content_id")
	}
	task := &dbdao.PublishTask{
		PublisherID: publisherID,
		ContentID:   contentID,
		Status:      "pending",
	}
	created, err := d.DB.CreatePublishTask(ctx, task)
	if err != nil {
		return 0, err
	}
	return created.ID, nil
}

func (d *PublishingDomain) GetPublishTask(ctx context.Context, id int64) (*PublishTask, error) {
	t, err := d.DB.GetPublishTask(ctx, id)
	if err != nil {
		return nil, err
	}
	return &PublishTask{
		ID:          t.ID,
		PublisherID: t.PublisherID,
		ContentID:   t.ContentID,
		Content:     t.Content,
		Status:      t.Status,
		Error:       t.Error,
		CreatedAt:   &t.CreatedAt,
		StartedAt:   t.StartedAt,
		CompletedAt: t.CompletedAt,
	}, nil
}

func (d *PublishingDomain) SearchPublishTasks(ctx context.Context, status, platform, createdAt string, page, limit int64) ([]*PublishTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	ts, total, err := d.DB.SearchPublishTasks(ctx, status, platform, createdAt, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*PublishTask, len(ts))
	for i, t := range ts {
		result[i] = &PublishTask{
			ID:          t.ID,
			PublisherID: t.PublisherID,
			ContentID:   t.ContentID,
			Content:     t.Content,
			Status:      t.Status,
			Error:       t.Error,
			CreatedAt:   &t.CreatedAt,
			StartedAt:   t.StartedAt,
			CompletedAt: t.CompletedAt,
		}
	}
	return result, total, nil
}

func (d *PublishingDomain) GetPublishLog(ctx context.Context, id int64) (*PublishLog, error) {
	l, err := d.DB.GetPublishLog(ctx, id)
	if err != nil {
		return nil, err
	}
	return &PublishLog{
		ID:            l.ID,
		PublishTaskID: l.PublishTaskID,
		LogType:       l.LogType,
		Message:       l.Message,
		CreatedAt:     &l.CreatedAt,
	}, nil
}

func (d *PublishingDomain) SearchPublishLogs(ctx context.Context, publishTaskID int64, logType string, page, limit int64) ([]*PublishLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	ls, total, err := d.DB.SearchPublishLogs(ctx, publishTaskID, logType, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*PublishLog, len(ls))
	for i, l := range ls {
		result[i] = &PublishLog{
			ID:            l.ID,
			PublishTaskID: l.PublishTaskID,
			LogType:       l.LogType,
			Message:       l.Message,
			CreatedAt:     &l.CreatedAt,
		}
	}
	return result, total, nil
}
