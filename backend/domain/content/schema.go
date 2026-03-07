package content

import "time"

// ===================== DOMAIN ENTITIES =====================

// Draft 内容草稿
type Draft struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	UserID        int64      `json:"user_id" gorm:"index"`
	Title         string     `json:"title" gorm:"size:255"`
	Content       string     `json:"content" gorm:"type:longtext"`
	ContentType   string     `json:"content_type" gorm:"size:20"` // text/code/image等
	Status        string     `json:"status" gorm:"size:20"`       // draft/archived
	Tags          string     `json:"tags" gorm:"type:json"`       // JSON数组
	Metadata      string     `json:"metadata" gorm:"type:json"`   // 自定义元数据
	LastEditedAt  *time.Time `json:"last_edited_at"`
	SavedVersionID *int64    `json:"saved_version_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// Content 正式发布的内容
type Content struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	UserID      int64      `json:"user_id" gorm:"index"`
	Title       string     `json:"title" gorm:"size:255"`
	Body        string     `json:"body" gorm:"type:longtext"`
	ContentType string     `json:"content_type" gorm:"size:20"`
	Status      string     `json:"status" gorm:"size:20"`
	Tags        string     `json:"tags" gorm:"type:json"`       // JSON数组
	Metadata    string     `json:"metadata" gorm:"type:json"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// ContentVersion 版本历史
type ContentVersion struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	ContentID     int64      `json:"content_id" gorm:"index"`
	DraftID       *int64     `json:"draft_id"`
	VersionNum    int        `json:"version_num"`
	Title         string     `json:"title" gorm:"size:255"`
	Content       string     `json:"content" gorm:"type:longtext"`
	ChangeSummary string     `json:"change_summary" gorm:"size:500"`
	CreatedBy     int64      `json:"created_by"`
	CreatedAt     *time.Time `json:"created_at"`
}

// Tag 标签
type Tag struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	UserID    int64      `json:"user_id" gorm:"index"`
	Name      string     `json:"name" gorm:"size:50"`
	Color     string     `json:"color" gorm:"size:20"`
	Count     int        `json:"count"` // 使用次数
	CreatedAt *time.Time `json:"created_at"`
}

// Category 分类
type Category struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	UserID    int64      `json:"user_id" gorm:"index"`
	Name      string     `json:"name" gorm:"size:100"`
	Icon      string     `json:"icon" gorm:"size:50"`
	Color     string     `json:"color" gorm:"size:20"`
	Order     int        `json:"order" gorm:"column:order_num"`
	CreatedAt *time.Time `json:"created_at"`
}

// ===================== PUBLISH RELATED =====================

// PublishTask 发布任务
type PublishTask struct {
	ID              int64       `json:"id" gorm:"primaryKey"`
	UserID          int64       `json:"user_id" gorm:"index"`
	ContentID       int64       `json:"content_id" gorm:"index"`
	DraftID         *int64      `json:"draft_id"`
	TargetPlatforms string      `json:"target_platforms" gorm:"type:json"` // ["juejin", "csdn"]
	Status          string      `json:"status" gorm:"size:20"`             // pending/publishing/published/failed
	ScheduledAt     *time.Time  `json:"scheduled_at" gorm:"index"`
	PublishedAt     *time.Time  `json:"published_at"`
	RetryCount      int         `json:"retry_count"`
	MaxRetries      int         `json:"max_retries"`
	ErrorMessage    string      `json:"error_message" gorm:"type:text"`
	PlatformResults string      `json:"platform_results" gorm:"type:json"`
	CreatedAt       *time.Time  `json:"created_at"`
	UpdatedAt       *time.Time  `json:"updated_at"`
}

// PlatformResult 单个平台的发布结果
type PlatformResult struct {
	Platform    string `json:"platform"`
	Status      string `json:"status"`
	PostID      string `json:"post_id"`
	PostURL     string `json:"post_url"`
	Error       string `json:"error,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Metrics     Metrics `json:"metrics,omitempty"`
}

// Metrics 发布效果数据
type Metrics struct {
	Views    int       `json:"views"`
	Likes    int       `json:"likes"`
	Comments int       `json:"comments"`
	Shares   int       `json:"shares"`
	Follows  int       `json:"follows"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// PublishAnalytic 发布数据统计
type PublishAnalytic struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	PublishTaskID int64      `json:"publish_task_id" gorm:"index"`
	ContentID     int64      `json:"content_id" gorm:"index"`
	UserID        int64      `json:"user_id" gorm:"index"`
	Platform      string     `json:"platform" gorm:"size:50"`
	PostID        string     `json:"post_id" gorm:"size:255"`
	Title         string     `json:"title" gorm:"size:255"`
	PostURL       string     `json:"post_url" gorm:"type:text"`
	Metrics       string     `json:"metrics" gorm:"type:json"`
	LastSyncedAt  *time.Time `json:"last_synced_at"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// ===================== AI TEMPLATES =====================

// PromptTemplate AI提示词模板
type PromptTemplate struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	UserID        *int64     `json:"user_id"`     // NULL表示系统模板
	Category      string     `json:"category" gorm:"size:100"`
	SubCategory   string     `json:"sub_category" gorm:"size:100"`
	Name          string     `json:"name" gorm:"size:255"`
	Description   string     `json:"description" gorm:"type:text"`
	PromptText    string     `json:"prompt_text" gorm:"type:longtext"`
	Variables     string     `json:"variables" gorm:"type:json"`       // ["{keyword}", "{tone}"]
	ExampleInput  string     `json:"example_input" gorm:"type:json"`
	ExampleOutput string     `json:"example_output" gorm:"type:text"`
	UsageCount    int        `json:"usage_count"`
	Rating        float32    `json:"rating"`
	Tags          string     `json:"tags" gorm:"type:json"`
	Status        string     `json:"status" gorm:"size:20"` // active/archived
	Version       int        `json:"version"`
	CreatedBy     int64      `json:"created_by"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// ===================== DRAFT API DTOs =====================

// GetDraftReq 获取草稿请求
type GetDraftReq struct {
	ID int64 `json:"id" binding:"required"`
}

// GetDraftRes 获取草稿响应
type GetDraftRes struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	ContentType  string     `json:"content_type"`
	Status       string     `json:"status"`
	Tags         []string   `json:"tags"`
	LastEditedAt *time.Time `json:"last_edited_at"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

// SearchDraftReq 搜索草稿请求
type SearchDraftReq struct {
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
	Tags    string `json:"tags"`     // 逗号分割
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
}

// SearchDraftRes 搜索草稿响应
type SearchDraftRes struct {
	Data  []*GetDraftRes `json:"data"`
	Total int64          `json:"total"`
}

// CreateDraftReq 创建草稿请求
type CreateDraftReq struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
	Tags        []string `json:"tags"`
}

// CreateDraftRes 创建草稿响应
type CreateDraftRes struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at"`
}

// UpdateDraftReq 更新草稿请求
type UpdateDraftReq struct {
	ID            int64    `json:"id" binding:"required"`
	Title         *string  `json:"title"`
	Content       *string  `json:"content"`
	ContentType   *string  `json:"content_type"`
	Tags          []string `json:"tags"`
	ChangeSummary string   `json:"change_summary"` // 本次修改摘要
}

// UpdateDraftRes 更新草稿响应
type UpdateDraftRes struct {
	ID        int64      `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
	VersionNum int       `json:"version_num"`
}

// DeleteDraftReq 删除草稿请求
type DeleteDraftReq struct {
	ID int64 `json:"id" binding:"required"`
}

// PublishDraftReq 草稿发布为正式内容请求
type PublishDraftReq struct {
	ID     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"` // published/archived
}

// GetDraftVersionsReq 获取版本历史请求
type GetDraftVersionsReq struct {
	ID    int64 `json:"id" binding:"required"`
	Limit int64 `json:"limit"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	VersionNum    int        `json:"version_num"`
	Title         string     `json:"title"`
	ChangeSummary string     `json:"change_summary"`
	CreatedAt     *time.Time `json:"created_at"`
	CreatedByUser string     `json:"created_by_user"` // 修改人名称
}

// GetDraftVersionsRes 获取版本历史响应
type GetDraftVersionsRes struct {
	Data []*VersionInfo `json:"data"`
}

// RevertDraftVersionReq 恢复版本请求
type RevertDraftVersionReq struct {
	ID         int64 `json:"id" binding:"required"`
	VersionNum int   `json:"version_num" binding:"required"`
}

// RevertDraftVersionRes 恢复版本响应
type RevertDraftVersionRes struct {
	ID         int64      `json:"id"`
	VersionNum int        `json:"version_num"` // 新创建的版本号
	RevertedFrom int      `json:"reverted_from"` // 从哪个版本恢复
	UpdatedAt  *time.Time `json:"updated_at"`
}

// DiffDraftVersionReq 版本对比请求
type DiffDraftVersionReq struct {
	ID          int64 `json:"id" binding:"required"`
	FromVersion int   `json:"from_version" binding:"required"`
	ToVersion   int   `json:"to_version" binding:"required"`
}

// DiffItem 差异项
type DiffItem struct {
	Type     string `json:"type"`     // added/removed/modified
	Position string `json:"position"` // 位置描述
	Old      string `json:"old"`
	New      string `json:"new"`
}

// ContentDiff 内容差异
type ContentDiff struct {
	Title   *DiffItem `json:"title,omitempty"`
	Content []*DiffItem `json:"content,omitempty"`
}

// DiffDraftVersionRes 版本对比响应
type DiffDraftVersionRes struct {
	FromVersion int          `json:"from_version"`
	ToVersion   int          `json:"to_version"`
	Diff        *ContentDiff `json:"diff"`
}

// ===================== CONTENT API DTOs =====================

// GetContentReq 获取内容请求
type GetContentReq struct {
	ID int64 `json:"id" binding:"required"`
}

// GetContentRes 获取内容响应
type GetContentRes struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	ContentType string     `json:"content_type"`
	Status      string     `json:"status"`
	Tags        []string   `json:"tags"`
	CreatedAt   *time.Time `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// SearchContentReq 搜索内容请求
type SearchContentReq struct {
	Query       string `json:"query"`
	ContentType string `json:"content_type"`
	Status      string `json:"status"`
	Tag         string `json:"tag"`
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
}

// SearchContentRes 搜索内容响应
type SearchContentRes struct {
	Contents []*GetContentRes `json:"contents"`
	Total    int64            `json:"total"`
}

// AddContentReq 创建内容请求
type AddContentReq struct {
	Title       string   `json:"title" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	ContentType string   `json:"content_type" binding:"required"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
}

// UpdateContentReq 更新内容请求
type UpdateContentReq struct {
	ID          int64    `json:"id" binding:"required"`
	Title       *string  `json:"title"`
	Body        *string  `json:"body"`
	ContentType *string  `json:"content_type"`
	Status      *string  `json:"status"`
	Tags        []string `json:"tags"`
}

// DeleteContentReq 删除内容请求
type DeleteContentReq struct {
	ID int64 `json:"id" binding:"required"`
}
