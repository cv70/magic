#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct NewsRecord {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub category: Option<String>,
    pub region: Option<String>,
    pub publish_date: Option<String>,
    pub source: Option<String>,
    pub likes: i64,
    pub views: i64,
}

impl DBDao {
    /// Get news by ID
    pub async fn get_news(&self, id: i64) -> Result<NewsRecord> {
        let row = sqlx::query(
            "SELECT id, title, content, category, region, publish_date, source, likes, views
             FROM news WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let news = NewsRecord {
            id: row.get("id"),
            title: row.get("title"),
            content: row.get("content"),
            category: row.get("category"),
            region: row.get("region"),
            publish_date: row.get("publish_date"),
            source: row.get("source"),
            likes: row.get("likes"),
            views: row.get("views"),
        };

        Ok(news)
    }

    /// Search news with pagination
    pub async fn search_news(
        &self,
        query: &str,
        category: Option<String>,
        region: Option<String>,
        _industry: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<NewsRecord>> {
        let offset = (page - 1) * limit;
        let like = format!("%{}%", query);
        let rows = sqlx::query(
            "SELECT id, title, content, category, region, publish_date, source, likes, views
             FROM news
             WHERE ($1 = '%%' OR title ILIKE $1 OR content ILIKE $1)
               AND ($2 IS NULL OR category = $2)
               AND ($3 IS NULL OR region = $3)
             OFFSET $4 LIMIT $5",
        )
        .bind(like)
        .bind(category)
        .bind(region)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let news = rows
            .into_iter()
            .map(|row| NewsRecord {
                id: row.get("id"),
                title: row.get("title"),
                content: row.get("content"),
                category: row.get("category"),
                region: row.get("region"),
                publish_date: row.get("publish_date"),
                source: row.get("source"),
                likes: row.get("likes"),
                views: row.get("views"),
            })
            .collect::<Vec<_>>();

        Ok(news)
    }

    /// Add news
    pub async fn add_news(
        &self,
        title: &str,
        content: &str,
        category: Option<String>,
        region: Option<String>,
        source: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO news (title, content, category, region, source)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id",
        )
        .bind(title)
        .bind(content)
        .bind(category)
        .bind(region)
        .bind(source)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
