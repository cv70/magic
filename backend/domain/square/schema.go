package square

import "time"

// ===================== DOMAIN ENTITIES =====================

// SquarePost 广场内容（从草稿发布）
type SquarePost struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	DraftID       int64      `json:"draft_id" gorm:"index"`                              // 关联草稿
	UserID        int64      `json:"user_id" gorm:"index"`                              // 发布者
	Title         string     `json:"title" gorm:"size:255;index"`                       // 标题
	PreviewText   string     `json:"preview_text" gorm:"type:text"`                     // 预览文本
	Domain        string     `json:"domain" gorm:"size:50;index"`                       // 领域分类
	Tags          string     `json:"tags" gorm:"type:json"`                             // JSON数组
	CoverImage    string     `json:"cover_image" gorm:"size:500"`                       // 封面图
	LikesCount    int        `json:"likes_count" gorm:"default:0;index"`                // 点赞数
	CommentsCount int        `json:"comments_count" gorm:"default:0"`                   // 评论数
	CollectsCount int        `json:"collects_count" gorm:"default:0;index"`            // 收藏数
	ViewsCount    int        `json:"views_count" gorm:"default:0"`                      // 浏览数
	CreatedAt     *time.Time `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SquareComment 广场评论
type SquareComment struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	PostID    int64      `json:"post_id" gorm:"index"`
	UserID    int64      `json:"user_id" gorm:"index"`
	Content   string     `json:"content" gorm:"type:text"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SquareLike 点赞记录
type SquareLike struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	PostID    int64      `json:"post_id" gorm:"index"`
	UserID    int64      `json:"user_id" gorm:"index"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	// 复合索引保证同一用户对同一帖子只能点赞一次
}

// SquareCollect 收藏记录
type SquareCollect struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	PostID    int64      `json:"post_id" gorm:"index"`
	UserID    int64      `json:"user_id" gorm:"index"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	// 复合索引保证同一用户对同一帖子只能收藏一次
}

// ===================== API REQUEST & RESPONSE =====================

// ListSquarePostsReq 浏览广场内容请求
type ListSquarePostsReq struct {
	Domain   string `json:"domain"`    // 可选，按领域筛选
	Keyword  string `json:"keyword"`   // 可选，关键词搜索
	Sort     string `json:"sort"`      // newest, hottest, trending (default: newest)
	Page     int    `json:"page"`      // 页码 (default: 1)
	PageSize int    `json:"page_size"` // 每页数量 (default: 20)
}

// ListSquarePostsResp 浏览广场内容响应
type ListSquarePostsResp struct {
	Data  []SquarePostDTO `json:"data"`
	Total int64           `json:"total"`
}

// SquarePostDTO 广场内容数据传输对象
type SquarePostDTO struct {
	ID            int64                  `json:"id"`
	DraftID       int64                  `json:"draft_id"`
	UserID        int64                  `json:"user_id"`
	User          *UserInfoDTO           `json:"user"`
	Title         string                 `json:"title"`
	PreviewText   string                 `json:"preview_text"`
	Domain        string                 `json:"domain"`
	Tags          []string               `json:"tags"`
	CoverImage    string                 `json:"cover_image"`
	LikesCount    int                    `json:"likes_count"`
	CommentsCount int                    `json:"comments_count"`
	CollectsCount int                    `json:"collects_count"`
	IsLiked       bool                   `json:"is_liked"`      // 当前用户是否点赞
	IsCollected   bool                   `json:"is_collected"`  // 当前用户是否收藏
	CreatedAt     *time.Time             `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at"`
}

// UserInfoDTO 用户信息数据传输对象
type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// GetSquarePostReq 获取广场内容详情
type GetSquarePostReq struct {
	ID int64 `json:"id"`
}

// PublishToSquareReq 发布草稿到广场
type PublishToSquareReq struct {
	DraftID int64 `json:"draft_id"`
}

// LikeReq 点赞请求
type LikeReq struct {
	PostID int64 `json:"post_id"`
}

// UnlikeReq 取消点赞请求
type UnlikeReq struct {
	PostID int64 `json:"post_id"`
}

// CollectReq 收藏请求
type CollectReq struct {
	PostID int64 `json:"post_id"`
}

// UncollectReq 取消收藏请求
type UncollectReq struct {
	PostID int64 `json:"post_id"`
}

// CommentReq 评论请求
type CommentReq struct {
	PostID  int64  `json:"post_id"`
	Content string `json:"content"`
}

// GetCommentsReq 获取评论列表请求
type GetCommentsReq struct {
	PostID   int64 `json:"post_id"`
	Page     int   `json:"page"`      // default: 1
	PageSize int   `json:"page_size"` // default: 20
}

// GetCommentsResp 评论列表响应
type GetCommentsResp struct {
	Data  []SquareCommentDTO `json:"data"`
	Total int64              `json:"total"`
}

// SquareCommentDTO 评论数据传输对象
type SquareCommentDTO struct {
	ID        int64      `json:"id"`
	PostID    int64      `json:"post_id"`
	UserID    int64      `json:"user_id"`
	User      *UserInfoDTO `json:"user"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ===================== Table Functions =====================

// TableName 指定 SquarePost 表名
func (s *SquarePost) TableName() string {
	return "square_posts"
}

// TableName 指定 SquareComment 表名
func (s *SquareComment) TableName() string {
	return "square_comments"
}

// TableName 指定 SquareLike 表名
func (s *SquareLike) TableName() string {
	return "square_likes"
}

// TableName 指定 SquareCollect 表名
func (s *SquareCollect) TableName() string {
	return "square_collects"
}
