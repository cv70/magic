package square

import (
	"backend/datasource/dbdao"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===================== SERVICE LAYER =====================

// SquareService 广场服务
type SquareService struct {
	postRepo    *SquarePostRepository
	commentRepo *SquareCommentRepository
	likeRepo    *SquareLikeRepository
	collectRepo *SquareCollectRepository
	db          *gorm.DB
}

// NewSquareService 创建广场服务
func NewSquareService(
	postRepo *SquarePostRepository,
	commentRepo *SquareCommentRepository,
	likeRepo *SquareLikeRepository,
	collectRepo *SquareCollectRepository,
	db *gorm.DB,
) *SquareService {
	return &SquareService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		collectRepo: collectRepo,
		db:          db,
	}
}

// ListPosts 浏览广场内容（带搜索、过滤、排序）
func (s *SquareService) ListPosts(ctx context.Context, req *ListSquarePostsReq, currentUserID int64) (*ListSquarePostsResp, error) {
	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "newest"
	}

	query := &ListQuery{
		Domain:   req.Domain,
		Keyword:  req.Keyword,
		Sort:     req.Sort,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	posts, total, err := s.postRepo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO，并添加用户点赞/收藏状态
	dtos := make([]SquarePostDTO, len(posts))
	for i, post := range posts {
		dto, err := s.convertToDTO(ctx, post, currentUserID)
		if err != nil {
			return nil, err
		}
		dtos[i] = *dto
	}

	return &ListSquarePostsResp{
		Data:  dtos,
		Total: total,
	}, nil
}

// GetPost 获取广场内容详情
func (s *SquareService) GetPost(ctx context.Context, id int64, currentUserID int64) (*SquarePostDTO, error) {
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// 增加浏览次数
	if err := s.db.WithContext(ctx).Model(post).Update("views_count", gorm.Expr("views_count + ?", 1)).Error; err != nil {
		// 不返回错误，继续
	}

	return s.convertToDTO(ctx, post, currentUserID)
}

// PublishToSquare 发布草稿到广场
// 注意：这里假设草稿数据已存在，实际需要从 content domain 获取
func (s *SquareService) PublishToSquare(ctx context.Context, req *PublishToSquareReq, userID int64) (*SquarePostDTO, error) {
	if req.DraftID == 0 {
		return nil, errors.New("draft_id is required")
	}

	// TODO: 从 content domain 获取草稿数据
	// draftService := ... // 需要依赖注入
	// draft := draftService.GetDraft(req.DraftID)

	// 这里暂时使用占位符逻辑
	post := &SquarePost{
		DraftID:     req.DraftID,
		UserID:      userID,
		Title:       "Sample Post", // 从草稿获取
		PreviewText: "Sample preview", // 从草稿获取
		Domain:      "film-drama", // 从草稿 metadata 获取
		Tags:        "[]",
		CoverImage:  "",
		LikesCount:  0,
		CommentsCount: 0,
		CollectsCount: 0,
		ViewsCount:  0,
		CreatedAt:   nil,
		UpdatedAt:   nil,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return s.convertToDTO(ctx, post, userID)
}

// Like 点赞
func (s *SquareService) Like(ctx context.Context, postID int64, userID int64) error {
	// 创建点赞记录
	if err := s.likeRepo.Create(ctx, postID, userID); err != nil {
		return err
	}

	// 更新帖子点赞计数
	if err := s.db.WithContext(ctx).
		Model(&SquarePost{}).
		Where("id = ?", postID).
		Update("likes_count", gorm.Expr("likes_count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// Unlike 取消点赞
func (s *SquareService) Unlike(ctx context.Context, postID int64, userID int64) error {
	// 删除点赞记录
	if err := s.likeRepo.Delete(ctx, postID, userID); err != nil {
		return err
	}

	// 更新帖子点赞计数
	if err := s.db.WithContext(ctx).
		Model(&SquarePost{}).
		Where("id = ?", postID).
		Update("likes_count", gorm.Expr("likes_count - ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// Collect 收藏
func (s *SquareService) Collect(ctx context.Context, postID int64, userID int64) error {
	// 创建收藏记录
	if err := s.collectRepo.Create(ctx, postID, userID); err != nil {
		return err
	}

	// 更新帖子收藏计数
	if err := s.db.WithContext(ctx).
		Model(&SquarePost{}).
		Where("id = ?", postID).
		Update("collects_count", gorm.Expr("collects_count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// Uncollect 取消收藏
func (s *SquareService) Uncollect(ctx context.Context, postID int64, userID int64) error {
	// 删除收藏记录
	if err := s.collectRepo.Delete(ctx, postID, userID); err != nil {
		return err
	}

	// 更新帖子收藏计数
	if err := s.db.WithContext(ctx).
		Model(&SquarePost{}).
		Where("id = ?", postID).
		Update("collects_count", gorm.Expr("collects_count - ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// AddComment 添加评论
func (s *SquareService) AddComment(ctx context.Context, req *CommentReq, userID int64) (*SquareCommentDTO, error) {
	if req.PostID == 0 {
		return nil, errors.New("post_id is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	// 验证帖子存在
	_, err := s.postRepo.Get(ctx, req.PostID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comment := &SquareComment{
		PostID:  req.PostID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// 更新帖子评论计数
	if err := s.db.WithContext(ctx).
		Model(&SquarePost{}).
		Where("id = ?", req.PostID).
		Update("comments_count", gorm.Expr("comments_count + ?", 1)).Error; err != nil {
		return nil, err
	}

	dto := &SquareCommentDTO{
		ID:      comment.ID,
		PostID:  comment.PostID,
		UserID:  comment.UserID,
		Content: comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	return dto, nil
}

// GetComments 获取评论列表
func (s *SquareService) GetComments(ctx context.Context, req *GetCommentsReq) (*GetCommentsResp, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	comments, total, err := s.commentRepo.ListByPostID(ctx, req.PostID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	dtos := make([]SquareCommentDTO, len(comments))
	for i, comment := range comments {
		dtos[i] = SquareCommentDTO{
			ID:        comment.ID,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	return &GetCommentsResp{
		Data:  dtos,
		Total: total,
	}, nil
}

// ===================== HELPER METHODS =====================

// convertToDTO 将 Post 转换为 DTO，包括点赞/收藏状态
func (s *SquareService) convertToDTO(ctx context.Context, post *SquarePost, currentUserID int64) (*SquarePostDTO, error) {
	tags, err := ParseTagsFromJSON(post.Tags)
	if err != nil {
		tags = []string{}
	}

	// 检查当前用户是否点赞和收藏
	isLiked, _ := s.likeRepo.Exists(ctx, post.ID, currentUserID)
	isCollected, _ := s.collectRepo.Exists(ctx, post.ID, currentUserID)

	// TODO: 从 identity domain 获取用户信息
	userInfo := &UserInfoDTO{
		ID:       post.UserID,
		Username: "User" + strconv.FormatInt(post.UserID, 10),
		Avatar:   "",
	}

	return &SquarePostDTO{
		ID:            post.ID,
		DraftID:       post.DraftID,
		UserID:        post.UserID,
		User:          userInfo,
		Title:         post.Title,
		PreviewText:   post.PreviewText,
		Domain:        post.Domain,
		Tags:          tags,
		CoverImage:    post.CoverImage,
		LikesCount:    post.LikesCount,
		CommentsCount: post.CommentsCount,
		CollectsCount: post.CollectsCount,
		IsLiked:       isLiked,
		IsCollected:   isCollected,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
	}, nil
}

// ===================== DOMAIN MAIN CLASS =====================

// SquareDomain Square Domain 主类
type SquareDomain struct {
	service *SquareService
}

// NewSquareDomain 创建 Square Domain
func NewSquareDomain(service *SquareService) *SquareDomain {
	return &SquareDomain{service: service}
}

// NewSquareDomainWithDB 从数据库连接创建 Square Domain
func NewSquareDomainWithDB(db *dbdao.DB) *SquareDomain {
	gormDB := db.DB()
	return newSquareDomainWithGormDB(gormDB)
}

// newSquareDomainWithGormDB 内部函数，从 GORM DB 创建 Square Domain
func newSquareDomainWithGormDB(db *gorm.DB) *SquareDomain {
	// 创建 repositories
	postRepo := NewSquarePostRepository(db)
	commentRepo := NewSquareCommentRepository(db)
	likeRepo := NewSquareLikeRepository(db)
	collectRepo := NewSquareCollectRepository(db)

	// 创建 service
	service := NewSquareService(postRepo, commentRepo, likeRepo, collectRepo, db)

	// 返回 SquareDomain
	return NewSquareDomain(service)
}

// ===================== API HANDLERS =====================

// ApiListSquarePosts 浏览广场内容
func (d *SquareDomain) ApiListSquarePosts(c *gin.Context) {
	var req ListSquarePostsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	resp, err := d.service.ListPosts(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resp.Data,
		"total": resp.Total,
	})
}

// ApiGetSquarePost 获取广场内容详情
func (d *SquareDomain) ApiGetSquarePost(c *gin.Context) {
	var req GetSquarePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	post, err := d.service.GetPost(c.Request.Context(), req.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": post,
	})
}

// ApiPublishToSquare 发布草稿到广场
func (d *SquareDomain) ApiPublishToSquare(c *gin.Context) {
	var req PublishToSquareReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	post, err := d.service.PublishToSquare(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": post,
	})
}

// ApiLike 点赞
func (d *SquareDomain) ApiLike(c *gin.Context) {
	var req LikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.service.Like(c.Request.Context(), req.PostID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "liked",
	})
}

// ApiUnlike 取消点赞
func (d *SquareDomain) ApiUnlike(c *gin.Context) {
	var req UnlikeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.service.Unlike(c.Request.Context(), req.PostID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "unliked",
	})
}

// ApiCollect 收藏
func (d *SquareDomain) ApiCollect(c *gin.Context) {
	var req CollectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.service.Collect(c.Request.Context(), req.PostID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "collected",
	})
}

// ApiUncollect 取消收藏
func (d *SquareDomain) ApiUncollect(c *gin.Context) {
	var req UncollectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.service.Uncollect(c.Request.Context(), req.PostID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "uncollected",
	})
}

// ApiAddComment 添加评论
func (d *SquareDomain) ApiAddComment(c *gin.Context) {
	var req CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	comment, err := d.service.AddComment(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": comment,
	})
}

// ApiGetComments 获取评论列表
func (d *SquareDomain) ApiGetComments(c *gin.Context) {
	var req GetCommentsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	resp, err := d.service.GetComments(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resp.Data,
		"total": resp.Total,
	})
}
