#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct GeneratorRecord {
    pub id: i64,
    pub name: String,
    pub provider: String,
    pub model: String,
    pub api_key: Option<String>,
    pub api_endpoint: Option<String>,
    pub enabled: bool,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct PromptTemplateRecord {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub template: String,
    pub input_variables: Option<String>,
    pub provider: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct GenerateTaskRecord {
    pub id: i64,
    pub generator_id: i64,
    pub input: String,
    pub output: Option<String>,
    pub status: Option<String>,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get generator by ID
    pub async fn get_generator(&self, id: i64) -> Result<GeneratorRecord> {
        let row = sqlx::query(
            "SELECT id, name, provider, model, api_key, api_endpoint, enabled, created_at
             FROM ai_generators WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let generator = GeneratorRecord {
            id: row.get("id"),
            name: row.get("name"),
            provider: row.get("provider"),
            model: row.get("model"),
            api_key: row.get("api_key"),
            api_endpoint: row.get("api_endpoint"),
            enabled: row.get("enabled"),
            created_at: row.get("created_at"),
        };

        Ok(generator)
    }

    /// Search generators with pagination
    pub async fn search_generators(
        &self,
        provider: Option<String>,
        _model: Option<String>,
        _enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<GeneratorRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, provider, model, api_key, api_endpoint, enabled, created_at
             FROM ai_generators
             WHERE ($1 IS NULL OR provider = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(provider)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let generators = rows
            .into_iter()
            .map(|row| GeneratorRecord {
                id: row.get("id"),
                name: row.get("name"),
                provider: row.get("provider"),
                model: row.get("model"),
                api_key: row.get("api_key"),
                api_endpoint: row.get("api_endpoint"),
                enabled: row.get("enabled"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(generators)
    }

    /// Add generator
    pub async fn add_generator(
        &self,
        name: &str,
        provider: &str,
        model: &str,
        api_key: Option<String>,
        api_endpoint: Option<String>,
        enabled: bool,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO ai_generators (name, provider, model, api_key, api_endpoint, enabled)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING id",
        )
        .bind(name)
        .bind(provider)
        .bind(model)
        .bind(api_key)
        .bind(api_endpoint)
        .bind(enabled)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update generator
    pub async fn update_generator(
        &self,
        id: i64,
        name: Option<String>,
        provider: Option<String>,
        model: Option<String>,
        api_key: Option<String>,
        api_endpoint: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE ai_generators
             SET name = COALESCE($1, name),
                 provider = COALESCE($2, provider),
                 model = COALESCE($3, model),
                 api_key = COALESCE($4, api_key),
                 api_endpoint = COALESCE($5, api_endpoint),
                 enabled = COALESCE($6, enabled)
             WHERE id = $7
             RETURNING id",
        )
        .bind(name)
        .bind(provider)
        .bind(model)
        .bind(api_key)
        .bind(api_endpoint)
        .bind(enabled)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Generate content
    pub async fn generate_content(&self, generator_id: i64, input: &str) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO ai_generation_tasks (generator_id, input)
             VALUES ($1, $2)
             RETURNING id",
        )
        .bind(generator_id)
        .bind(input)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Get prompt template by ID
    pub async fn get_prompt_template(&self, id: i64) -> Result<PromptTemplateRecord> {
        let row = sqlx::query(
            "SELECT id, name, description, template, input_variables, provider, created_at
             FROM ai_prompt_templates WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let template = PromptTemplateRecord {
            id: row.get("id"),
            name: row.get("name"),
            description: row.get("description"),
            template: row.get("template"),
            input_variables: row.get("input_variables"),
            provider: row.get("provider"),
            created_at: row.get("created_at"),
        };

        Ok(template)
    }

    /// Search prompt templates with pagination
    pub async fn search_prompt_templates(
        &self,
        name: Option<String>,
        _provider: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<PromptTemplateRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, description, template, input_variables, provider, created_at
             FROM ai_prompt_templates
             WHERE ($1 IS NULL OR name LIKE $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(name)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let templates = rows
            .into_iter()
            .map(|row| PromptTemplateRecord {
                id: row.get("id"),
                name: row.get("name"),
                description: row.get("description"),
                template: row.get("template"),
                input_variables: row.get("input_variables"),
                provider: row.get("provider"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(templates)
    }

    /// Add prompt template
    pub async fn add_prompt_template(
        &self,
        name: &str,
        description: Option<String>,
        template: &str,
        input_variables: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO ai_prompt_templates (name, description, template, input_variables)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .bind(template)
        .bind(input_variables)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update prompt template
    pub async fn update_prompt_template(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
        template: Option<String>,
        input_variables: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE ai_prompt_templates
             SET name = COALESCE($1, name),
                 description = COALESCE($2, description),
                 template = COALESCE($3, template),
                 input_variables = COALESCE($4, input_variables)
             WHERE id = $5
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .bind(template)
        .bind(input_variables)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
