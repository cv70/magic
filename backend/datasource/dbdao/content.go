package dbdao

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ContentModel 数据库模型
type ContentModel struct {
	ID          int64     `gorm:"primarykey;column:id"`
	Title       string    `gorm:"column:title;type:varchar(255);not null"`
	Body        string    `gorm:"column:body;type:longtext"`
	ContentType string    `gorm:"column:content_type;type:varchar(100);not null"`
	Status      string    `gorm:"column:status;type:varchar(50)"`
	Tags        string    `gorm:"column:tags;type:json"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
	PublishedAt *time.Time `gorm:"column:published_at"`
}

// TableName 指定表名
func (ContentModel) TableName() string {
	return "contents"
}

// Content DAO层返回的内容数据结构
type Content struct {
	ID          int64
	Title       string
	Body        string
	ContentType string
	Status      string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}

// GetContent 获取内容
func (d *DB) GetContent(ctx context.Context, id int64) (*Content, error) {
	var model ContentModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("content not found")
		}
		return nil, result.Error
	}

	var tags []string
	if len(model.Tags) > 0 {
		if err := json.Unmarshal([]byte(model.Tags), &tags); err != nil {
			tags = []string{}
		}
	}

	return &Content{
		ID:          model.ID,
		Title:       model.Title,
		Body:        model.Body,
		ContentType: model.ContentType,
		Status:      model.Status,
		Tags:        tags,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		PublishedAt: model.PublishedAt,
	}, nil
}

// SearchContents 搜索内容
func (d *DB) SearchContents(
	ctx context.Context,
	query, contentType, status, tag string,
	page, limit int64,
) ([]*Content, int64, error) {
	var models []ContentModel
	q := d.DB().WithContext(ctx)

	if query != "" {
		q = q.Where("title LIKE ? OR body LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if contentType != "" {
		q = q.Where("content_type = ?", contentType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if tag != "" {
		q = q.Where("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
	}

	var total int64
	q.Model(&ContentModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	contents := make([]*Content, len(models))
	for i, model := range models {
		var tags []string
		if len(model.Tags) > 0 {
			if err := json.Unmarshal([]byte(model.Tags), &tags); err != nil {
				tags = []string{}
			}
		}

		contents[i] = &Content{
			ID:          model.ID,
			Title:       model.Title,
			Body:        model.Body,
			ContentType: model.ContentType,
			Status:      model.Status,
			Tags:        tags,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
			PublishedAt: model.PublishedAt,
		}
	}

	return contents, total, nil
}

// CreateContent 创建内容
func (d *DB) CreateContent(ctx context.Context, content *Content) (*Content, error) {
	tagsJSON, err := json.Marshal(content.Tags)
	if err != nil {
		tagsJSON = []byte("[]")
	}

	model := &ContentModel{
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        string(tagsJSON),
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &Content{
		ID:          model.ID,
		Title:       model.Title,
		Body:        model.Body,
		ContentType: model.ContentType,
		Status:      model.Status,
		Tags:        content.Tags,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		PublishedAt: model.PublishedAt,
	}, nil
}

// UpdateContent 更新内容
func (d *DB) UpdateContent(ctx context.Context, content *Content) (*Content, error) {
	updates := map[string]interface{}{}

	if content.Title != "" {
		updates["title"] = content.Title
	}
	if content.Body != "" {
		updates["body"] = content.Body
	}
	if content.ContentType != "" {
		updates["content_type"] = content.ContentType
	}
	if content.Status != "" {
		updates["status"] = content.Status
	}
	if len(content.Tags) > 0 {
		tagsJSON, _ := json.Marshal(content.Tags)
		updates["tags"] = string(tagsJSON)
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := d.DB().WithContext(ctx).Model(&ContentModel{}).Where("id = ?", content.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("content not found")
	}

	return d.GetContent(ctx, content.ID)
}

// DeleteContent 删除内容
func (d *DB) DeleteContent(ctx context.Context, id int64) error {
	result := d.DB().WithContext(ctx).Where("id = ?", id).Delete(&ContentModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("content not found")
	}

	return nil
}
