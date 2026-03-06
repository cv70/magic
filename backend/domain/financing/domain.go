package financing

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type FinancingDomain struct {
	DB *dbdao.DB
}

func NewFinancingDomain(db *dbdao.DB) *FinancingDomain {
	return &FinancingDomain{DB: db}
}

// GetBusinessPlan 获取商业计划
func (d *FinancingDomain) GetBusinessPlan(ctx context.Context, id int64) (*BusinessPlan, error) {
	plan, err := d.DB.GetBusinessPlan(ctx, id)
	if err != nil {
		return nil, err
	}

	return &BusinessPlan{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       &plan.CreatedAt,
		UpdatedAt:       &plan.UpdatedAt,
	}, nil
}

// GetBusinessPlans 批量获取商业计划
func (d *FinancingDomain) GetBusinessPlans(ctx context.Context, ids []int64) ([]*BusinessPlan, error) {
	if len(ids) == 0 {
		return []*BusinessPlan{}, nil
	}

	plans, err := d.DB.GetBusinessPlans(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]*BusinessPlan, len(plans))
	for i, plan := range plans {
		result[i] = &BusinessPlan{
			ID:              plan.ID,
			Title:           plan.Title,
			Content:         plan.Content,
			Industry:        plan.Industry,
			Region:          plan.Region,
			FinancingAmount: plan.FinancingAmount,
			CompanySize:     plan.CompanySize,
			CreatedAt:       &plan.CreatedAt,
			UpdatedAt:       &plan.UpdatedAt,
		}
	}

	return result, nil
}

// CreateBusinessPlan 创建商业计划
func (d *FinancingDomain) CreateBusinessPlan(
	ctx context.Context,
	title, content, industry, region, companySize string,
	financingAmount float64,
) (*BusinessPlan, error) {
	if title == "" || content == "" || industry == "" || region == "" {
		return nil, errors.New("missing required fields")
	}

	plan := &dbdao.BusinessPlan{
		Title:           title,
		Content:         content,
		Industry:        industry,
		Region:          region,
		FinancingAmount: financingAmount,
		CompanySize:     companySize,
	}

	created, err := d.DB.CreateBusinessPlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &BusinessPlan{
		ID:              created.ID,
		Title:           created.Title,
		Content:         created.Content,
		Industry:        created.Industry,
		Region:          created.Region,
		FinancingAmount: created.FinancingAmount,
		CompanySize:     created.CompanySize,
		CreatedAt:       &created.CreatedAt,
		UpdatedAt:       &created.UpdatedAt,
	}, nil
}

// UpdateBusinessPlan 更新商业计划
func (d *FinancingDomain) UpdateBusinessPlan(
	ctx context.Context,
	id int64,
	title, content, industry, region, companySize *string,
	financingAmount *float64,
) (*BusinessPlan, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	// 构建更新对象，只包含非nil字段
	plan := &dbdao.BusinessPlan{
		ID: id,
	}

	if title != nil {
		plan.Title = *title
	}
	if content != nil {
		plan.Content = *content
	}
	if industry != nil {
		plan.Industry = *industry
	}
	if region != nil {
		plan.Region = *region
	}
	if companySize != nil {
		plan.CompanySize = *companySize
	}
	if financingAmount != nil {
		plan.FinancingAmount = *financingAmount
	}

	updated, err := d.DB.UpdateBusinessPlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &BusinessPlan{
		ID:              updated.ID,
		Title:           updated.Title,
		Content:         updated.Content,
		Industry:        updated.Industry,
		Region:          updated.Region,
		FinancingAmount: updated.FinancingAmount,
		CompanySize:     updated.CompanySize,
		CreatedAt:       &updated.CreatedAt,
		UpdatedAt:       &updated.UpdatedAt,
	}, nil
}

// DeleteBusinessPlan 删除商业计划
func (d *FinancingDomain) DeleteBusinessPlan(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return d.DB.DeleteBusinessPlan(ctx, id)
}
