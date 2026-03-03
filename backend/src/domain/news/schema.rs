#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// News structures
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct News {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub category: Option<String>,
    pub region: Option<String>,
    pub industry: Option<String>,
    pub publish_date: Option<String>,
    pub source: Option<String>,
    pub likes: i64,
    pub views: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetNewsReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetNewsRes {
    pub news: News,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SearchNewsReq {
    pub query: String,
    pub category: Option<String>,
    pub region: Option<String>,
    pub industry: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct SearchNewsRes {
    pub news: Vec<News>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct AddNewsReq {
    pub title: String,
    pub content: String,
    pub category: Option<String>,
    pub region: Option<String>,
    pub industry: Option<String>,
    pub source: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct AddNewsRes {
    pub id: i64,
}
