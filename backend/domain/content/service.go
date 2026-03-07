package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DraftService Draft业务逻辑
type DraftService struct {
	draftRepo    DraftRepository
	versionRepo  ContentVersionRepository
	db           *gorm.DB
}

// NewDraftService 创建Draft Service
func NewDraftService(draftRepo DraftRepository, versionRepo ContentVersionRepository, db *gorm.DB) *DraftService {
	return &DraftService{
		draftRepo:   draftRepo,
		versionRepo: versionRepo,
		db:          db,
	}
}

// CreateDraft 创建新草稿
func (s *DraftService) CreateDraft(ctx context.Context, userID int64, req *CreateDraftReq) (*Draft, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	draft := &Draft{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		ContentType: req.ContentType,
		Status:      "draft",
		Tags:        tagsToJSON(req.Tags),
	}

	if err := s.draftRepo.Create(ctx, draft); err != nil {
		return nil, err
	}

	return draft, nil
}

// GetDraft 获取草稿详情
func (s *DraftService) GetDraft(ctx context.Context, userID, draftID int64) (*GetDraftRes, error) {
	draft, err := s.draftRepo.GetByID(ctx, draftID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if draft.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return draftToResponse(draft), nil
}

// SearchDrafts 搜索草稿
func (s *DraftService) SearchDrafts(ctx context.Context, userID int64, req *SearchDraftReq) (*SearchDraftRes, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	var tags []string
	if req.Tags != "" {
		// 解析逗号分割的标签
		// tags = strings.Split(req.Tags, ",")
	}

	// 如果没有status，默认搜索draft
	if req.Status == "" {
		req.Status = "draft"
	}

	drafts, total, err := s.draftRepo.Search(ctx, userID, req.Keyword, req.Status, tags, int(req.Limit), int(offset))
	if err != nil {
		return nil, err
	}

	res := &SearchDraftRes{
		Data:  make([]*GetDraftRes, len(drafts)),
		Total: total,
	}

	for i, draft := range drafts {
		res.Data[i] = draftToResponse(draft)
	}

	return res, nil
}

// UpdateDraft 更新草稿
// 每次更新都会自动创建一个新版本
func (s *DraftService) UpdateDraft(ctx context.Context, userID int64, req *UpdateDraftReq) (*UpdateDraftRes, error) {
	if req.ID <= 0 {
		return nil, errors.New("invalid draft id")
	}

	// 获取原草稿（权限检查）
	draft, err := s.draftRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if draft.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 更新字段
	if req.Title != nil {
		draft.Title = *req.Title
	}
	if req.Content != nil {
		draft.Content = *req.Content
	}
	if req.ContentType != nil {
		draft.ContentType = *req.ContentType
	}
	if len(req.Tags) > 0 {
		draft.Tags = tagsToJSON(req.Tags)
	}

	// 在事务中保存草稿和版本
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新草稿
		if err := s.draftRepo.Update(ctx, draft); err != nil {
			return err
		}

		// 创建版本记录
		version := &ContentVersion{
			ContentID:     draft.ID, // 暂时使用draft.ID，实际场景中应该是content_id
			DraftID:       &draft.ID,
			Title:         draft.Title,
			Content:       draft.Content,
			ChangeSummary: req.ChangeSummary,
			CreatedBy:     userID,
		}

		if err := s.versionRepo.Create(ctx, version); err != nil {
			return err
		}

		// 更新SavedVersionID
		draft.SavedVersionID = &version.ID
		if err := s.draftRepo.Update(ctx, draft); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return &UpdateDraftRes{
		ID:        draft.ID,
		UpdatedAt: draft.UpdatedAt,
		VersionNum: 1, // TODO: 计算实际版本号
	}, nil
}

// DeleteDraft 删除草稿
func (s *DraftService) DeleteDraft(ctx context.Context, userID, draftID int64) error {
	// 权限检查：确保只能删除自己的草稿
	return s.draftRepo.DeleteByUserID(ctx, userID, draftID)
}

// GetDraftVersions 获取草稿的版本历史
func (s *DraftService) GetDraftVersions(ctx context.Context, userID, draftID int64, limit int) (*GetDraftVersionsRes, error) {
	// 先检查草稿是否存在且属于该用户
	draft, err := s.draftRepo.GetByID(ctx, draftID)
	if err != nil {
		return nil, err
	}

	if draft.UserID != userID {
		return nil, errors.New("permission denied")
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	versions, err := s.versionRepo.GetByContentID(ctx, draftID, limit)
	if err != nil {
		return nil, err
	}

	res := &GetDraftVersionsRes{
		Data: make([]*VersionInfo, len(versions)),
	}

	for i, v := range versions {
		res.Data[i] = &VersionInfo{
			VersionNum:    v.VersionNum,
			Title:         v.Title,
			ChangeSummary: v.ChangeSummary,
			CreatedAt:     v.CreatedAt,
			// TODO: 关联查询用户信息获取创建人名称
			CreatedByUser: fmt.Sprintf("user_%d", v.CreatedBy),
		}
	}

	return res, nil
}

// RevertDraftVersion 恢复草稿到指定版本
func (s *DraftService) RevertDraftVersion(ctx context.Context, userID, draftID int64, versionNum int) (*RevertDraftVersionRes, error) {
	// 权限检查
	draft, err := s.draftRepo.GetByID(ctx, draftID)
	if err != nil {
		return nil, err
	}

	if draft.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 获取指定版本
	version, err := s.versionRepo.GetByVersion(ctx, draftID, versionNum)
	if err != nil {
		return nil, err
	}

	// 恢复内容
	draft.Title = version.Title
	draft.Content = version.Content

	// 在事务中保存
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新草稿
		if err := s.draftRepo.Update(ctx, draft); err != nil {
			return err
		}

		// 创建新版本记录（恢复操作本身也会创建版本）
		newVersion := &ContentVersion{
			ContentID:     draftID,
			DraftID:       &draftID,
			Title:         draft.Title,
			Content:       draft.Content,
			ChangeSummary: fmt.Sprintf("恢复到版本%d", versionNum),
			CreatedBy:     userID,
		}

		if err := s.versionRepo.Create(ctx, newVersion); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &RevertDraftVersionRes{
		ID:         draft.ID,
		VersionNum: 1, // TODO: 计算实际的新版本号
		RevertedFrom: versionNum,
		UpdatedAt:  draft.UpdatedAt,
	}, nil
}

// PublishDraft 将草稿发布为正式内容
func (s *DraftService) PublishDraft(ctx context.Context, userID, draftID int64) (*Content, error) {
	// 获取草稿
	draft, err := s.draftRepo.GetByID(ctx, draftID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if draft.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 创建Content
	var tagsArray []string
	if err := json.Unmarshal([]byte(draft.Tags), &tagsArray); err != nil {
		tagsArray = []string{}
	}

	content := &Content{
		UserID:      userID,
		Title:       draft.Title,
		Body:        draft.Content,
		ContentType: draft.ContentType,
		Status:      "published",
		Tags:        draft.Tags,
		Metadata:    draft.Metadata,
		PublishedAt: ptrTime(time.Now()),
	}

	// 在事务中保存
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// TODO: 使用ContentRepository创建content
		if err := tx.Create(content).Error; err != nil {
			return err
		}

		// 更新草稿状态
		draft.Status = "archived"
		if err := s.draftRepo.Update(ctx, draft); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return content, nil
}

// ===== ContentService =====

// ContentService Content业务逻辑
type ContentService struct {
	contentRepo ContentRepository
	versionRepo ContentVersionRepository
	db          *gorm.DB
}

// NewContentService 创建Content Service
func NewContentService(contentRepo ContentRepository, versionRepo ContentVersionRepository, db *gorm.DB) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
		versionRepo: versionRepo,
		db:          db,
	}
}

// CreateContent 创建内容
func (s *ContentService) CreateContent(ctx context.Context, userID int64, req *AddContentReq) (*Content, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if req.Title == "" || req.Body == "" {
		return nil, errors.New("title and body cannot be empty")
	}

	content := &Content{
		UserID:      userID,
		Title:       req.Title,
		Body:        req.Body,
		ContentType: req.ContentType,
		Status:      req.Status,
		Tags:        tagsToJSON(req.Tags),
	}

	if err := s.contentRepo.Create(ctx, content); err != nil {
		return nil, err
	}

	return content, nil
}

// GetContent 获取内容
func (s *ContentService) GetContent(ctx context.Context, userID, contentID int64) (*Content, error) {
	content, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if content.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return content, nil
}

// SearchContents 搜索内容
func (s *ContentService) SearchContents(ctx context.Context, userID int64, req *SearchContentReq) (*SearchContentRes, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	contents, total, err := s.contentRepo.Search(ctx, userID, req.Query, req.ContentType, req.Status, int(req.Limit), int(offset))
	if err != nil {
		return nil, err
	}

	res := &SearchContentRes{
		Contents: make([]*GetContentRes, len(contents)),
		Total:    total,
	}

	for i, content := range contents {
		res.Contents[i] = contentToResponse(content)
	}

	return res, nil
}

// UpdateContent 更新内容
func (s *ContentService) UpdateContent(ctx context.Context, userID int64, req *UpdateContentReq) (*Content, error) {
	if req.ID <= 0 {
		return nil, errors.New("invalid content id")
	}

	content, err := s.contentRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if content.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 更新字段
	if req.Title != nil {
		content.Title = *req.Title
	}
	if req.Body != nil {
		content.Body = *req.Body
	}
	if req.ContentType != nil {
		content.ContentType = *req.ContentType
	}
	if req.Status != nil {
		content.Status = *req.Status
	}
	if len(req.Tags) > 0 {
		content.Tags = tagsToJSON(req.Tags)
	}

	if err := s.contentRepo.Update(ctx, content); err != nil {
		return nil, err
	}

	return content, nil
}

// DeleteContent 删除内容
func (s *ContentService) DeleteContent(ctx context.Context, userID, contentID int64) error {
	return s.contentRepo.DeleteByUserID(ctx, userID, contentID)
}

// ===== 工具函数 =====

func draftToResponse(draft *Draft) *GetDraftRes {
	var tags []string
	if err := json.Unmarshal([]byte(draft.Tags), &tags); err != nil {
		tags = []string{}
	}

	return &GetDraftRes{
		ID:           draft.ID,
		Title:        draft.Title,
		Content:      draft.Content,
		ContentType:  draft.ContentType,
		Status:       draft.Status,
		Tags:         tags,
		LastEditedAt: draft.LastEditedAt,
		CreatedAt:    draft.CreatedAt,
		UpdatedAt:    draft.UpdatedAt,
	}
}

func contentToResponse(content *Content) *GetContentRes {
	var tags []string
	if err := json.Unmarshal([]byte(content.Tags), &tags); err != nil {
		tags = []string{}
	}

	return &GetContentRes{
		ID:          content.ID,
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        tags,
		CreatedAt:   content.CreatedAt,
		PublishedAt: content.PublishedAt,
	}
}
