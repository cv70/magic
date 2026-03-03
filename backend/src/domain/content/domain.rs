#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::content::schema::Content,
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct ContentDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl ContentDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_content(&self, id: i64) -> Option<Content> {
        self.db_dao
            .get_content(id)
            .await
            .ok()
            .map(|record| Content {
                id: record.id,
                title: record.title,
                content: record.content,
                content_type: record.content_type,
                status: record.status,
                tags: record
                    .tags
                    .map(|s| s.split(',').map(|t| t.trim().to_string()).collect()),
                created_at: record.created_at,
                updated_at: record.updated_at,
                published_at: record.published_at,
            })
    }

    pub async fn search_contents(
        &self,
        query: Option<String>,
        content_type: Option<String>,
        status: Option<String>,
        tag: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<Content>, i64) {
        let query = query.map(|q| format!("%{}%", q));
        match self
            .db_dao
            .search_contents(query.as_deref(), content_type, status, tag, page, limit)
            .await
        {
            Ok(records) => {
                let contents: Vec<Content> = records
                    .into_iter()
                    .map(|record| Content {
                        id: record.id,
                        title: record.title,
                        content: record.content,
                        content_type: record.content_type,
                        status: record.status,
                        tags: record
                            .tags
                            .map(|s| s.split(',').map(|t| t.trim().to_string()).collect()),
                        created_at: record.created_at,
                        updated_at: record.updated_at,
                        published_at: record.published_at,
                    })
                    .collect();
                let total = contents.len() as i64;
                (contents, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_content(
        &self,
        title: String,
        content: String,
        content_type: Option<String>,
        tags: Option<Vec<String>>,
    ) -> Result<i64, String> {
        let tags = tags.map(|v| v.join(","));
        self.db_dao
            .add_content(title.as_str(), content.as_str(), content_type, tags)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_content(
        &self,
        id: i64,
        title: Option<String>,
        content: Option<String>,
        content_type: Option<String>,
        tags: Option<Vec<String>>,
        _status: Option<String>,
    ) -> Result<i64, String> {
        let tags = tags.map(|v| v.join(","));
        self.db_dao
            .update_content(id, title, content, content_type, tags)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn delete_content(&self, id: i64) -> Result<i64, String> {
        self.db_dao
            .delete_content(id)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn create_content_version(
        &self,
        content_id: i64,
        content: String,
        created_by: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .create_content_version(content_id, content.as_str(), created_by)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn publish_content(
        &self,
        id: i64,
        published_at: Option<String>,
        published_by: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .publish_content(id, published_at, published_by)
            .await
            .map_err(|e| e.to_string())
    }
}
