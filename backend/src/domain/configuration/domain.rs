#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::configuration::schema::{ProviderConfig, SystemConfig},
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct ConfigurationDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl ConfigurationDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_system_config(&self, id: i64) -> Option<SystemConfig> {
        self.db_dao
            .get_system_config(id)
            .await
            .ok()
            .map(|record| SystemConfig {
                id: record.id,
                key: record.key,
                value: record.value,
                description: record.description,
                category: record.category,
                created_at: record.created_at,
                updated_at: None,
            })
    }

    pub async fn search_system_configs(
        &self,
        key: Option<String>,
        category: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<SystemConfig>, i64) {
        match self
            .db_dao
            .search_system_configs(key, category, page, limit)
            .await
        {
            Ok(records) => {
                let configs: Vec<SystemConfig> = records
                    .into_iter()
                    .map(|record| SystemConfig {
                        id: record.id,
                        key: record.key,
                        value: record.value,
                        description: record.description,
                        category: record.category,
                        created_at: record.created_at,
                        updated_at: None,
                    })
                    .collect();
                let total = configs.len() as i64;
                (configs, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_system_config(
        &self,
        key: String,
        value: String,
        description: Option<String>,
        category: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_system_config(key.as_str(), value.as_str(), description, category)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_system_config(
        &self,
        id: i64,
        key: Option<String>,
        value: Option<String>,
        description: Option<String>,
        category: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_system_config(id, key, value, description, category)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn get_provider_config(&self, id: i64) -> Option<ProviderConfig> {
        self.db_dao
            .get_provider_config(id)
            .await
            .ok()
            .map(|record| ProviderConfig {
                id: record.id,
                provider_name: record.provider_name,
                config_key: record.config_key,
                config_value: record.config_value,
                description: record.description,
                created_at: record.created_at,
                updated_at: None,
            })
    }

    pub async fn search_provider_configs(
        &self,
        provider_name: Option<String>,
        config_key: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<ProviderConfig>, i64) {
        match self
            .db_dao
            .search_provider_configs(provider_name, config_key, page, limit)
            .await
        {
            Ok(records) => {
                let configs: Vec<ProviderConfig> = records
                    .into_iter()
                    .map(|record| ProviderConfig {
                        id: record.id,
                        provider_name: record.provider_name,
                        config_key: record.config_key,
                        config_value: record.config_value,
                        description: record.description,
                        created_at: record.created_at,
                        updated_at: None,
                    })
                    .collect();
                let total = configs.len() as i64;
                (configs, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_provider_config(
        &self,
        provider_name: String,
        config_key: String,
        config_value: String,
        description: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_provider_config(
                provider_name.as_str(),
                config_key.as_str(),
                config_value.as_str(),
                description,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_provider_config(
        &self,
        id: i64,
        provider_name: Option<String>,
        config_key: Option<String>,
        config_value: Option<String>,
        description: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_provider_config(id, provider_name, config_key, config_value, description)
            .await
            .map_err(|e| e.to_string())
    }
}
