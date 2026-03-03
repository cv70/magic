#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::financing::schema::BusinessPlan,
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct FinancingDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl FinancingDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_business_plan(&self, id: i64) -> Result<BusinessPlan, String> {
        self.db_dao
            .get_business_plan(id)
            .await
            .map(|record| BusinessPlan {
                id: record.id,
                title: record.title,
                content: record.content,
                industry: record.industry,
                region: record.region,
                financing_amount: record.financing_amount,
                company_size: record.company_size,
                created_at: record.created_at,
            })
            .map_err(|e| e.to_string())
    }

    pub async fn get_business_plans(&self, ids: Vec<i64>) -> Result<Vec<BusinessPlan>, String> {
        self.db_dao
            .get_business_plans(ids)
            .await
            .map(|records| {
                records
                    .into_iter()
                    .map(|record| BusinessPlan {
                        id: record.id,
                        title: record.title,
                        content: record.content,
                        industry: record.industry,
                        region: record.region,
                        financing_amount: record.financing_amount,
                        company_size: record.company_size,
                        created_at: record.created_at,
                    })
                    .collect()
            })
            .map_err(|e| e.to_string())
    }

    pub async fn create_business_plan(
        &self,
        title: &str,
        content: &str,
        industry: &str,
        region: &str,
        financing_amount: f64,
        company_size: &str,
    ) -> Result<i64, String> {
        self.db_dao
            .create_business_plan(
                title,
                content,
                industry,
                region,
                financing_amount,
                company_size,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_business_plan(
        &self,
        id: i64,
        title: Option<String>,
        content: Option<String>,
        industry: Option<String>,
        region: Option<String>,
        financing_amount: Option<f64>,
        company_size: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_business_plan(
                id,
                title,
                content,
                industry,
                region,
                financing_amount,
                company_size,
            )
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn delete_business_plan(&self, id: i64) -> Result<i64, String> {
        self.db_dao
            .delete_business_plan(id)
            .await
            .map_err(|e| e.to_string())
    }
}
