#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Configuration domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SystemConfig {
    pub id: i64,
    pub key: String,
    pub value: Option<String>,
    pub description: Option<String>,
    pub category: Option<String>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ProviderConfig {
    pub id: i64,
    pub provider_name: String,
    pub config_key: String,
    pub config_value: Option<String>,
    pub description: Option<String>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetSystemConfigReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetSystemConfigRes {
    pub system_config: SystemConfig,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchSystemConfigsReq {
    pub key: Option<String>,
    pub category: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchSystemConfigsRes {
    pub system_configs: Vec<SystemConfig>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddSystemConfigReq {
    pub key: String,
    pub value: String,
    pub description: Option<String>,
    pub category: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddSystemConfigRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateSystemConfigReq {
    pub id: i64,
    pub key: Option<String>,
    pub value: Option<String>,
    pub description: Option<String>,
    pub category: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateSystemConfigRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetProviderConfigReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetProviderConfigRes {
    pub provider_config: ProviderConfig,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchProviderConfigsReq {
    pub provider_name: Option<String>,
    pub config_key: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchProviderConfigsRes {
    pub provider_configs: Vec<ProviderConfig>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddProviderConfigReq {
    pub provider_name: String,
    pub config_key: String,
    pub config_value: String,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddProviderConfigRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateProviderConfigReq {
    pub id: i64,
    pub provider_name: Option<String>,
    pub config_key: Option<String>,
    pub config_value: Option<String>,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateProviderConfigRes {
    pub id: i64,
}
