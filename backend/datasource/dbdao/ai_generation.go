package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GeneratorModel struct {
	ID          int64     `gorm:"primarykey;column:id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Provider    string    `gorm:"column:provider;type:varchar(100);not null"`
	Model       string    `gorm:"column:model;type:varchar(100);not null"`
	APIKey      string    `gorm:"column:api_key;type:varchar(500)"`
	APIEndpoint string    `gorm:"column:api_endpoint;type:varchar(500)"`
	Enabled     bool      `gorm:"column:enabled;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (GeneratorModel) TableName() string {
	return "ai_generators"
}

type Generator struct {
	ID          int64
	Name        string
	Provider    string
	Model       string
	APIKey      string
	APIEndpoint string
	Enabled     bool
	CreatedAt   time.Time
}

func (d *DB) GetGenerator(ctx context.Context, id int64) (*Generator, error) {
	var model GeneratorModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("generator not found")
		}
		return nil, result.Error
	}

	return &Generator{
		ID:          model.ID,
		Name:        model.Name,
		Provider:    model.Provider,
		Model:       model.Model,
		APIKey:      model.APIKey,
		APIEndpoint: model.APIEndpoint,
		Enabled:     model.Enabled,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) SearchGenerators(
	ctx context.Context,
	provider string, model string, enabled *bool,
	page, limit int64,
) ([]*Generator, int64, error) {
	var models []GeneratorModel
	q := d.DB().WithContext(ctx)

	if provider != "" {
		q = q.Where("provider = ?", provider)
	}
	if model != "" {
		q = q.Where("model = ?", model)
	}
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}

	var total int64
	q.Model(&GeneratorModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	gens := make([]*Generator, len(models))
	for i, m := range models {
		gens[i] = &Generator{
			ID:          m.ID,
			Name:        m.Name,
			Provider:    m.Provider,
			Model:       m.Model,
			APIKey:      m.APIKey,
			APIEndpoint: m.APIEndpoint,
			Enabled:     m.Enabled,
			CreatedAt:   m.CreatedAt,
		}
	}

	return gens, total, nil
}

func (d *DB) CreateGenerator(ctx context.Context, gen *Generator) (*Generator, error) {
	model := &GeneratorModel{
		Name:        gen.Name,
		Provider:    gen.Provider,
		Model:       gen.Model,
		APIKey:      gen.APIKey,
		APIEndpoint: gen.APIEndpoint,
		Enabled:     gen.Enabled,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &Generator{
		ID:          model.ID,
		Name:        model.Name,
		Provider:    model.Provider,
		Model:       model.Model,
		APIKey:      model.APIKey,
		APIEndpoint: model.APIEndpoint,
		Enabled:     model.Enabled,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) UpdateGenerator(ctx context.Context, gen *Generator) (*Generator, error) {
	updates := map[string]interface{}{}

	if gen.Name != "" {
		updates["name"] = gen.Name
	}
	if gen.Provider != "" {
		updates["provider"] = gen.Provider
	}
	if gen.Model != "" {
		updates["model"] = gen.Model
	}
	if gen.APIKey != "" {
		updates["api_key"] = gen.APIKey
	}
	if gen.APIEndpoint != "" {
		updates["api_endpoint"] = gen.APIEndpoint
	}

	result := d.DB().WithContext(ctx).Model(&GeneratorModel{}).Where("id = ?", gen.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("generator not found")
	}

	return d.GetGenerator(ctx, gen.ID)
}

type PromptTemplateModel struct {
	ID             int64     `gorm:"primarykey;column:id"`
	Name           string    `gorm:"column:name;type:varchar(255);not null"`
	Description    string    `gorm:"column:description;type:varchar(500)"`
	Template       string    `gorm:"column:template;type:longtext;not null"`
	InputVariables string    `gorm:"column:input_variables;type:json"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (PromptTemplateModel) TableName() string {
	return "ai_prompt_templates"
}

type PromptTemplate struct {
	ID             int64
	Name           string
	Description    string
	Template       string
	InputVariables string
	CreatedAt      time.Time
}

func (d *DB) GetPromptTemplate(ctx context.Context, id int64) (*PromptTemplate, error) {
	var model PromptTemplateModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("prompt template not found")
		}
		return nil, result.Error
	}

	return &PromptTemplate{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description,
		Template:       model.Template,
		InputVariables: model.InputVariables,
		CreatedAt:      model.CreatedAt,
	}, nil
}

func (d *DB) SearchPromptTemplates(
	ctx context.Context,
	name, provider string,
	page, limit int64,
) ([]*PromptTemplate, int64, error) {
	var models []PromptTemplateModel
	q := d.DB().WithContext(ctx)

	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	q.Model(&PromptTemplateModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	pts := make([]*PromptTemplate, len(models))
	for i, m := range models {
		pts[i] = &PromptTemplate{
			ID:             m.ID,
			Name:           m.Name,
			Description:    m.Description,
			Template:       m.Template,
			InputVariables: m.InputVariables,
			CreatedAt:      m.CreatedAt,
		}
	}

	return pts, total, nil
}

func (d *DB) CreatePromptTemplate(ctx context.Context, pt *PromptTemplate) (*PromptTemplate, error) {
	model := &PromptTemplateModel{
		Name:           pt.Name,
		Description:    pt.Description,
		Template:       pt.Template,
		InputVariables: pt.InputVariables,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &PromptTemplate{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description,
		Template:       model.Template,
		InputVariables: model.InputVariables,
		CreatedAt:      model.CreatedAt,
	}, nil
}

func (d *DB) UpdatePromptTemplate(ctx context.Context, pt *PromptTemplate) (*PromptTemplate, error) {
	updates := map[string]interface{}{}

	if pt.Name != "" {
		updates["name"] = pt.Name
	}
	if pt.Description != "" {
		updates["description"] = pt.Description
	}
	if pt.Template != "" {
		updates["template"] = pt.Template
	}
	if pt.InputVariables != "" {
		updates["input_variables"] = pt.InputVariables
	}

	result := d.DB().WithContext(ctx).Model(&PromptTemplateModel{}).Where("id = ?", pt.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("prompt template not found")
	}

	return d.GetPromptTemplate(ctx, pt.ID)
}

type GenerateTaskModel struct {
	ID          int64      `gorm:"primarykey;column:id"`
	GeneratorID int64      `gorm:"column:generator_id;index"`
	Input       string     `gorm:"column:input;type:longtext"`
	Output      string     `gorm:"column:output;type:longtext"`
	Status      string     `gorm:"column:status;type:varchar(50)"`
	Error       string     `gorm:"column:error;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime:milli"`
	StartedAt   *time.Time `gorm:"column:started_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

func (GenerateTaskModel) TableName() string {
	return "ai_generate_tasks"
}

type GenerateTask struct {
	ID          int64
	GeneratorID int64
	Input       string
	Output      string
	Status      string
	Error       string
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

func (d *DB) CreateGenerateTask(ctx context.Context, task *GenerateTask) (*GenerateTask, error) {
	model := &GenerateTaskModel{
		GeneratorID: task.GeneratorID,
		Input:       task.Input,
		Status:      task.Status,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &GenerateTask{
		ID:          model.ID,
		GeneratorID: model.GeneratorID,
		Input:       model.Input,
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
	}, nil
}
