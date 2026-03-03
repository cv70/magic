use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct BusinessPlanRecord {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub industry: String,
    pub region: String,
    pub financing_amount: f64,
    pub company_size: String,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get business plan by ID
    pub async fn get_business_plan(&self, id: i64) -> Result<BusinessPlanRecord> {
        let row = sqlx::query(
            "SELECT id, title, content, industry, region, financing_amount, company_size, created_at
             FROM business_plans WHERE id = $1"
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let plan = BusinessPlanRecord {
            id: row.get("id"),
            title: row.get("title"),
            content: row.get("content"),
            industry: row.get("industry"),
            region: row.get("region"),
            financing_amount: row.get("financing_amount"),
            company_size: row.get("company_size"),
            created_at: row.get("created_at"),
        };

        Ok(plan)
    }

    /// Get multiple business plans by IDs
    pub async fn get_business_plans(&self, ids: Vec<i64>) -> Result<Vec<BusinessPlanRecord>> {
        let rows = sqlx::query(
            "SELECT id, title, content, industry, region, financing_amount, company_size, created_at
             FROM business_plans WHERE id = ANY($1)")
        .bind(ids)
        .fetch_all(&self.db)
        .await?;

        let plans = rows
            .into_iter()
            .map(|row| BusinessPlanRecord {
                id: row.get("id"),
                title: row.get("title"),
                content: row.get("content"),
                industry: row.get("industry"),
                region: row.get("region"),
                financing_amount: row.get("financing_amount"),
                company_size: row.get("company_size"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(plans)
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
        let row = sqlx::query(
            "INSERT INTO business_plans (title, content, industry, region, financing_amount, company_size)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING id")
        .bind(title)
        .bind(content)
        .bind(industry)
        .bind(region)
        .bind(financing_amount)
        .bind(company_size)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
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
        let row = sqlx::query(
            "UPDATE business_plans
             SET title = COALESCE($1, title),
                 content = COALESCE($2, content),
                 industry = COALESCE($3, industry),
                 region = COALESCE($4, region),
                 financing_amount = COALESCE($5, financing_amount),
                 company_size = COALESCE($6, company_size),
                 created_at = NOW()
             WHERE id = $7
             RETURNING id",
        )
        .bind(title)
        .bind(content)
        .bind(industry)
        .bind(region)
        .bind(financing_amount)
        .bind(company_size)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Delete business plan
    pub async fn delete_business_plan(&self, id: i64) -> Result<i64> {
        let row = sqlx::query("DELETE FROM business_plans WHERE id = $1 RETURNING id")
            .bind(id)
            .fetch_one(&self.db)
            .await?;

        Ok(row.get("id"))
    }
}
