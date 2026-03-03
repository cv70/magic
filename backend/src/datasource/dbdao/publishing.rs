#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct PublisherRecord {
    pub id: i64,
    pub name: String,
    pub platform: String,
    pub platform_id: Option<String>,
    pub enabled: bool,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct PublishTaskRecord {
    pub id: i64,
    pub publisher_id: i64,
    pub content_id: i64,
    pub content: String,
    pub status: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct PublishLogRecord {
    pub id: i64,
    pub publish_task_id: i64,
    pub log_type: String,
    pub message: String,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get publisher by ID
    pub async fn get_publisher(&self, id: i64) -> Result<PublisherRecord> {
        let row = sqlx::query(
            "SELECT id, name, platform, platform_id, enabled, created_at
             FROM publishers WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let publisher = PublisherRecord {
            id: row.get("id"),
            name: row.get("name"),
            platform: row.get("platform"),
            platform_id: row.get("platform_id"),
            enabled: row.get("enabled"),
            created_at: row.get("created_at"),
        };

        Ok(publisher)
    }

    /// Search publishers with pagination
    pub async fn search_publishers(
        &self,
        platform: Option<String>,
        _enabled: Option<bool>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<PublisherRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, platform, platform_id, enabled, created_at
             FROM publishers
             WHERE ($1 IS NULL OR platform = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(platform)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let publishers = rows
            .into_iter()
            .map(|row| PublisherRecord {
                id: row.get("id"),
                name: row.get("name"),
                platform: row.get("platform"),
                platform_id: row.get("platform_id"),
                enabled: row.get("enabled"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(publishers)
    }

    /// Add publisher
    pub async fn add_publisher(
        &self,
        name: &str,
        platform: &str,
        platform_id: Option<String>,
        enabled: bool,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO publishers (name, platform, platform_id, enabled)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(name)
        .bind(platform)
        .bind(platform_id)
        .bind(enabled)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update publisher
    pub async fn update_publisher(
        &self,
        id: i64,
        name: Option<String>,
        platform: Option<String>,
        platform_id: Option<String>,
        enabled: Option<bool>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE publishers
             SET name = COALESCE($1, name),
                 platform = COALESCE($2, platform),
                 platform_id = COALESCE($3, platform_id),
                 enabled = COALESCE($4, enabled)
             WHERE id = $5
             RETURNING id",
        )
        .bind(name)
        .bind(platform)
        .bind(platform_id)
        .bind(enabled)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Create publish task for content
    pub async fn create_publish_task(&self, publisher_id: i64, content_id: i64) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO publish_tasks (publisher_id, content_id)
             VALUES ($1, $2)
             RETURNING id",
        )
        .bind(publisher_id)
        .bind(content_id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Get publish task by ID
    pub async fn get_publish_task(&self, id: i64) -> Result<PublishTaskRecord> {
        let row = sqlx::query(
            "SELECT id, publisher_id, content_id, content, status, created_at
             FROM publish_tasks WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let task = PublishTaskRecord {
            id: row.get("id"),
            publisher_id: row.get("publisher_id"),
            content_id: row.get("content_id"),
            content: row.get("content"),
            status: row.get("status"),
            created_at: row.get("created_at"),
        };

        Ok(task)
    }

    /// Search publish tasks with pagination
    pub async fn search_publish_tasks(
        &self,
        status: Option<String>,
        _platform: Option<String>,
        _created_at: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<PublishTaskRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, publisher_id, content_id, content, status, created_at
             FROM publish_tasks
         WHERE ($1 IS NULL OR status = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(status)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let tasks = rows
            .into_iter()
            .map(|row| PublishTaskRecord {
                id: row.get("id"),
                publisher_id: row.get("publisher_id"),
                content_id: row.get("content_id"),
                content: row.get("content"),
                status: row.get("status"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(tasks)
    }

    /// Get publish log by ID
    pub async fn get_publish_log(&self, id: i64) -> Result<PublishLogRecord> {
        let row = sqlx::query(
            "SELECT id, publish_task_id, log_type, message, created_at
             FROM publish_logs WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let log = PublishLogRecord {
            id: row.get("id"),
            publish_task_id: row.get("publish_task_id"),
            log_type: row.get("log_type"),
            message: row.get("message"),
            created_at: row.get("created_at"),
        };

        Ok(log)
    }

    /// Search publish logs with pagination
    pub async fn search_publish_logs(
        &self,
        publish_task_id: i64,
        _log_type: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<PublishLogRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, publish_task_id, log_type, message, created_at
             FROM publish_logs
             WHERE publish_task_id = $1
             OFFSET $2 LIMIT $3",
        )
        .bind(publish_task_id)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let logs = rows
            .into_iter()
            .map(|row| PublishLogRecord {
                id: row.get("id"),
                publish_task_id: row.get("publish_task_id"),
                log_type: row.get("log_type"),
                message: row.get("message"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(logs)
    }
}
