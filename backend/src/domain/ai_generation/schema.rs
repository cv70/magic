#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// AI Generation domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Generator {
    pub id: i64,
    pub name: String,
    pub provider: String,
    pub model: String,
    pub api_key: Option<String>,
    pub api_endpoint: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AIProvider {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub api_key: Option<String>,
    pub api_endpoint: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AIModel {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub provider: String,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PromptTemplate {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub template: String,
    pub input_variables: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GenerateTask {
    pub id: i64,
    pub generator_id: i64,
    pub input: String,
    pub output: Option<String>,
    pub status: Option<String>,
    pub error: Option<String>,
    pub created_at: Option<String>,
    pub started_at: Option<String>,
    pub completed_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetGeneratorReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetGeneratorRes {
    pub generator: Generator,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchGeneratorsReq {
    pub provider: Option<String>,
    pub model: Option<String>,
    pub enabled: Option<bool>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchGeneratorsRes {
    pub generators: Vec<Generator>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddGeneratorReq {
    pub name: String,
    pub provider: String,
    pub model: String,
    pub api_key: Option<String>,
    pub api_endpoint: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddGeneratorRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateGeneratorReq {
    pub id: i64,
    pub name: Option<String>,
    pub provider: Option<String>,
    pub model: Option<String>,
    pub api_key: Option<String>,
    pub api_endpoint: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateGeneratorRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GenerateContentReq {
    pub generator_id: i64,
    pub input: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GenerateContentRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPromptTemplateReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPromptTemplateRes {
    pub prompt_template: PromptTemplate,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPromptTemplatesReq {
    pub name: Option<String>,
    pub provider: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPromptTemplatesRes {
    pub prompt_templates: Vec<PromptTemplate>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPromptTemplateReq {
    pub name: String,
    pub description: Option<String>,
    pub template: String,
    pub input_variables: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPromptTemplateRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePromptTemplateReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
    pub template: Option<String>,
    pub input_variables: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePromptTemplateRes {
    pub id: i64,
}
