#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::ai_generation::schema::{Generator, PromptTemplate},
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct AIGenerationDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl AIGenerationDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_generator(&self, id: i64) -> Option<Generator> {
        self.db_dao
            .get_generator(id)
            .await
            .ok()
            .map(|record| Generator {
                id: record.id,
                name: record.name,
                provider: record.provider,
                model: record.model,
                api_key: record.api_key,
                api_endpoint: record.api_endpoint,
                enabled: Some(record.enabled),
                created_at: record.created_at,
            })
    }

    pub async fn search_generators(
        &self,
        provider: Option<String>,
        model: Option<String>,
        enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> (Vec<Generator>, i64) {
        match self
            .db_dao
            .search_generators(provider, model, enabled, page, limit)
            .await
        {
            Ok(records) => {
                let generators: Vec<Generator> = records
                    .into_iter()
                    .map(|record| Generator {
                        id: record.id,
                        name: record.name,
                        provider: record.provider,
                        model: record.model,
                        api_key: record.api_key,
                        api_endpoint: record.api_endpoint,
                        enabled: Some(record.enabled),
                        created_at: record.created_at,
                    })
                    .collect();
                let total = generators.len() as i64;
                (generators, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_generator(
        &self,
        name: String,
        provider: String,
        model: String,
        api_key: Option<String>,
        api_endpoint: Option<String>,
        enabled: bool,
    ) -> Result<i64, String> {
        self.db_dao
            .add_generator(
                name.as_str(),
                provider.as_str(),
                model.as_str(),
                api_key,
                api_endpoint,
                enabled,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_generator(
        &self,
        id: i64,
        name: Option<String>,
        provider: Option<String>,
        model: Option<String>,
        api_key: Option<String>,
        api_endpoint: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_generator(id, name, provider, model, api_key, api_endpoint, enabled)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn generate_content(&self, generator_id: i64, input: String) -> Result<i64, String> {
        self.db_dao
            .generate_content(generator_id, input.as_str())
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn get_prompt_template(&self, id: i64) -> Option<PromptTemplate> {
        self.db_dao
            .get_prompt_template(id)
            .await
            .ok()
            .map(|record| PromptTemplate {
                id: record.id,
                name: record.name,
                description: record.description,
                template: record.template,
                input_variables: record.input_variables,
                created_at: record.created_at,
            })
    }

    pub async fn search_prompt_templates(
        &self,
        name: Option<String>,
        provider: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<PromptTemplate>, i64) {
        match self
            .db_dao
            .search_prompt_templates(name, provider, page, limit)
            .await
        {
            Ok(records) => {
                let templates: Vec<PromptTemplate> = records
                    .into_iter()
                    .map(|record| PromptTemplate {
                        id: record.id,
                        name: record.name,
                        description: record.description,
                        template: record.template,
                        input_variables: record.input_variables,
                        created_at: record.created_at,
                    })
                    .collect();
                let total = templates.len() as i64;
                (templates, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_prompt_template(
        &self,
        name: String,
        description: Option<String>,
        template: String,
        input_variables: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_prompt_template(
                name.as_str(),
                description,
                template.as_str(),
                input_variables,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_prompt_template(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
        template: Option<String>,
        input_variables: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_prompt_template(id, name, description, template, input_variables)
            .await
            .map_err(|e| e.to_string())
    }
}
