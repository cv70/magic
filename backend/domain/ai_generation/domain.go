package ai_generation

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type AIGenerationDomain struct {
	DB *dbdao.DB
}

func NewAIGenerationDomain(db *dbdao.DB) *AIGenerationDomain {
	return &AIGenerationDomain{DB: db}
}

func (d *AIGenerationDomain) GetGenerator(ctx context.Context, id int64) (*Generator, error) {
	gen, err := d.DB.GetGenerator(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Generator{
		ID:          gen.ID,
		Name:        gen.Name,
		Provider:    gen.Provider,
		Model:       gen.Model,
		APIKey:      gen.APIKey,
		APIEndpoint: gen.APIEndpoint,
		Enabled:     gen.Enabled,
		CreatedAt:   &gen.CreatedAt,
	}, nil
}

func (d *AIGenerationDomain) SearchGenerators(
	ctx context.Context,
	provider string, model string, enabled *bool,
	page, limit int64,
) ([]*Generator, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	gens, total, err := d.DB.SearchGenerators(ctx, provider, model, enabled, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*Generator, len(gens))
	for i, g := range gens {
		result[i] = &Generator{
			ID:          g.ID,
			Name:        g.Name,
			Provider:    g.Provider,
			Model:       g.Model,
			APIKey:      g.APIKey,
			APIEndpoint: g.APIEndpoint,
			Enabled:     g.Enabled,
			CreatedAt:   &g.CreatedAt,
		}
	}

	return result, total, nil
}

func (d *AIGenerationDomain) AddGenerator(
	ctx context.Context,
	name, provider, model, apiKey, apiEndpoint string,
	enabled bool,
) (*Generator, error) {
	if name == "" || provider == "" || model == "" {
		return nil, errors.New("missing required fields")
	}

	gen := &dbdao.Generator{
		Name:        name,
		Provider:    provider,
		Model:       model,
		APIKey:      apiKey,
		APIEndpoint: apiEndpoint,
		Enabled:     enabled,
	}

	created, err := d.DB.CreateGenerator(ctx, gen)
	if err != nil {
		return nil, err
	}

	return &Generator{
		ID:          created.ID,
		Name:        created.Name,
		Provider:    created.Provider,
		Model:       created.Model,
		APIKey:      created.APIKey,
		APIEndpoint: created.APIEndpoint,
		Enabled:     created.Enabled,
		CreatedAt:   &created.CreatedAt,
	}, nil
}

func (d *AIGenerationDomain) UpdateGenerator(
	ctx context.Context,
	id int64,
	name, provider, model, apiKey, apiEndpoint *string,
	enabled *bool,
) (*Generator, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	gen := &dbdao.Generator{ID: id}

	if name != nil {
		gen.Name = *name
	}
	if provider != nil {
		gen.Provider = *provider
	}
	if model != nil {
		gen.Model = *model
	}
	if apiKey != nil {
		gen.APIKey = *apiKey
	}
	if apiEndpoint != nil {
		gen.APIEndpoint = *apiEndpoint
	}
	if enabled != nil {
		gen.Enabled = *enabled
	}

	updated, err := d.DB.UpdateGenerator(ctx, gen)
	if err != nil {
		return nil, err
	}

	return &Generator{
		ID:          updated.ID,
		Name:        updated.Name,
		Provider:    updated.Provider,
		Model:       updated.Model,
		APIKey:      updated.APIKey,
		APIEndpoint: updated.APIEndpoint,
		Enabled:     updated.Enabled,
		CreatedAt:   &updated.CreatedAt,
	}, nil
}

func (d *AIGenerationDomain) GenerateContent(ctx context.Context, generatorID int64, input string) (int64, error) {
	if generatorID <= 0 || input == "" {
		return 0, errors.New("invalid generator_id or input")
	}

	task := &dbdao.GenerateTask{
		GeneratorID: generatorID,
		Input:       input,
		Status:      "pending",
	}

	created, err := d.DB.CreateGenerateTask(ctx, task)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (d *AIGenerationDomain) GetPromptTemplate(ctx context.Context, id int64) (*PromptTemplate, error) {
	pt, err := d.DB.GetPromptTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	return &PromptTemplate{
		ID:             pt.ID,
		Name:           pt.Name,
		Description:    pt.Description,
		Template:       pt.Template,
		InputVariables: pt.InputVariables,
		CreatedAt:      &pt.CreatedAt,
	}, nil
}

func (d *AIGenerationDomain) SearchPromptTemplates(
	ctx context.Context,
	name, provider string,
	page, limit int64,
) ([]*PromptTemplate, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	pts, total, err := d.DB.SearchPromptTemplates(ctx, name, provider, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*PromptTemplate, len(pts))
	for i, p := range pts {
		result[i] = &PromptTemplate{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			Template:       p.Template,
			InputVariables: p.InputVariables,
			CreatedAt:      &p.CreatedAt,
		}
	}

	return result, total, nil
}

func (d *AIGenerationDomain) AddPromptTemplate(
	ctx context.Context,
	name, description, template, inputVariables string,
) (*PromptTemplate, error) {
	if name == "" || template == "" {
		return nil, errors.New("missing required fields")
	}

	pt := &dbdao.PromptTemplate{
		Name:           name,
		Description:    description,
		Template:       template,
		InputVariables: inputVariables,
	}

	created, err := d.DB.CreatePromptTemplate(ctx, pt)
	if err != nil {
		return nil, err
	}

	return &PromptTemplate{
		ID:             created.ID,
		Name:           created.Name,
		Description:    created.Description,
		Template:       created.Template,
		InputVariables: created.InputVariables,
		CreatedAt:      &created.CreatedAt,
	}, nil
}

func (d *AIGenerationDomain) UpdatePromptTemplate(
	ctx context.Context,
	id int64,
	name, description, template, inputVariables *string,
) (*PromptTemplate, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	pt := &dbdao.PromptTemplate{ID: id}

	if name != nil {
		pt.Name = *name
	}
	if description != nil {
		pt.Description = *description
	}
	if template != nil {
		pt.Template = *template
	}
	if inputVariables != nil {
		pt.InputVariables = *inputVariables
	}

	updated, err := d.DB.UpdatePromptTemplate(ctx, pt)
	if err != nil {
		return nil, err
	}

	return &PromptTemplate{
		ID:             updated.ID,
		Name:           updated.Name,
		Description:    updated.Description,
		Template:       updated.Template,
		InputVariables: updated.InputVariables,
		CreatedAt:      &updated.CreatedAt,
	}, nil
}
