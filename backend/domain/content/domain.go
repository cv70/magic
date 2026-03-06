package content

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type ContentDomain struct {
	DB *dbdao.DB
}

func NewContentDomain(db *dbdao.DB) *ContentDomain {
	return &ContentDomain{DB: db}
}

// GetContent 获取内容详情
func (d *ContentDomain) GetContent(ctx context.Context, id int64) (*Content, error) {
	content, err := d.DB.GetContent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Content{
		ID:          content.ID,
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        content.Tags,
		CreatedAt:   &content.CreatedAt,
		UpdatedAt:   &content.UpdatedAt,
		PublishedAt: content.PublishedAt,
	}, nil
}

// SearchContents 搜索内容
func (d *ContentDomain) SearchContents(
	ctx context.Context,
	query, contentType, status, tag string,
	page, limit int64,
) ([]*Content, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	contents, total, err := d.DB.SearchContents(ctx, query, contentType, status, tag, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*Content, len(contents))
	for i, c := range contents {
		result[i] = &Content{
			ID:          c.ID,
			Title:       c.Title,
			Body:        c.Body,
			ContentType: c.ContentType,
			Status:      c.Status,
			Tags:        c.Tags,
			CreatedAt:   &c.CreatedAt,
			UpdatedAt:   &c.UpdatedAt,
			PublishedAt: c.PublishedAt,
		}
	}

	return result, total, nil
}

// AddContent 添加内容
func (d *ContentDomain) AddContent(
	ctx context.Context,
	title, body, contentType, status string,
	tags []string,
) (*Content, error) {
	if title == "" || body == "" || contentType == "" {
		return nil, errors.New("missing required fields")
	}

	content := &dbdao.Content{
		Title:       title,
		Body:        body,
		ContentType: contentType,
		Status:      status,
		Tags:        tags,
	}

	created, err := d.DB.CreateContent(ctx, content)
	if err != nil {
		return nil, err
	}

	return &Content{
		ID:          created.ID,
		Title:       created.Title,
		Body:        created.Body,
		ContentType: created.ContentType,
		Status:      created.Status,
		Tags:        created.Tags,
		CreatedAt:   &created.CreatedAt,
		UpdatedAt:   &created.UpdatedAt,
		PublishedAt: created.PublishedAt,
	}, nil
}

// UpdateContent 更新内容
func (d *ContentDomain) UpdateContent(
	ctx context.Context,
	id int64,
	title, body, contentType, status *string,
	tags []string,
) (*Content, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	content := &dbdao.Content{
		ID:   id,
		Tags: tags,
	}

	if title != nil {
		content.Title = *title
	}
	if body != nil {
		content.Body = *body
	}
	if contentType != nil {
		content.ContentType = *contentType
	}
	if status != nil {
		content.Status = *status
	}

	updated, err := d.DB.UpdateContent(ctx, content)
	if err != nil {
		return nil, err
	}

	return &Content{
		ID:          updated.ID,
		Title:       updated.Title,
		Body:        updated.Body,
		ContentType: updated.ContentType,
		Status:      updated.Status,
		Tags:        updated.Tags,
		CreatedAt:   &updated.CreatedAt,
		UpdatedAt:   &updated.UpdatedAt,
		PublishedAt: updated.PublishedAt,
	}, nil
}

// DeleteContent 删除内容
func (d *ContentDomain) DeleteContent(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return d.DB.DeleteContent(ctx, id)
}
