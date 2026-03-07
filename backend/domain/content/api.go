package content

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiGetContent 获取内容详情
func (d *ContentDomain) ApiGetContent(c *gin.Context) {
	var req GetContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	content, err := d.GetContent(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to get content: %v", err)
		utils.RespError(c, 500, "failed to get content")
		return
	}

	res := &GetContentRes{
		ID:          content.ID,
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        content.Tags,
		CreatedAt:   content.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiSearchContents 搜索内容
func (d *ContentDomain) ApiSearchContents(c *gin.Context) {
	var req SearchContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	contents, total, err := d.SearchContents(
		c.Request.Context(),
		req.Query,
		req.ContentType,
		req.Status,
		req.Tag,
		req.Page,
		req.Limit,
	)
	if err != nil {
		klog.Errorf("failed to search contents: %v", err)
		utils.RespError(c, 500, "failed to search contents")
		return
	}

	responses := make([]*GetContentRes, len(contents))
	for i, content := range contents {
		responses[i] = &GetContentRes{
			ID:          content.ID,
			Title:       content.Title,
			Body:        content.Body,
			ContentType: content.ContentType,
			Status:      content.Status,
			Tags:        content.Tags,
			CreatedAt:   content.CreatedAt,
		}
	}

	utils.RespSuccess(c, SearchContentRes{Contents: responses, Total: total})
}

// ApiAddContent 添加内容
func (d *ContentDomain) ApiAddContent(c *gin.Context) {
	var req AddContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	content, err := d.AddContent(
		c.Request.Context(),
		req.Title,
		req.Body,
		req.ContentType,
		req.Status,
		req.Tags,
	)
	if err != nil {
		klog.Errorf("failed to add content: %v", err)
		utils.RespError(c, 500, "failed to add content")
		return
	}

	res := &GetContentRes{
		ID:          content.ID,
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        content.Tags,
		CreatedAt:   content.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiUpdateContent 更新内容
func (d *ContentDomain) ApiUpdateContent(c *gin.Context) {
	var req UpdateContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	content, err := d.UpdateContent(
		c.Request.Context(),
		req.ID,
		req.Title,
		req.Body,
		req.ContentType,
		req.Status,
		req.Tags,
	)
	if err != nil {
		klog.Errorf("failed to update content: %v", err)
		utils.RespError(c, 500, "failed to update content")
		return
	}

	res := &GetContentRes{
		ID:          content.ID,
		Title:       content.Title,
		Body:        content.Body,
		ContentType: content.ContentType,
		Status:      content.Status,
		Tags:        content.Tags,
		CreatedAt:   content.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiDeleteContent 删除内容
func (d *ContentDomain) ApiDeleteContent(c *gin.Context) {
	var req DeleteContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	err := d.DeleteContent(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to delete content: %v", err)
		utils.RespError(c, 500, "failed to delete content")
		return
	}

	utils.RespSuccess(c, map[string]int64{"id": req.ID})
}
