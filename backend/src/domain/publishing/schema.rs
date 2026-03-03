#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Publishing domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Publisher {
    pub id: i64,
    pub name: String,
    pub platform: String,
    pub platform_id: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishPlatform {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishTask {
    pub id: i64,
    pub publisher_id: i64,
    pub content_id: i64,
    pub content: String,
    pub status: Option<String>,
    pub error: Option<String>,
    pub created_at: Option<String>,
    pub started_at: Option<String>,
    pub completed_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishLog {
    pub id: i64,
    pub publish_task_id: i64,
    pub log_type: String,
    pub message: String,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublisherReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublisherRes {
    pub publisher: Publisher,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPublishersReq {
    pub platform: Option<String>,
    pub enabled: Option<bool>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPublishersRes {
    pub publishers: Vec<Publisher>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPublisherReq {
    pub name: String,
    pub platform: String,
    pub platform_id: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPublisherRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePublisherReq {
    pub id: i64,
    pub name: Option<String>,
    pub platform: Option<String>,
    pub platform_id: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePublisherRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishContentReq {
    pub publisher_id: i64,
    pub content_id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishContentRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublishTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublishTaskRes {
    pub publish_task: PublishTask,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPublishTasksReq {
    pub status: Option<String>,
    pub platform: Option<String>,
    pub created_at: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPublishTasksRes {
    pub publish_tasks: Vec<PublishTask>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublishLogReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPublishLogRes {
    pub publish_log: PublishLog,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPublishLogsReq {
    pub publish_task_id: i64,
    pub log_type: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Debug, Clone)]
pub struct SearchPublishLogsRes {
    pub publish_logs: Vec<PublishLog>,
    pub total: i64,
}
