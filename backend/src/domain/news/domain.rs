#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::news::schema::News,
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct NewsDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl NewsDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_news(&self, id: i64) -> Option<News> {
        self.db_dao.get_news(id).await.ok().map(|record| News {
            id: record.id,
            title: record.title,
            content: record.content,
            category: record.category,
            region: record.region,
            industry: None,
            publish_date: record.publish_date,
            source: record.source,
            likes: record.likes,
            views: record.views,
        })
    }

    pub async fn search_news(
        &self,
        query: String,
        category: Option<String>,
        region: Option<String>,
        industry: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<News>, i64) {
        match self
            .db_dao
            .search_news(query.as_str(), category, region, industry, page, limit)
            .await
        {
            Ok(records) => {
                let news: Vec<News> = records
                    .into_iter()
                    .map(|record| News {
                        id: record.id,
                        title: record.title,
                        content: record.content,
                        category: record.category,
                        region: record.region,
                        industry: None,
                        publish_date: record.publish_date,
                        source: record.source,
                        likes: record.likes,
                        views: record.views,
                    })
                    .collect();
                let total = news.len() as i64;
                (news, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_news(
        &self,
        title: String,
        content: String,
        category: Option<String>,
        region: Option<String>,
        source: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_news(title.as_str(), content.as_str(), category, region, source)
            .await
            .map_err(|e| e.to_string())
    }
}
