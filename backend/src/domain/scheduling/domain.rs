#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::scheduling::task::Task,
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct SchedulingDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl SchedulingDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_task(&self, id: i64) -> Option<Task> {
        self.db_dao.get_task(id).await.ok().map(|record| Task {
            id: record.id,
            name: record.name,
            description: None,
            task_type: record.task_type,
            scheduler_id: record.scheduler_id,
            cron_expression: record.cron_expression,
            next_run_at: record.next_run_at,
            last_run_at: record.last_run_at,
            last_run_status: record.last_run_status,
            last_run_error: record.last_run_error,
            last_run_duration: record.last_run_duration,
            last_run_result: None,
            enabled: Some(record.enabled),
            created_at: record.created_at,
            updated_at: None,
        })
    }

    pub async fn search_tasks(
        &self,
        name: Option<String>,
        task_type: Option<String>,
        scheduler_id: Option<i64>,
        enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> (Vec<Task>, i64) {
        match self
            .db_dao
            .search_tasks(name, task_type, scheduler_id, enabled, page, limit)
            .await
        {
            Ok(records) => {
                let tasks: Vec<Task> = records
                    .into_iter()
                    .map(|record| Task {
                        id: record.id,
                        name: record.name,
                        description: None,
                        task_type: record.task_type,
                        scheduler_id: record.scheduler_id,
                        cron_expression: record.cron_expression,
                        next_run_at: record.next_run_at,
                        last_run_at: record.last_run_at,
                        last_run_status: record.last_run_status,
                        last_run_error: record.last_run_error,
                        last_run_duration: record.last_run_duration,
                        last_run_result: None,
                        enabled: Some(record.enabled),
                        created_at: record.created_at,
                        updated_at: None,
                    })
                    .collect();
                let total = tasks.len() as i64;
                (tasks, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_task(
        &self,
        name: String,
        task_type: String,
        scheduler_id: i64,
        cron_expression: Option<String>,
        enabled: bool,
    ) -> Result<i64, String> {
        self.db_dao
            .add_task(
                name.as_str(),
                task_type.as_str(),
                scheduler_id,
                cron_expression,
                enabled,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_task(
        &self,
        id: i64,
        name: Option<String>,
        task_type: Option<String>,
        scheduler_id: Option<i64>,
        cron_expression: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_task(id, name, task_type, scheduler_id, cron_expression, enabled)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn run_task(&self, id: i64) -> Result<i64, String> {
        self.db_dao.run_task(id).await.map_err(|e| e.to_string())
    }

    pub async fn stop_task(&self, id: i64) -> Result<i64, String> {
        self.db_dao.stop_task(id).await.map_err(|e| e.to_string())
    }

    pub async fn restart_task(&self, id: i64) -> Result<i64, String> {
        self.db_dao
            .restart_task(id)
            .await
            .map_err(|e| e.to_string())
    }
}
