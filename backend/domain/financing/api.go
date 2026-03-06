package financing

import (
	"backend/utils"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiGetBusinessPlan 获取商业计划详情
func (d *FinancingDomain) ApiGetBusinessPlan(c *gin.Context) {
	var req GetBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.GetBusinessPlan(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to get business plan: %v", err)
		utils.RespError(c, 500, "failed to get business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiGetBusinessPlans 批量获取商业计划
func (d *FinancingDomain) ApiGetBusinessPlans(c *gin.Context) {
	var req GetBusinessPlansReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	plans, err := d.GetBusinessPlans(c.Request.Context(), req.IDs)
	if err != nil {
		klog.Errorf("failed to get business plans: %v", err)
		utils.RespError(c, 500, "failed to get business plans")
		return
	}

	responses := make([]*GetBusinessPlanRes, len(plans))
	for i, plan := range plans {
		responses[i] = &GetBusinessPlanRes{
			ID:              plan.ID,
			Title:           plan.Title,
			Content:         plan.Content,
			Industry:        plan.Industry,
			Region:          plan.Region,
			FinancingAmount: plan.FinancingAmount,
			CompanySize:     plan.CompanySize,
			CreatedAt:       plan.CreatedAt,
		}
	}

	utils.RespSuccess(c, responses)
}

// ApiCreateBusinessPlan 创建商业计划
func (d *FinancingDomain) ApiCreateBusinessPlan(c *gin.Context) {
	var req CreateBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.CreateBusinessPlan(
		c.Request.Context(),
		req.Title,
		req.Content,
		req.Industry,
		req.Region,
		req.CompanySize,
		req.FinancingAmount,
	)
	if err != nil {
		klog.Errorf("failed to create business plan: %v", err)
		utils.RespError(c, 500, "failed to create business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiUpdateBusinessPlan 更新商业计划
func (d *FinancingDomain) ApiUpdateBusinessPlan(c *gin.Context) {
	var req UpdateBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.UpdateBusinessPlan(
		c.Request.Context(),
		req.ID,
		req.Title,
		req.Content,
		req.Industry,
		req.Region,
		req.CompanySize,
		req.FinancingAmount,
	)
	if err != nil {
		klog.Errorf("failed to update business plan: %v", err)
		utils.RespError(c, 500, "failed to update business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	utils.RespSuccess(c, res)
}

// ApiDeleteBusinessPlan 删除商业计划
func (d *FinancingDomain) ApiDeleteBusinessPlan(c *gin.Context) {
	var req DeleteBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	err := d.DeleteBusinessPlan(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to delete business plan: %v", err)
		utils.RespError(c, 500, "failed to delete business plan")
		return
	}

	utils.RespSuccess(c, map[string]int64{"id": req.ID})
}
