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

// ===================== Analytics API Handlers =====================

func (d *PublishingDomain) ApiAnalyticsSummary(c *gin.Context) {
	var req AnalyticsSummaryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Days = 30 // 默认 30 天
	}

	// 简单的统计实现：从 publish_tasks 聚合
	// 在生产环境中应该从专业的分析表或业务元数据中获取
	tasks, _, err := d.SearchPublishTasks(c, "", "", "", 1, 1000)
	if err != nil {
		klog.Errorf("failed to fetch publish tasks: %v", err)
		utils.RespError(c, 500, "failed to fetch analytics")
		return
	}

	summary := &AnalyticsSummaryRes{
		TotalPublished:    int64(len(tasks)),
		TotalViews:        0, // TODO: 从真实数据源获取
		TotalLikes:        0,
		TotalComments:     0,
		AvgLikesPerPost:   0,
		TotalNewFollowers: 0,
	}

	utils.RespSuccess(c, summary)
}

func (d *PublishingDomain) ApiAnalyticsRanking(c *gin.Context) {
	var req AnalyticsRankingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Metric = "views"
		req.Days = 30
		req.Limit = 10
	}

	// 简单实现：返回空列表（在实际环境中应聚合内容性能数据）
	resp := &AnalyticsRankingRes{
		Data: []*ContentRankingItem{},
	}

	utils.RespSuccess(c, resp)
}

func (d *PublishingDomain) ApiAnalyticsPlatformComparison(c *gin.Context) {
	// 获取所有发布者
	publishers, _, err := d.SearchPublishers(c, "", nil, 1, 100)
	if err != nil {
		klog.Errorf("failed to fetch publishers: %v", err)
		utils.RespError(c, 500, "failed to fetch platform metrics")
		return
	}

	metrics := make([]*PlatformMetrics, 0)

	// 为每个平台计算统计
	for _, pub := range publishers {
		tasks, _, err := d.SearchPublishTasks(c, "", pub.Platform, "", 1, 1000)
		if err != nil {
			continue
		}

		successful := int64(0)
		failed := int64(0)

		for _, task := range tasks {
			if task.Status == "published" {
				successful++
			} else if task.Status == "failed" {
				failed++
			}
		}

		total := successful + failed
		successRate := 0.0
		if total > 0 {
			successRate = float64(successful) / float64(total) * 100
		}

		metrics = append(metrics, &PlatformMetrics{
			Platform:    pub.Platform,
			Count:       int64(len(tasks)),
			Success:     successful,
			Failed:      failed,
			SuccessRate: successRate,
		})
	}

	utils.RespSuccess(c, AnalyticsPlatformComparisonRes{Data: metrics})
}
