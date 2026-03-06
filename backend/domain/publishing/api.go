package publishing

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func (d *PublishingDomain) ApiGetPublisher(c *gin.Context) {
	var req GetPublisherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	publisher, err := d.GetPublisher(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get publisher: %v", err)
		utils.RespError(c, 500, "failed to get publisher")
		return
	}

	utils.RespSuccess(c, GetPublisherRes{Publisher: publisher})
}

func (d *PublishingDomain) ApiSearchPublishers(c *gin.Context) {
	var req SearchPublishersReq
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

	publishers, total, err := d.SearchPublishers(c, req.Platform, req.Enabled, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search publishers: %v", err)
		utils.RespError(c, 500, "failed to search publishers")
		return
	}

	utils.RespSuccess(c, SearchPublishersRes{Publishers: publishers, Total: total})
}

func (d *PublishingDomain) ApiAddPublisher(c *gin.Context) {
	var req AddPublisherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	publisher, err := d.AddPublisher(c, req.Name, req.Platform, req.PlatformID, enabled)
	if err != nil {
		klog.Errorf("failed to add publisher: %v", err)
		utils.RespError(c, 500, "failed to add publisher")
		return
	}

	utils.RespSuccess(c, AddPublisherRes{ID: publisher.ID})
}

func (d *PublishingDomain) ApiUpdatePublisher(c *gin.Context) {
	var req UpdatePublisherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	_, err := d.UpdatePublisher(c, req.ID, req.Name, req.Platform, req.PlatformID, req.Enabled)
	if err != nil {
		klog.Errorf("failed to update publisher: %v", err)
		utils.RespError(c, 500, "failed to update publisher")
		return
	}

	utils.RespSuccess(c, UpdatePublisherRes{ID: req.ID})
}

func (d *PublishingDomain) ApiPublishContent(c *gin.Context) {
	var req PublishContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	taskID, err := d.PublishContent(c, req.PublisherID, req.ContentID)
	if err != nil {
		klog.Errorf("failed to publish content: %v", err)
		utils.RespError(c, 500, "failed to publish content")
		return
	}

	utils.RespSuccess(c, PublishContentRes{ID: taskID})
}

func (d *PublishingDomain) ApiGetPublishTask(c *gin.Context) {
	var req GetPublishTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	task, err := d.GetPublishTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get publish task: %v", err)
		utils.RespError(c, 500, "failed to get publish task")
		return
	}

	utils.RespSuccess(c, GetPublishTaskRes{PublishTask: task})
}

func (d *PublishingDomain) ApiSearchPublishTasks(c *gin.Context) {
	var req SearchPublishTasksReq
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

	tasks, total, err := d.SearchPublishTasks(c, req.Status, req.Platform, req.CreatedAt, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search publish tasks: %v", err)
		utils.RespError(c, 500, "failed to search publish tasks")
		return
	}

	utils.RespSuccess(c, SearchPublishTasksRes{PublishTasks: tasks, Total: total})
}

func (d *PublishingDomain) ApiGetPublishLog(c *gin.Context) {
	var req GetPublishLogReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	log, err := d.GetPublishLog(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get publish log: %v", err)
		utils.RespError(c, 500, "failed to get publish log")
		return
	}

	utils.RespSuccess(c, GetPublishLogRes{PublishLog: log})
}

func (d *PublishingDomain) ApiSearchPublishLogs(c *gin.Context) {
	var req SearchPublishLogsReq
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

	logs, total, err := d.SearchPublishLogs(c, req.PublishTaskID, req.LogType, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search publish logs: %v", err)
		utils.RespError(c, 500, "failed to search publish logs")
		return
	}

	utils.RespSuccess(c, SearchPublishLogsRes{PublishLogs: logs, Total: total})
}
