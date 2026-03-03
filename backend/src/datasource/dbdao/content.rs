#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct ContentRecord {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub content_type: Option<String>,
    pub status: Option<String>,
    pub tags: Option<String>,
    pub created_at: Option<String>,
    pub published_at: Option<String>,
    pub updated_at: Option<String>,
}

impl DBDao {
    /// Get content by ID
    pub async fn get_content(&self, id: i64) -> Result<ContentRecord> {
        let row = sqlx::query(
            "SELECT id, title, content, content_type, status, tags, created_at, published_at, updated_at
             FROM content WHERE id = $1")
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let content = ContentRecord {
            id: row.get("id"),
            title: row.get("title"),
            content: row.get("content"),
            content_type: row.get("content_type"),
            status: row.get("status"),
            tags: row.get("tags"),
            created_at: row.get("created_at"),
            published_at: row.get("published_at"),
            updated_at: row.get("updated_at"),
        };

        Ok(content)
    }

    /// Search content with pagination
    pub async fn search_contents(
        &self,
        query: Option<&str>,
        _content_type: Option<String>,
        _status: Option<String>,
        _tag: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<ContentRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, title, content, content_type, status, tags, created_at, published_at, updated_at
             FROM content
             WHERE $1 IS NULL OR content LIKE $1
             OFFSET $2 LIMIT $3")
        .bind(query.unwrap_or("%%"))
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let contents = rows
            .into_iter()
            .map(|row| ContentRecord {
                id: row.get("id"),
                title: row.get("title"),
                content: row.get("content"),
                content_type: row.get("content_type"),
                status: row.get("status"),
                tags: row.get("tags"),
                created_at: row.get("created_at"),
                published_at: row.get("published_at"),
                updated_at: row.get("updated_at"),
            })
            .collect::<Vec<_>>();

        Ok(contents)
    }

    /// Add content
    pub async fn add_content(
        &self,
        title: &str,
        content: &str,
        content_type: Option<String>,
        tags: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO content (title, content, content_type, tags)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(title)
        .bind(content)
        .bind(content_type)
        .bind(tags)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update content
    pub async fn update_content(
        &self,
        id: i64,
        title: Option<String>,
        content: Option<String>,
        content_type: Option<String>,
        tags: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE content
             SET title = COALESCE($1, title),
                 content = COALESCE($2, content),
                 content_type = COALESCE($3, content_type),
                 tags = COALESCE($4, tags),
                 updated_at = NOW()
             WHERE id = $5
             RETURNING id",
        )
        .bind(title)
        .bind(content)
        .bind(content_type)
        .bind(tags)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Delete content
    pub async fn delete_content(&self, id: i64) -> Result<i64> {
        let row = sqlx::query("DELETE FROM content WHERE id = $1 RETURNING id")
            .bind(id)
            .fetch_one(&self.db)
            .await?;

        Ok(row.get("id"))
    }

    /// Create content version
    pub async fn create_content_version(
        &self,
        content_id: i64,
        content: &str,
        created_by: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO content_versions (content_id, content, created_by)
             VALUES ($1, $2, $3)
             RETURNING id",
        )
        .bind(content_id)
        .bind(content)
        .bind(created_by)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Publish content
    pub async fn publish_content(
        &self,
        id: i64,
        published_at: Option<String>,
        published_by: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE content
             SET status = 'published',
                 published_at = COALESCE($1, published_at),
                 published_by = COALESCE($2, published_by)
             WHERE id = $3
             RETURNING id",
        )
        .bind(published_at)
        .bind(published_by)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
