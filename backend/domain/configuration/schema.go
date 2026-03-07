package configuration

import (
	"time"
)

type SystemConfig struct {
	ID          int64      `json:"id"`
	Key         string     `json:"key"`
	Value       string     `json:"value"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProviderConfig struct {
	ID           int64      `json:"id"`
	ProviderName string     `json:"provider_name"`
	ConfigKey    string     `json:"config_key"`
	ConfigValue  string     `json:"config_value"`
	Description  string     `json:"description"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type GetSystemConfigReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetSystemConfigRes struct {
	SystemConfig *SystemConfig `json:"system_config"`
}

type SearchSystemConfigsReq struct {
	Key      string `json:"key"`
	Category string `json:"category"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type SearchSystemConfigsRes struct {
	SystemConfigs []*SystemConfig `json:"system_configs"`
	Total         int64           `json:"total"`
}

type AddSystemConfigReq struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type AddSystemConfigRes struct {
	ID int64 `json:"id"`
}

type UpdateSystemConfigReq struct {
	ID          int64   `json:"id" binding:"required"`
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
}

type UpdateSystemConfigRes struct {
	ID int64 `json:"id"`
}

type GetProviderConfigReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetProviderConfigRes struct {
	ProviderConfig *ProviderConfig `json:"provider_config"`
}

type SearchProviderConfigsReq struct {
	ProviderName string `json:"provider_name"`
	ConfigKey    string `json:"config_key"`
	Page         int64  `json:"page"`
	Limit        int64  `json:"limit"`
}

type SearchProviderConfigsRes struct {
	ProviderConfigs []*ProviderConfig `json:"provider_configs"`
	Total           int64             `json:"total"`
}

type AddProviderConfigReq struct {
	ProviderName string `json:"provider_name" binding:"required"`
	ConfigKey    string `json:"config_key" binding:"required"`
	ConfigValue  string `json:"config_value" binding:"required"`
	Description  string `json:"description"`
}

type AddProviderConfigRes struct {
	ID int64 `json:"id"`
}

type UpdateProviderConfigReq struct {
	ID           int64   `json:"id" binding:"required"`
	ProviderName *string `json:"provider_name"`
	ConfigKey    *string `json:"config_key"`
	ConfigValue  *string `json:"config_value"`
	Description  *string `json:"description"`
}

type UpdateProviderConfigRes struct {
	ID int64 `json:"id"`
}
