#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct TaskRecord {
    pub id: i64,
    pub name: String,
    pub task_type: String,
    pub scheduler_id: i64,
    pub cron_expression: Option<String>,
    pub next_run_at: Option<String>,
    pub last_run_at: Option<String>,
    pub last_run_status: Option<String>,
    pub last_run_error: Option<String>,
    pub last_run_duration: Option<i64>,
    pub enabled: bool,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get task by ID
    pub async fn get_task(&self, id: i64) -> Result<TaskRecord> {
        let row = sqlx::query(
            "SELECT id, name, task_type, scheduler_id, cron_expression, next_run_at, last_run_at, last_run_status, last_run_error, last_run_duration, enabled, created_at
             FROM tasks WHERE id = $1")
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let task = TaskRecord {
            id: row.get("id"),
            name: row.get("name"),
            task_type: row.get("task_type"),
            scheduler_id: row.get("scheduler_id"),
            cron_expression: row.get("cron_expression"),
            next_run_at: row.get("next_run_at"),
            last_run_at: row.get("last_run_at"),
            last_run_status: row.get("last_run_status"),
            last_run_error: row.get("last_run_error"),
            last_run_duration: row.get("last_run_duration"),
            enabled: row.get("enabled"),
            created_at: row.get("created_at"),
        };

        Ok(task)
    }

    /// Search tasks with pagination
    pub async fn search_tasks(
        &self,
        name: Option<String>,
        _task_type: Option<String>,
        _scheduler_id: Option<i64>,
        _enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<TaskRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, task_type, scheduler_id, cron_expression, next_run_at, last_run_at, last_run_status, last_run_error, last_run_duration, enabled, created_at
             FROM tasks
             WHERE ($1 IS NULL OR name = $1)
             OFFSET $2 LIMIT $3")
        .bind(name)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let tasks = rows
            .into_iter()
            .map(|row| TaskRecord {
                id: row.get("id"),
                name: row.get("name"),
                task_type: row.get("task_type"),
                scheduler_id: row.get("scheduler_id"),
                cron_expression: row.get("cron_expression"),
                next_run_at: row.get("next_run_at"),
                last_run_at: row.get("last_run_at"),
                last_run_status: row.get("last_run_status"),
                last_run_error: row.get("last_run_error"),
                last_run_duration: row.get("last_run_duration"),
                enabled: row.get("enabled"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(tasks)
    }

    /// Add task
    pub async fn add_task(
        &self,
        name: &str,
        task_type: &str,
        scheduler_id: i64,
        cron_expression: Option<String>,
        enabled: bool,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO tasks (name, task_type, scheduler_id, cron_expression, enabled)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id",
        )
        .bind(name)
        .bind(task_type)
        .bind(scheduler_id)
        .bind(cron_expression)
        .bind(enabled)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update task
    pub async fn update_task(
        &self,
        id: i64,
        name: Option<String>,
        task_type: Option<String>,
        scheduler_id: Option<i64>,
        cron_expression: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE tasks
             SET name = COALESCE($1, name),
                 task_type = COALESCE($2, task_type),
                 scheduler_id = COALESCE($3, scheduler_id),
                 cron_expression = COALESCE($4, cron_expression),
                 enabled = COALESCE($5, enabled)
             WHERE id = $6
             RETURNING id",
        )
        .bind(name)
        .bind(task_type)
        .bind(scheduler_id)
        .bind(cron_expression)
        .bind(enabled)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Run task
    pub async fn run_task(&self, id: i64) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE tasks
             SET last_run_status = 'running'
             WHERE id = $1
             RETURNING id",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Stop task
    pub async fn stop_task(&self, id: i64) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE tasks
             SET last_run_status = 'stopped'
             WHERE id = $1
             RETURNING id",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Restart task
    pub async fn restart_task(&self, id: i64) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE tasks
             SET last_run_status = 'restarted'
             WHERE id = $1
             RETURNING id",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
