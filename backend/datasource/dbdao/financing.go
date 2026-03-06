package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// BusinessPlanModel 数据库模型
type BusinessPlanModel struct {
	ID              int64     `gorm:"primarykey;column:id"`
	Title           string    `gorm:"column:title;type:varchar(255);not null"`
	Content         string    `gorm:"column:content;type:longtext"`
	Industry        string    `gorm:"column:industry;type:varchar(100);not null"`
	Region          string    `gorm:"column:region;type:varchar(100);not null"`
	FinancingAmount float64   `gorm:"column:financing_amount;type:decimal(15,2)"`
	CompanySize     string    `gorm:"column:company_size;type:varchar(50)"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
}

// TableName 指定表名
func (BusinessPlanModel) TableName() string {
	return "business_plans"
}

// BusinessPlan DAO层返回的业务计划数据结构
type BusinessPlan struct {
	ID              int64
	Title           string
	Content         string
	Industry        string
	Region          string
	FinancingAmount float64
	CompanySize     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GetBusinessPlan 获取商业计划
func (d *DB) GetBusinessPlan(ctx context.Context, id int64) (*BusinessPlan, error) {
	var model BusinessPlanModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("business plan not found")
		}
		return nil, result.Error
	}

	return &BusinessPlan{
		ID:              model.ID,
		Title:           model.Title,
		Content:         model.Content,
		Industry:        model.Industry,
		Region:          model.Region,
		FinancingAmount: model.FinancingAmount,
		CompanySize:     model.CompanySize,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}, nil
}

// GetBusinessPlans 批量获取商业计划
func (d *DB) GetBusinessPlans(ctx context.Context, ids []int64) ([]*BusinessPlan, error) {
	var models []BusinessPlanModel
	result := d.DB().WithContext(ctx).Where("id IN ?", ids).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	plans := make([]*BusinessPlan, len(models))
	for i, model := range models {
		plans[i] = &BusinessPlan{
			ID:              model.ID,
			Title:           model.Title,
			Content:         model.Content,
			Industry:        model.Industry,
			Region:          model.Region,
			FinancingAmount: model.FinancingAmount,
			CompanySize:     model.CompanySize,
			CreatedAt:       model.CreatedAt,
			UpdatedAt:       model.UpdatedAt,
		}
	}

	return plans, nil
}

// CreateBusinessPlan 创建商业计划
func (d *DB) CreateBusinessPlan(ctx context.Context, plan *BusinessPlan) (*BusinessPlan, error) {
	model := &BusinessPlanModel{
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &BusinessPlan{
		ID:              model.ID,
		Title:           model.Title,
		Content:         model.Content,
		Industry:        model.Industry,
		Region:          model.Region,
		FinancingAmount: model.FinancingAmount,
		CompanySize:     model.CompanySize,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}, nil
}

// UpdateBusinessPlan 更新商业计划
func (d *DB) UpdateBusinessPlan(ctx context.Context, plan *BusinessPlan) (*BusinessPlan, error) {
	updates := map[string]interface{}{}

	if plan.Title != "" {
		updates["title"] = plan.Title
	}
	if plan.Content != "" {
		updates["content"] = plan.Content
	}
	if plan.Industry != "" {
		updates["industry"] = plan.Industry
	}
	if plan.Region != "" {
		updates["region"] = plan.Region
	}
	if plan.FinancingAmount > 0 {
		updates["financing_amount"] = plan.FinancingAmount
	}
	if plan.CompanySize != "" {
		updates["company_size"] = plan.CompanySize
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := d.DB().WithContext(ctx).Model(&BusinessPlanModel{}).Where("id = ?", plan.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("business plan not found")
	}

	// 获取更新后的记录
	return d.GetBusinessPlan(ctx, plan.ID)
}

// DeleteBusinessPlan 删除商业计划
func (d *DB) DeleteBusinessPlan(ctx context.Context, id int64) error {
	result := d.DB().WithContext(ctx).Where("id = ?", id).Delete(&BusinessPlanModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("business plan not found")
	}

	return nil
}
