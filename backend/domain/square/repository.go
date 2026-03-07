package square

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// ===================== REPOSITORY LAYER =====================

// SquarePostRepository 广场帖子存储库
type SquarePostRepository struct {
	db *gorm.DB
}

// NewSquarePostRepository 创建帖子存储库
func NewSquarePostRepository(db *gorm.DB) *SquarePostRepository {
	return &SquarePostRepository{db: db}
}

// Create 创建帖子
func (r *SquarePostRepository) Create(ctx context.Context, post *SquarePost) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// Get 获取帖子详情
func (r *SquarePostRepository) Get(ctx context.Context, id int64) (*SquarePost, error) {
	var post *SquarePost
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// List 列表查询
func (r *SquarePostRepository) List(ctx context.Context, query *ListQuery) ([]*SquarePost, int64, error) {
	var posts []*SquarePost
	var total int64

	db := r.db.WithContext(ctx)

	// 构建查询条件
	if query.Domain != "" {
		db = db.Where("domain = ?", query.Domain)
	}

	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR preview_text LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 计算总数
	if err := db.Model(&SquarePost{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch query.Sort {
	case "hottest":
		db = db.Order("likes_count DESC, created_at DESC")
	case "trending":
		db = db.Order("(likes_count + comments_count + collects_count) DESC, created_at DESC")
	default: // newest
		db = db.Order("created_at DESC")
	}

	// 分页
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// SquareCommentRepository 评论存储库
type SquareCommentRepository struct {
	db *gorm.DB
}

// NewSquareCommentRepository 创建评论存储库
func NewSquareCommentRepository(db *gorm.DB) *SquareCommentRepository {
	return &SquareCommentRepository{db: db}
}

// Create 创建评论
func (r *SquareCommentRepository) Create(ctx context.Context, comment *SquareComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// ListByPostID 按帖子 ID 列表查询
func (r *SquareCommentRepository) ListByPostID(ctx context.Context, postID int64, page, pageSize int) ([]*SquareComment, int64, error) {
	var comments []*SquareComment
	var total int64

	db := r.db.WithContext(ctx).Where("post_id = ?", postID)

	// 计算总数
	if err := db.Model(&SquareComment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// SquareLikeRepository 点赞存储库
type SquareLikeRepository struct {
	db *gorm.DB
}

// NewSquareLikeRepository 创建点赞存储库
func NewSquareLikeRepository(db *gorm.DB) *SquareLikeRepository {
	return &SquareLikeRepository{db: db}
}

// Create 创建点赞记录
func (r *SquareLikeRepository) Create(ctx context.Context, postID, userID int64) error {
	// 检查是否已点赞
	var count int64
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Model(&SquareLike{}).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("already liked")
	}

	like := &SquareLike{PostID: postID, UserID: userID}
	return r.db.WithContext(ctx).Create(like).Error
}

// Delete 删除点赞记录
func (r *SquareLikeRepository) Delete(ctx context.Context, postID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&SquareLike{}).Error
}

// Exists 检查是否已点赞
func (r *SquareLikeRepository) Exists(ctx context.Context, postID, userID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Model(&SquareLike{}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// SquareCollectRepository 收藏存储库
type SquareCollectRepository struct {
	db *gorm.DB
}

// NewSquareCollectRepository 创建收藏存储库
func NewSquareCollectRepository(db *gorm.DB) *SquareCollectRepository {
	return &SquareCollectRepository{db: db}
}

// Create 创建收藏记录
func (r *SquareCollectRepository) Create(ctx context.Context, postID, userID int64) error {
	// 检查是否已收藏
	var count int64
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Model(&SquareCollect{}).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("already collected")
	}

	collect := &SquareCollect{PostID: postID, UserID: userID}
	return r.db.WithContext(ctx).Create(collect).Error
}

// Delete 删除收藏记录
func (r *SquareCollectRepository) Delete(ctx context.Context, postID, userID int64) error {
	return r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&SquareCollect{}).Error
}

// Exists 检查是否已收藏
func (r *SquareCollectRepository) Exists(ctx context.Context, postID, userID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Model(&SquareCollect{}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ===================== HELPER TYPES =====================

// ListQuery 查询参数
type ListQuery struct {
	Domain   string
	Keyword  string
	Sort     string
	Page     int
	PageSize int
}

// ParseTagsFromJSON 从 JSON 字符串解析标签数组
func ParseTagsFromJSON(tagsJSON string) ([]string, error) {
	var tags []string
	if tagsJSON == "" {
		return tags, nil
	}
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// TagsToJSON 将标签数组转换为 JSON 字符串
func TagsToJSON(tags []string) (string, error) {
	if len(tags) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(tags)
	return string(data), err
}
