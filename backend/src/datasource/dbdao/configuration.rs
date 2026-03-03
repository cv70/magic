#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct SystemConfigRecord {
    pub id: i64,
    pub key: String,
    pub value: Option<String>,
    pub description: Option<String>,
    pub category: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct ProviderConfigRecord {
    pub id: i64,
    pub provider_name: String,
    pub config_key: String,
    pub config_value: Option<String>,
    pub description: Option<String>,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get system config by ID
    pub async fn get_system_config(&self, id: i64) -> Result<SystemConfigRecord> {
        let row = sqlx::query(
            "SELECT id, key, value, description, category, created_at
             FROM system_config WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let config = SystemConfigRecord {
            id: row.get("id"),
            key: row.get("key"),
            value: row.get("value"),
            description: row.get("description"),
            category: row.get("category"),
            created_at: row.get("created_at"),
        };

        Ok(config)
    }

    /// Search system configs with pagination
    pub async fn search_system_configs(
        &self,
        key: Option<String>,
        category: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<SystemConfigRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, key, value, description, category, created_at
             FROM system_config
             WHERE ($1 IS NULL OR key = $1)
               AND ($2 IS NULL OR category = $2)
             OFFSET $3 LIMIT $4",
        )
        .bind(key)
        .bind(category)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let configs = rows
            .into_iter()
            .map(|row| SystemConfigRecord {
                id: row.get("id"),
                key: row.get("key"),
                value: row.get("value"),
                description: row.get("description"),
                category: row.get("category"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(configs)
    }

    /// Add system config
    pub async fn add_system_config(
        &self,
        key: &str,
        value: &str,
        description: Option<String>,
        category: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO system_config (key, value, description, category)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(key)
        .bind(value)
        .bind(description)
        .bind(category)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update system config
    pub async fn update_system_config(
        &self,
        id: i64,
        key: Option<String>,
        value: Option<String>,
        description: Option<String>,
        category: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE system_config
             SET key = COALESCE($1, key),
                 value = COALESCE($2, value),
                 description = COALESCE($3, description),
                 category = COALESCE($4, category)
             WHERE id = $5
             RETURNING id",
        )
        .bind(key)
        .bind(value)
        .bind(description)
        .bind(category)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Get provider config by ID
    pub async fn get_provider_config(&self, id: i64) -> Result<ProviderConfigRecord> {
        let row = sqlx::query(
            "SELECT id, provider_name, config_key, config_value, description, created_at
             FROM provider_config WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let config = ProviderConfigRecord {
            id: row.get("id"),
            provider_name: row.get("provider_name"),
            config_key: row.get("config_key"),
            config_value: row.get("config_value"),
            description: row.get("description"),
            created_at: row.get("created_at"),
        };

        Ok(config)
    }

    /// Search provider configs with pagination
    pub async fn search_provider_configs(
        &self,
        provider_name: Option<String>,
        config_key: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<ProviderConfigRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, provider_name, config_key, config_value, description, created_at
             FROM provider_config
             WHERE ($1 IS NULL OR provider_name = $1)
               AND ($2 IS NULL OR config_key = $2)
             OFFSET $3 LIMIT $4",
        )
        .bind(provider_name)
        .bind(config_key)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let configs = rows
            .into_iter()
            .map(|row| ProviderConfigRecord {
                id: row.get("id"),
                provider_name: row.get("provider_name"),
                config_key: row.get("config_key"),
                config_value: row.get("config_value"),
                description: row.get("description"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(configs)
    }

    /// Add provider config
    pub async fn add_provider_config(
        &self,
        provider_name: &str,
        config_key: &str,
        config_value: &str,
        description: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO provider_config (provider_name, config_key, config_value, description)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(provider_name)
        .bind(config_key)
        .bind(config_value)
        .bind(description)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update provider config
    pub async fn update_provider_config(
        &self,
        id: i64,
        provider_name: Option<String>,
        config_key: Option<String>,
        config_value: Option<String>,
        description: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE provider_config
             SET provider_name = COALESCE($1, provider_name),
                 config_key = COALESCE($2, config_key),
                 config_value = COALESCE($3, config_value),
                 description = COALESCE($4, description)
             WHERE id = $5
             RETURNING id",
        )
        .bind(provider_name)
        .bind(config_key)
        .bind(config_value)
        .bind(description)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
