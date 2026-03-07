package ai_generation

import (
	"time"
)

type Generator struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Provider    string     `json:"provider"`
	Model       string     `json:"model"`
	APIKey      string     `json:"api_key"`
	APIEndpoint string     `json:"api_endpoint"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   *time.Time `json:"created_at"`
}

type PromptTemplate struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Template       string     `json:"template"`
	InputVariables string     `json:"input_variables"`
	CreatedAt      *time.Time `json:"created_at"`
}

type GenerateTask struct {
	ID          int64      `json:"id"`
	GeneratorID int64      `json:"generator_id"`
	Input       string     `json:"input"`
	Output      string     `json:"output"`
	Status      string     `json:"status"`
	Error       string     `json:"error"`
	CreatedAt   *time.Time `json:"created_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}


type GetGeneratorReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetGeneratorRes struct {
	Generator *Generator `json:"generator"`
}

type SearchGeneratorsReq struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Enabled  *bool  `json:"enabled"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type SearchGeneratorsRes struct {
	Generators []*Generator `json:"generators"`
	Total      int64        `json:"total"`
}

type AddGeneratorReq struct {
	Name        string `json:"name" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	Model       string `json:"model" binding:"required"`
	APIKey      string `json:"api_key"`
	APIEndpoint string `json:"api_endpoint"`
	Enabled     *bool  `json:"enabled"`
}

type AddGeneratorRes struct {
	ID int64 `json:"id"`
}

type UpdateGeneratorReq struct {
	ID          int64   `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Provider    *string `json:"provider"`
	Model       *string `json:"model"`
	APIKey      *string `json:"api_key"`
	APIEndpoint *string `json:"api_endpoint"`
	Enabled     *bool   `json:"enabled"`
}

type UpdateGeneratorRes struct {
	ID int64 `json:"id"`
}

type GenerateContentReq struct {
	GeneratorID int64  `json:"generator_id" binding:"required"`
	Input       string `json:"input" binding:"required"`
}

type GenerateContentRes struct {
	ID int64 `json:"id"`
}

type GetPromptTemplateReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetPromptTemplateRes struct {
	PromptTemplate *PromptTemplate `json:"prompt_template"`
}

type SearchPromptTemplatesReq struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type SearchPromptTemplatesRes struct {
	PromptTemplates []*PromptTemplate `json:"prompt_templates"`
	Total           int64             `json:"total"`
}

type AddPromptTemplateReq struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Template       string `json:"template" binding:"required"`
	InputVariables string `json:"input_variables"`
}

type AddPromptTemplateRes struct {
	ID int64 `json:"id"`
}

type UpdatePromptTemplateReq struct {
	ID             int64   `json:"id" binding:"required"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Template       *string `json:"template"`
	InputVariables *string `json:"input_variables"`
}

type UpdatePromptTemplateRes struct {
	ID int64 `json:"id"`
}

// ===================== Generate Task API =====================

type GetGenerateTaskReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetGenerateTaskRes struct {
	GenerateTask *GenerateTask `json:"generate_task"`
}

