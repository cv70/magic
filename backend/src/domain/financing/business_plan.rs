use anyhow::Result;
use crate::domain::financing::schema::BusinessPlan;
use super::domain::FinancingDomain;

impl FinancingDomain {
    /// Get business plan by ID
    pub async fn get_business_plan(&self, id: i64) -> Result<BusinessPlan> {
        let db_record = self.db_dao.get_business_plan(id).await?;
        Ok(BusinessPlan {
            id: db_record.id,
            title: db_record.title,
            content: db_record.content,
            industry: db_record.industry,
            region: db_record.region,
            financing_amount: db_record.financing_amount,
            company_size: db_record.company_size,
            created_at: db_record.created_at,
        })
    }

    /// Get multiple business plans by IDs
    pub async fn get_business_plans(&self, ids: Vec<i64>) -> Result<Vec<BusinessPlan>> {
        let db_records = self.db_dao.get_business_plans(ids).await?;
        Ok(db_records.into_iter().map(|r| BusinessPlan {
            id: r.id,
            title: r.title,
            content: r.content,
            industry: r.industry,
            region: r.region,
            financing_amount: r.financing_amount,
            company_size: r.company_size,
            created_at: r.created_at,
        }).collect())
    }

    /// Create a new business plan
    pub async fn create_business_plan(
        &self,
        title: &str,
        content: &str,
        industry: &str,
        region: &str,
        financing_amount: f64,
        company_size: &str,
    ) -> Result<i64> {
        self.db_dao.create_business_plan(title, content, industry, region, financing_amount, company_size).await
    }

    /// Update business plan
    pub async fn update_business_plan(
        &self,
        id: i64,
        title: Option<String>,
        content: Option<String>,
        industry: Option<String>,
        region: Option<String>,
        financing_amount: Option<f64>,
        company_size: Option<String>,
    ) -> Result<i64> {
        self.db_dao.update_business_plan(id, title, content, industry, region, financing_amount, company_size).await
    }

    /// Delete business plan
    pub async fn delete_business_plan(&self, id: i64) -> Result<i64> {
        self.db_dao.delete_business_plan(id).await
    }
}
