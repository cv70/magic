#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Business Plan Record from database
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct BusinessPlan {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub industry: String,
    pub region: String,
    pub financing_amount: f64,
    pub company_size: String,
    pub created_at: Option<String>,
}

// API Request/Response structures
#[derive(Serialize, Deserialize, Debug)]
pub struct GetBusinessPlanReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetBusinessPlanRes {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub industry: String,
    pub region: String,
    pub financing_amount: f64,
    pub company_size: String,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetBusinessPlansReq {
    pub ids: Vec<i64>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct CreateBusinessPlanReq {
    pub title: String,
    pub content: String,
    pub industry: String,
    pub region: String,
    pub financing_amount: f64,
    pub company_size: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct UpdateBusinessPlanReq {
    pub id: i64,
    pub title: Option<String>,
    pub content: Option<String>,
    pub industry: Option<String>,
    pub region: Option<String>,
    pub financing_amount: Option<f64>,
    pub company_size: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DeleteBusinessPlanReq {
    pub id: i64,
}
