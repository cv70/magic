package publishing

import "time"

type Publisher struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Platform   string     `json:"platform"`
	PlatformID string     `json:"platform_id"`
	Enabled    bool       `json:"enabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type PublishTask struct {
	ID          int64      `json:"id"`
	PublisherID int64      `json:"publisher_id"`
	ContentID   int64      `json:"content_id"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	Error       string     `json:"error"`
	CreatedAt   *time.Time `json:"created_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type PublishLog struct {
	ID            int64      `json:"id"`
	PublishTaskID int64      `json:"publish_task_id"`
	LogType       string     `json:"log_type"`
	Message       string     `json:"message"`
	CreatedAt     *time.Time `json:"created_at"`
}

type GetPublisherReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetPublisherRes struct {
	Publisher *Publisher `json:"publisher"`
}

type SearchPublishersReq struct {
	Platform string `json:"platform"`
	Enabled  *bool  `json:"enabled"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type SearchPublishersRes struct {
	Publishers []*Publisher `json:"publishers"`
	Total      int64        `json:"total"`
}

type AddPublisherReq struct {
	Name       string `json:"name" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	PlatformID string `json:"platform_id"`
	Enabled    *bool  `json:"enabled"`
}

type AddPublisherRes struct {
	ID int64 `json:"id"`
}

type UpdatePublisherReq struct {
	ID         int64   `json:"id" binding:"required"`
	Name       *string `json:"name"`
	Platform   *string `json:"platform"`
	PlatformID *string `json:"platform_id"`
	Enabled    *bool   `json:"enabled"`
}

type UpdatePublisherRes struct {
	ID int64 `json:"id"`
}

type PublishContentReq struct {
	PublisherID int64 `json:"publisher_id" binding:"required"`
	ContentID   int64 `json:"content_id" binding:"required"`
}

type PublishContentRes struct {
	ID int64 `json:"id"`
}

type GetPublishTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetPublishTaskRes struct {
	PublishTask *PublishTask `json:"publish_task"`
}

type SearchPublishTasksReq struct {
	Status    string `json:"status"`
	Platform  string `json:"platform"`
	CreatedAt string `json:"created_at"`
	Page      int64  `json:"page"`
	Limit     int64  `json:"limit"`
}

type SearchPublishTasksRes struct {
	PublishTasks []*PublishTask `json:"publish_tasks"`
	Total        int64          `json:"total"`
}

type GetPublishLogReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetPublishLogRes struct {
	PublishLog *PublishLog `json:"publish_log"`
}

type SearchPublishLogsReq struct {
	PublishTaskID int64  `json:"publish_task_id" binding:"required"`
	LogType       string `json:"log_type"`
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
}

type SearchPublishLogsRes struct {
	PublishLogs []*PublishLog `json:"publish_logs"`
	Total       int64         `json:"total"`
}

// ===================== Analytics API =====================

type AnalyticsSummaryReq struct {
	Days int `json:"days"`
}

type AnalyticsSummaryRes struct {
	TotalPublished    int64 `json:"total_published"`
	TotalViews        int64 `json:"total_views"`
	TotalLikes        int64 `json:"total_likes"`
	TotalComments     int64 `json:"total_comments"`
	AvgLikesPerPost   int64 `json:"avg_likes_per_post"`
	TotalNewFollowers int64 `json:"total_new_followers"`
}

type ContentRankingItem struct {
	ContentID int64  `json:"content_id"`
	Title     string `json:"title"`
	Views     int64  `json:"views"`
	Likes     int64  `json:"likes"`
	Comments  int64  `json:"comments"`
}

type AnalyticsRankingReq struct {
	Metric string `json:"metric"`
	Days   int    `json:"days"`
	Limit  int    `json:"limit"`
}

type AnalyticsRankingRes struct {
	Data []*ContentRankingItem `json:"data"`
}

type PlatformMetrics struct {
	Platform  string `json:"platform"`
	Count     int64  `json:"count"`
	Success   int64  `json:"success"`
	Failed    int64  `json:"failed"`
	SuccessRate float64 `json:"success_rate"`
}

type AnalyticsPlatformComparisonRes struct {
	Data []*PlatformMetrics `json:"data"`
}

