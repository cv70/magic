package content

import "time"

// Content 内容实体
type Content struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	ContentType string     `json:"content_type"`
	Status      string     `json:"status"`
	Tags        []string   `json:"tags"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// API 请求/响应

type GetContentReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetContentRes struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	ContentType string     `json:"content_type"`
	Status      string     `json:"status"`
	Tags        []string   `json:"tags"`
	CreatedAt   *time.Time `json:"created_at"`
}

type SearchContentReq struct {
	Query       string `json:"query"`
	ContentType string `json:"content_type"`
	Status      string `json:"status"`
	Tag         string `json:"tag"`
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
}

type SearchContentRes struct {
	Contents []*GetContentRes `json:"contents"`
	Total    int64            `json:"total"`
}

type AddContentReq struct {
	Title       string   `json:"title" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	ContentType string   `json:"content_type" binding:"required"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
}

type UpdateContentReq struct {
	ID          int64    `json:"id" binding:"required"`
	Title       *string  `json:"title"`
	Body        *string  `json:"body"`
	ContentType *string  `json:"content_type"`
	Status      *string  `json:"status"`
	Tags        []string `json:"tags"`
}

type DeleteContentReq struct {
	ID int64 `json:"id" binding:"required"`
}
