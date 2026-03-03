#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Content domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Content {
    pub id: i64,
    pub title: String,
    pub content: String,
    pub content_type: Option<String>,
    pub status: Option<String>,
    pub tags: Option<Vec<String>>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
    pub published_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ContentVersion {
    pub id: i64,
    pub content_id: i64,
    pub version: i64,
    pub content: String,
    pub created_by: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ContentType {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub fields: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ContentStatus {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Tag {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetContentReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetContentRes {
    pub content: Content,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchContentReq {
    pub query: Option<String>,
    pub content_type: Option<String>,
    pub status: Option<String>,
    pub tag: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchContentRes {
    pub contents: Vec<Content>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddContentReq {
    pub title: String,
    pub content: String,
    pub content_type: Option<String>,
    pub tags: Option<Vec<String>>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddContentRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateContentReq {
    pub id: i64,
    pub title: Option<String>,
    pub content: Option<String>,
    pub content_type: Option<String>,
    pub tags: Option<Vec<String>>,
    pub status: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateContentRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DeleteContentReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct DeleteContentRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CreateContentVersionReq {
    pub content_id: i64,
    pub content: String,
    pub created_by: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CreateContentVersionRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishContentReq {
    pub id: i64,
    pub published_at: Option<String>,
    pub published_by: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct PublishContentRes {
    pub id: i64,
}
