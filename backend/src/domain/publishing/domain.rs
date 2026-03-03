#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::publishing::schema::{PublishLog, PublishTask, Publisher},
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct PublishingDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl PublishingDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_publisher(&self, id: i64) -> Option<Publisher> {
        self.db_dao
            .get_publisher(id)
            .await
            .ok()
            .map(|record| Publisher {
                id: record.id,
                name: record.name,
                platform: record.platform,
                platform_id: record.platform_id,
                enabled: Some(record.enabled),
                created_at: record.created_at,
                updated_at: None,
            })
    }

    pub async fn search_publishers(
        &self,
        platform: Option<String>,
        enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> (Vec<Publisher>, i64) {
        match self
            .db_dao
            .search_publishers(platform, enabled, page, limit)
            .await
        {
            Ok(records) => {
                let publishers: Vec<Publisher> = records
                    .into_iter()
                    .map(|record| Publisher {
                        id: record.id,
                        name: record.name,
                        platform: record.platform,
                        platform_id: record.platform_id,
                        enabled: Some(record.enabled),
                        created_at: record.created_at,
                        updated_at: None,
                    })
                    .collect();
                let total = publishers.len() as i64;
                (publishers, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_publisher(
        &self,
        name: String,
        platform: String,
        platform_id: Option<String>,
        enabled: bool,
    ) -> Result<i64, String> {
        self.db_dao
            .add_publisher(name.as_str(), platform.as_str(), platform_id, enabled)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_publisher(
        &self,
        id: i64,
        name: Option<String>,
        platform: Option<String>,
        platform_id: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_publisher(id, name, platform, platform_id, enabled)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn publish_content(&self, publisher_id: i64, content_id: i64) -> Result<i64, String> {
        self.db_dao
            .create_publish_task(publisher_id, content_id)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn get_publish_task(&self, id: i64) -> Option<PublishTask> {
        self.db_dao
            .get_publish_task(id)
            .await
            .ok()
            .map(|record| PublishTask {
                id: record.id,
                publisher_id: record.publisher_id,
                content_id: record.content_id,
                content: record.content,
                status: record.status,
                error: None,
                created_at: record.created_at,
                started_at: None,
                completed_at: None,
            })
    }

    pub async fn search_publish_tasks(
        &self,
        status: Option<String>,
        platform: Option<String>,
        created_at: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<PublishTask>, i64) {
        match self
            .db_dao
            .search_publish_tasks(status, platform, created_at, page, limit)
            .await
        {
            Ok(records) => {
                let tasks: Vec<PublishTask> = records
                    .into_iter()
                    .map(|record| PublishTask {
                        id: record.id,
                        publisher_id: record.publisher_id,
                        content_id: record.content_id,
                        content: record.content,
                        status: record.status,
                        error: None,
                        created_at: record.created_at,
                        started_at: None,
                        completed_at: None,
                    })
                    .collect();
                let total = tasks.len() as i64;
                (tasks, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn get_publish_log(&self, id: i64) -> Option<PublishLog> {
        self.db_dao
            .get_publish_log(id)
            .await
            .ok()
            .map(|record| PublishLog {
                id: record.id,
                publish_task_id: record.publish_task_id,
                log_type: record.log_type,
                message: record.message,
                created_at: record.created_at,
            })
    }

    pub async fn search_publish_logs(
        &self,
        publish_task_id: i64,
        log_type: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<PublishLog>, i64) {
        match self
            .db_dao
            .search_publish_logs(publish_task_id, log_type, page, limit)
            .await
        {
            Ok(records) => {
                let logs: Vec<PublishLog> = records
                    .into_iter()
                    .map(|record| PublishLog {
                        id: record.id,
                        publish_task_id: record.publish_task_id,
                        log_type: record.log_type,
                        message: record.message,
                        created_at: record.created_at,
                    })
                    .collect();
                let total = logs.len() as i64;
                (logs, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }
}
