use serde::{Deserialize, Serialize};

// Founder structures
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Founder {
    pub id: i64,
    pub name: String,
    pub email: String,
    pub phone: Option<String>,
    pub industry: Vec<String>,
    pub region: Option<String>,
    pub company_size: Option<String>,
    pub experience: Option<String>,
    pub skills: Option<String>,
    pub investment_stage: Option<String>,
    pub investment_amount_lower: Option<i64>,
    pub investment_amount_upper: Option<i64>,
    pub introduction: Option<String>,
    pub embedding: Option<Vec<f32>>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Investor {
    pub id: i64,
    pub name: String,
    pub company: String,
    pub industry: Vec<String>,
    pub region: Option<String>,
    pub investment_stage: Option<String>,
    pub investment_amount_lower: Option<i64>,
    pub investment_amount_upper: Option<i64>,
    pub introduction: Option<String>,
    pub embedding: Option<Vec<f32>>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct MatchedInvestor {
    pub id: i64,
    pub name: String,
    pub company: String,
    pub industry: Vec<String>,
    pub region: Option<String>,
    pub investment_stage: Option<String>,
    pub investment_amount_lower: Option<i64>,
    pub investment_amount_upper: Option<i64>,
    pub match_reason: String,
    pub score: f32,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct MatchedFounder {
    pub id: i64,
    pub name: String,
    pub company: String,
    pub industry: Vec<String>,
    pub region: Option<String>,
    pub investment_stage: Option<String>,
    pub investment_amount_lower: Option<i64>,
    pub investment_amount_upper: Option<i64>,
    pub match_reason: String,
    pub score: f32,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct MatchInvestorsReq {
    pub founder_id: i64,
    pub top_n: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct MatchFoundersReq {
    pub investor_id: i64,
    pub top_n: i64,
}
