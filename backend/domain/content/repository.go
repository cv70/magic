package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DraftRepository Draft数据访问接口
type DraftRepository interface {
	Create(ctx context.Context, draft *Draft) error
	GetByID(ctx context.Context, id int64) (*Draft, error)
	GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]*Draft, int64, error)
	Search(ctx context.Context, userID int64, keyword, status string, tags []string, limit int, offset int) ([]*Draft, int64, error)
	Update(ctx context.Context, draft *Draft) error
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64, id int64) error
}

// draftRepositoryImpl Draft Repository实现
type draftRepositoryImpl struct {
	db *gorm.DB
}

// NewDraftRepository 创建Draft Repository
func NewDraftRepository(db *gorm.DB) DraftRepository {
	return &draftRepositoryImpl{db: db}
}

// Create 创建新草稿
func (r *draftRepositoryImpl) Create(ctx context.Context, draft *Draft) error {
	if draft == nil {
		return errors.New("draft cannot be nil")
	}

	if draft.Title == "" {
		return errors.New("draft title cannot be empty")
	}

	draft.CreatedAt = ptrTime(time.Now())
	draft.UpdatedAt = ptrTime(time.Now())
	draft.LastEditedAt = ptrTime(time.Now())
	draft.Status = "draft"

	if err := r.db.WithContext(ctx).Create(draft).Error; err != nil {
		return fmt.Errorf("failed to create draft: %w", err)
	}

	return nil
}

// GetByID 根据ID获取草稿
func (r *draftRepositoryImpl) GetByID(ctx context.Context, id int64) (*Draft, error) {
	var draft Draft

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&draft).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("draft not found")
		}
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	return &draft, nil
}

// GetByUserID 获取用户的草稿列表
func (r *draftRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]*Draft, int64, error) {
	var drafts []*Draft
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "draft")

	// 获取总数
	if err := query.Model(&Draft{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count drafts: %w", err)
	}

	// 获取数据
	if err := query.
		Order("last_edited_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&drafts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get drafts: %w", err)
	}

	return drafts, total, nil
}

// Search 搜索草稿
func (r *draftRepositoryImpl) Search(ctx context.Context, userID int64, keyword, status string, tags []string, limit int, offset int) ([]*Draft, int64, error) {
	var drafts []*Draft
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 标签筛选
	if len(tags) > 0 {
		// 这是一个简化的实现，实际上JSON搜索在不同数据库中有不同方式
		for _, tag := range tags {
			query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, tag))
		}
	}

	// 获取总数
	if err := query.Model(&Draft{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count drafts: %w", err)
	}

	// 获取数据
	if err := query.
		Order("last_edited_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&drafts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search drafts: %w", err)
	}

	return drafts, total, nil
}

// Update 更新草稿
func (r *draftRepositoryImpl) Update(ctx context.Context, draft *Draft) error {
	if draft == nil || draft.ID == 0 {
		return errors.New("draft cannot be nil or id must be set")
	}

	draft.UpdatedAt = ptrTime(time.Now())
	draft.LastEditedAt = ptrTime(time.Now())

	result := r.db.WithContext(ctx).
		Model(&Draft{}).
		Where("id = ?", draft.ID).
		Updates(draft)

	if result.Error != nil {
		return fmt.Errorf("failed to update draft: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("draft not found")
	}

	return nil
}

// Delete 删除草稿
func (r *draftRepositoryImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid draft id")
	}

	if err := r.db.WithContext(ctx).Delete(&Draft{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	return nil
}

// DeleteByUserID 删除用户的草稿（确保权限）
func (r *draftRepositoryImpl) DeleteByUserID(ctx context.Context, userID int64, id int64) error {
	if id <= 0 {
		return errors.New("invalid draft id")
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Draft{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete draft: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("draft not found or permission denied")
	}

	return nil
}

// ===== ContentRepository 实现 =====

// ContentRepository Content数据访问接口
type ContentRepository interface {
	Create(ctx context.Context, content *Content) error
	GetByID(ctx context.Context, id int64) (*Content, error)
	GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]*Content, int64, error)
	Search(ctx context.Context, userID int64, keyword, contentType, status string, limit int, offset int) ([]*Content, int64, error)
	Update(ctx context.Context, content *Content) error
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64, id int64) error
}

type contentRepositoryImpl struct {
	db *gorm.DB
}

// NewContentRepository 创建Content Repository
func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepositoryImpl{db: db}
}

func (r *contentRepositoryImpl) Create(ctx context.Context, content *Content) error {
	if content == nil || content.Title == "" {
		return errors.New("content cannot be nil and title must not be empty")
	}

	content.CreatedAt = ptrTime(time.Now())
	content.UpdatedAt = ptrTime(time.Now())
	content.PublishedAt = ptrTime(time.Now())

	if err := r.db.WithContext(ctx).Create(content).Error; err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	return nil
}

func (r *contentRepositoryImpl) GetByID(ctx context.Context, id int64) (*Content, error) {
	var content Content

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&content).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("content not found")
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	return &content, nil
}

func (r *contentRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit int, offset int) ([]*Content, int64, error) {
	var contents []*Content
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if err := query.Model(&Content{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contents: %w", err)
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&contents).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get contents: %w", err)
	}

	return contents, total, nil
}

func (r *contentRepositoryImpl) Search(ctx context.Context, userID int64, keyword, contentType, status string, limit int, offset int) ([]*Content, int64, error) {
	var contents []*Content
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if keyword != "" {
		query = query.Where("title LIKE ? OR body LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&Content{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contents: %w", err)
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&contents).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search contents: %w", err)
	}

	return contents, total, nil
}

func (r *contentRepositoryImpl) Update(ctx context.Context, content *Content) error {
	if content == nil || content.ID == 0 {
		return errors.New("content cannot be nil or id must be set")
	}

	content.UpdatedAt = ptrTime(time.Now())

	result := r.db.WithContext(ctx).
		Model(&Content{}).
		Where("id = ?", content.ID).
		Updates(content)

	if result.Error != nil {
		return fmt.Errorf("failed to update content: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("content not found")
	}

	return nil
}

func (r *contentRepositoryImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid content id")
	}

	if err := r.db.WithContext(ctx).Delete(&Content{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}

func (r *contentRepositoryImpl) DeleteByUserID(ctx context.Context, userID int64, id int64) error {
	if id <= 0 {
		return errors.New("invalid content id")
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Content{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete content: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("content not found or permission denied")
	}

	return nil
}

// ===== ContentVersionRepository 实现 =====

type ContentVersionRepository interface {
	Create(ctx context.Context, version *ContentVersion) error
	GetByContentID(ctx context.Context, contentID int64, limit int) ([]*ContentVersion, error)
	GetByVersion(ctx context.Context, contentID int64, versionNum int) (*ContentVersion, error)
}

type contentVersionRepositoryImpl struct {
	db *gorm.DB
}

func NewContentVersionRepository(db *gorm.DB) ContentVersionRepository {
	return &contentVersionRepositoryImpl{db: db}
}

func (r *contentVersionRepositoryImpl) Create(ctx context.Context, version *ContentVersion) error {
	if version == nil {
		return errors.New("version cannot be nil")
	}

	version.CreatedAt = ptrTime(time.Now())

	if err := r.db.WithContext(ctx).Create(version).Error; err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}

	return nil
}

func (r *contentVersionRepositoryImpl) GetByContentID(ctx context.Context, contentID int64, limit int) ([]*ContentVersion, error) {
	var versions []*ContentVersion

	if err := r.db.WithContext(ctx).
		Where("content_id = ?", contentID).
		Order("version_num DESC").
		Limit(limit).
		Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	return versions, nil
}

func (r *contentVersionRepositoryImpl) GetByVersion(ctx context.Context, contentID int64, versionNum int) (*ContentVersion, error) {
	var version ContentVersion

	if err := r.db.WithContext(ctx).
		Where("content_id = ? AND version_num = ?", contentID, versionNum).
		First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &version, nil
}

// ===== 工具函数 =====

func ptrTime(t time.Time) *time.Time {
	return &t
}

// parseJSONTags 将JSON字符串解析为字符串数组
func parseJSONTags(tagsJSON string) []string {
	if tagsJSON == "" {
		return []string{}
	}

	var tags []string
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		return []string{}
	}

	return tags
}

// tagsToJSON 将字符串数组转换为JSON
func tagsToJSON(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}

	data, _ := json.Marshal(tags)
	return string(data)
}
