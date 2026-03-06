package financing

import "time"

// BusinessPlan 商业计划实体
type BusinessPlan struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Industry        string     `json:"industry"`
	Region          string     `json:"region"`
	FinancingAmount float64    `json:"financing_amount"`
	CompanySize     string     `json:"company_size"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

// 请求/响应结构体

type GetBusinessPlanReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetBusinessPlanRes struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Industry        string     `json:"industry"`
	Region          string     `json:"region"`
	FinancingAmount float64    `json:"financing_amount"`
	CompanySize     string     `json:"company_size"`
	CreatedAt       *time.Time `json:"created_at"`
}

type GetBusinessPlansReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type CreateBusinessPlanReq struct {
	Title           string  `json:"title" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	Industry        string  `json:"industry" binding:"required"`
	Region          string  `json:"region" binding:"required"`
	FinancingAmount float64 `json:"financing_amount" binding:"required"`
	CompanySize     string  `json:"company_size" binding:"required"`
}

type UpdateBusinessPlanReq struct {
	ID              int64    `json:"id" binding:"required"`
	Title           *string  `json:"title"`
	Content         *string  `json:"content"`
	Industry        *string  `json:"industry"`
	Region          *string  `json:"region"`
	FinancingAmount *float64 `json:"financing_amount"`
	CompanySize     *string  `json:"company_size"`
}

type DeleteBusinessPlanReq struct {
	ID int64 `json:"id" binding:"required"`
}
