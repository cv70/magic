package configuration

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type ConfigurationDomain struct {
	DB *dbdao.DB
}

func NewConfigurationDomain(db *dbdao.DB) *ConfigurationDomain {
	return &ConfigurationDomain{DB: db}
}

func (d *ConfigurationDomain) GetSystemConfig(ctx context.Context, id int64) (*SystemConfig, error) {
	cfg, err := d.DB.GetSystemConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	return &SystemConfig{
		ID:          cfg.ID,
		Key:         cfg.Key,
		Value:       cfg.Value,
		Description: cfg.Description,
		Category:    cfg.Category,
		CreatedAt:   &cfg.CreatedAt,
		UpdatedAt:   &cfg.UpdatedAt,
	}, nil
}

func (d *ConfigurationDomain) SearchSystemConfigs(
	ctx context.Context,
	key, category string,
	page, limit int64,
) ([]*SystemConfig, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	cfgs, total, err := d.DB.SearchSystemConfigs(ctx, key, category, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*SystemConfig, len(cfgs))
	for i, c := range cfgs {
		result[i] = &SystemConfig{
			ID:          c.ID,
			Key:         c.Key,
			Value:       c.Value,
			Description: c.Description,
			Category:    c.Category,
			CreatedAt:   &c.CreatedAt,
			UpdatedAt:   &c.UpdatedAt,
		}
	}

	return result, total, nil
}

func (d *ConfigurationDomain) AddSystemConfig(
	ctx context.Context,
	key, value, description, category string,
) (*SystemConfig, error) {
	if key == "" || value == "" {
		return nil, errors.New("missing required fields")
	}

	cfg := &dbdao.SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
		Category:    category,
	}

	created, err := d.DB.CreateSystemConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &SystemConfig{
		ID:          created.ID,
		Key:         created.Key,
		Value:       created.Value,
		Description: created.Description,
		Category:    created.Category,
		CreatedAt:   &created.CreatedAt,
		UpdatedAt:   &created.UpdatedAt,
	}, nil
}

func (d *ConfigurationDomain) UpdateSystemConfig(
	ctx context.Context,
	id int64,
	key, value, description, category *string,
) (*SystemConfig, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	cfg := &dbdao.SystemConfig{ID: id}

	if key != nil {
		cfg.Key = *key
	}
	if value != nil {
		cfg.Value = *value
	}
	if description != nil {
		cfg.Description = *description
	}
	if category != nil {
		cfg.Category = *category
	}

	updated, err := d.DB.UpdateSystemConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &SystemConfig{
		ID:          updated.ID,
		Key:         updated.Key,
		Value:       updated.Value,
		Description: updated.Description,
		Category:    updated.Category,
		CreatedAt:   &updated.CreatedAt,
		UpdatedAt:   &updated.UpdatedAt,
	}, nil
}

func (d *ConfigurationDomain) GetProviderConfig(ctx context.Context, id int64) (*ProviderConfig, error) {
	cfg, err := d.DB.GetProviderConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ProviderConfig{
		ID:           cfg.ID,
		ProviderName: cfg.ProviderName,
		ConfigKey:    cfg.ConfigKey,
		ConfigValue:  cfg.ConfigValue,
		Description:  cfg.Description,
		CreatedAt:    &cfg.CreatedAt,
		UpdatedAt:    &cfg.UpdatedAt,
	}, nil
}

func (d *ConfigurationDomain) SearchProviderConfigs(
	ctx context.Context,
	providerName, configKey string,
	page, limit int64,
) ([]*ProviderConfig, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	cfgs, total, err := d.DB.SearchProviderConfigs(ctx, providerName, configKey, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*ProviderConfig, len(cfgs))
	for i, c := range cfgs {
		result[i] = &ProviderConfig{
			ID:           c.ID,
			ProviderName: c.ProviderName,
			ConfigKey:    c.ConfigKey,
			ConfigValue:  c.ConfigValue,
			Description:  c.Description,
			CreatedAt:    &c.CreatedAt,
			UpdatedAt:    &c.UpdatedAt,
		}
	}

	return result, total, nil
}

func (d *ConfigurationDomain) AddProviderConfig(
	ctx context.Context,
	providerName, configKey, configValue, description string,
) (*ProviderConfig, error) {
	if providerName == "" || configKey == "" || configValue == "" {
		return nil, errors.New("missing required fields")
	}

	cfg := &dbdao.ProviderConfig{
		ProviderName: providerName,
		ConfigKey:    configKey,
		ConfigValue:  configValue,
		Description:  description,
	}

	created, err := d.DB.CreateProviderConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &ProviderConfig{
		ID:           created.ID,
		ProviderName: created.ProviderName,
		ConfigKey:    created.ConfigKey,
		ConfigValue:  created.ConfigValue,
		Description:  created.Description,
		CreatedAt:    &created.CreatedAt,
		UpdatedAt:    &created.UpdatedAt,
	}, nil
}

func (d *ConfigurationDomain) UpdateProviderConfig(
	ctx context.Context,
	id int64,
	providerName, configKey, configValue, description *string,
) (*ProviderConfig, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	cfg := &dbdao.ProviderConfig{ID: id}

	if providerName != nil {
		cfg.ProviderName = *providerName
	}
	if configKey != nil {
		cfg.ConfigKey = *configKey
	}
	if configValue != nil {
		cfg.ConfigValue = *configValue
	}
	if description != nil {
		cfg.Description = *description
	}

	updated, err := d.DB.UpdateProviderConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &ProviderConfig{
		ID:           updated.ID,
		ProviderName: updated.ProviderName,
		ConfigKey:    updated.ConfigKey,
		ConfigValue:  updated.ConfigValue,
		Description:  updated.Description,
		CreatedAt:    &updated.CreatedAt,
		UpdatedAt:    &updated.UpdatedAt,
	}, nil
}
