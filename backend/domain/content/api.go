package content

import (
	"backend/datasource/dbdao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContentDomain 内容domain主类
type ContentDomain struct {
	draftService   *DraftService
	contentService *ContentService
}

// NewContentDomain 创建Content Domain
func NewContentDomain(draftService *DraftService, contentService *ContentService) *ContentDomain {
	return &ContentDomain{
		draftService:   draftService,
		contentService: contentService,
	}
}

// NewContentDomainWithDB 从数据库连接创建Content Domain (便利函数用于main.go)
func NewContentDomainWithDB(db *dbdao.DB) *ContentDomain {
	gormDB := db.DB()
	return newContentDomainWithGormDB(gormDB)
}

// newContentDomainWithGormDB 内部函数，从GORM DB创建Content Domain
func newContentDomainWithGormDB(db *gorm.DB) *ContentDomain {
	// 创建repositories
	draftRepo := NewDraftRepository(db)
	contentRepo := NewContentRepository(db)
	versionRepo := NewContentVersionRepository(db)

	// 创建services
	draftService := NewDraftService(draftRepo, versionRepo, db)
	contentService := NewContentService(contentRepo, versionRepo, db)

	// 返回ContentDomain
	return NewContentDomain(draftService, contentService)
}

// ===================== Draft API Handlers =====================

// ApiCreateDraft 创建草稿
func (d *ContentDomain) ApiCreateDraft(c *gin.Context) {
	var req CreateDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	draft, err := d.draftService.CreateDraft(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": draft,
	})
}

// ApiGetDraft 获取草稿
func (d *ContentDomain) ApiGetDraft(c *gin.Context) {
	var req GetDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	draft, err := d.draftService.GetDraft(c.Request.Context(), userID, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": draft,
	})
}

// ApiSearchDrafts 搜索草稿
func (d *ContentDomain) ApiSearchDrafts(c *gin.Context) {
	var req SearchDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	result, err := d.draftService.SearchDrafts(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// ApiUpdateDraft 更新草稿
func (d *ContentDomain) ApiUpdateDraft(c *gin.Context) {
	var req UpdateDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	result, err := d.draftService.UpdateDraft(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// ApiDeleteDraft 删除草稿
func (d *ContentDomain) ApiDeleteDraft(c *gin.Context) {
	var req DeleteDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.draftService.DeleteDraft(c.Request.Context(), userID, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "draft deleted successfully",
	})
}

// ApiGetDraftVersions 获取草稿版本历史
func (d *ContentDomain) ApiGetDraftVersions(c *gin.Context) {
	var req GetDraftVersionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取
	limit := 20
	if req.Limit > 0 {
		limit = int(req.Limit)
	}

	result, err := d.draftService.GetDraftVersions(c.Request.Context(), userID, req.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// ApiRevertDraftVersion 恢复草稿版本
func (d *ContentDomain) ApiRevertDraftVersion(c *gin.Context) {
	draftIDStr := c.Param("id")
	versionStr := c.Param("version")

	draftID, err := strconv.ParseInt(draftIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	result, err := d.draftService.RevertDraftVersion(c.Request.Context(), userID, draftID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// ApiPublishDraft 发布草稿
func (d *ContentDomain) ApiPublishDraft(c *gin.Context) {
	draftIDStr := c.Param("id")

	draftID, err := strconv.ParseInt(draftIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid draft id"})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	content, err := d.draftService.PublishDraft(c.Request.Context(), userID, draftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": content,
	})
}

// ===================== Content API Handlers =====================

// ApiCreateContent 创建内容
func (d *ContentDomain) ApiCreateContent(c *gin.Context) {
	var req AddContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	content, err := d.contentService.CreateContent(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": content,
	})
}

// ApiGetContent 获取内容
func (d *ContentDomain) ApiGetContent(c *gin.Context) {
	var req GetContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	content, err := d.contentService.GetContent(c.Request.Context(), userID, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": content,
	})
}

// ApiSearchContents 搜索内容
func (d *ContentDomain) ApiSearchContents(c *gin.Context) {
	var req SearchContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	result, err := d.contentService.SearchContents(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// ApiUpdateContent 更新内容
func (d *ContentDomain) ApiUpdateContent(c *gin.Context) {
	var req UpdateContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	content, err := d.contentService.UpdateContent(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": content,
	})
}

// ApiDeleteContent 删除内容
func (d *ContentDomain) ApiDeleteContent(c *gin.Context) {
	var req DeleteContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证信息中获取

	if err := d.contentService.DeleteContent(c.Request.Context(), userID, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "content deleted successfully",
	})
}
