package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SystemConfigModel struct {
	ID          int64     `gorm:"primarykey;column:id"`
	Key         string    `gorm:"column:key;type:varchar(255);not null;uniqueIndex"`
	Value       string    `gorm:"column:value;type:text"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	Category    string    `gorm:"column:category;type:varchar(100)"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (SystemConfigModel) TableName() string {
	return "system_configs"
}

type SystemConfig struct {
	ID          int64
	Key         string
	Value       string
	Description string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *DB) GetSystemConfig(ctx context.Context, id int64) (*SystemConfig, error) {
	var model SystemConfigModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("system config not found")
		}
		return nil, result.Error
	}

	return &SystemConfig{
		ID:          model.ID,
		Key:         model.Key,
		Value:       model.Value,
		Description: model.Description,
		Category:    model.Category,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (d *DB) SearchSystemConfigs(ctx context.Context, key, category string, page, limit int64) ([]*SystemConfig, int64, error) {
	var models []SystemConfigModel
	q := d.DB().WithContext(ctx)

	if key != "" {
		q = q.Where("key LIKE ?", "%"+key+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}

	var total int64
	q.Model(&SystemConfigModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	cfgs := make([]*SystemConfig, len(models))
	for i, m := range models {
		cfgs[i] = &SystemConfig{
			ID:          m.ID,
			Key:         m.Key,
			Value:       m.Value,
			Description: m.Description,
			Category:    m.Category,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}

	return cfgs, total, nil
}

func (d *DB) CreateSystemConfig(ctx context.Context, cfg *SystemConfig) (*SystemConfig, error) {
	model := &SystemConfigModel{
		Key:         cfg.Key,
		Value:       cfg.Value,
		Description: cfg.Description,
		Category:    cfg.Category,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &SystemConfig{
		ID:          model.ID,
		Key:         model.Key,
		Value:       model.Value,
		Description: model.Description,
		Category:    model.Category,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (d *DB) UpdateSystemConfig(ctx context.Context, cfg *SystemConfig) (*SystemConfig, error) {
	updates := map[string]interface{}{}

	if cfg.Key != "" {
		updates["key"] = cfg.Key
	}
	if cfg.Value != "" {
		updates["value"] = cfg.Value
	}
	if cfg.Description != "" {
		updates["description"] = cfg.Description
	}
	if cfg.Category != "" {
		updates["category"] = cfg.Category
	}

	result := d.DB().WithContext(ctx).Model(&SystemConfigModel{}).Where("id = ?", cfg.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("system config not found")
	}

	return d.GetSystemConfig(ctx, cfg.ID)
}

type ProviderConfigModel struct {
	ID           int64     `gorm:"primarykey;column:id"`
	ProviderName string    `gorm:"column:provider_name;type:varchar(100);not null;index"`
	ConfigKey    string    `gorm:"column:config_key;type:varchar(255);not null"`
	ConfigValue  string    `gorm:"column:config_value;type:text"`
	Description  string    `gorm:"column:description;type:varchar(500)"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (ProviderConfigModel) TableName() string {
	return "provider_configs"
}

type ProviderConfig struct {
	ID           int64
	ProviderName string
	ConfigKey    string
	ConfigValue  string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (d *DB) GetProviderConfig(ctx context.Context, id int64) (*ProviderConfig, error) {
	var model ProviderConfigModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("provider config not found")
		}
		return nil, result.Error
	}

	return &ProviderConfig{
		ID:           model.ID,
		ProviderName: model.ProviderName,
		ConfigKey:    model.ConfigKey,
		ConfigValue:  model.ConfigValue,
		Description:  model.Description,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

func (d *DB) SearchProviderConfigs(ctx context.Context, providerName, configKey string, page, limit int64) ([]*ProviderConfig, int64, error) {
	var models []ProviderConfigModel
	q := d.DB().WithContext(ctx)

	if providerName != "" {
		q = q.Where("provider_name = ?", providerName)
	}
	if configKey != "" {
		q = q.Where("config_key LIKE ?", "%"+configKey+"%")
	}

	var total int64
	q.Model(&ProviderConfigModel{}).Count(&total)

	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	cfgs := make([]*ProviderConfig, len(models))
	for i, m := range models {
		cfgs[i] = &ProviderConfig{
			ID:           m.ID,
			ProviderName: m.ProviderName,
			ConfigKey:    m.ConfigKey,
			ConfigValue:  m.ConfigValue,
			Description:  m.Description,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}
	}

	return cfgs, total, nil
}

func (d *DB) CreateProviderConfig(ctx context.Context, cfg *ProviderConfig) (*ProviderConfig, error) {
	model := &ProviderConfigModel{
		ProviderName: cfg.ProviderName,
		ConfigKey:    cfg.ConfigKey,
		ConfigValue:  cfg.ConfigValue,
		Description:  cfg.Description,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ProviderConfig{
		ID:           model.ID,
		ProviderName: model.ProviderName,
		ConfigKey:    model.ConfigKey,
		ConfigValue:  model.ConfigValue,
		Description:  model.Description,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

func (d *DB) UpdateProviderConfig(ctx context.Context, cfg *ProviderConfig) (*ProviderConfig, error) {
	updates := map[string]interface{}{}

	if cfg.ProviderName != "" {
		updates["provider_name"] = cfg.ProviderName
	}
	if cfg.ConfigKey != "" {
		updates["config_key"] = cfg.ConfigKey
	}
	if cfg.ConfigValue != "" {
		updates["config_value"] = cfg.ConfigValue
	}
	if cfg.Description != "" {
		updates["description"] = cfg.Description
	}

	result := d.DB().WithContext(ctx).Model(&ProviderConfigModel{}).Where("id = ?", cfg.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("provider config not found")
	}

	return d.GetProviderConfig(ctx, cfg.ID)
}
