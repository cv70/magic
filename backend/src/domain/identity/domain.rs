#![allow(dead_code)]
use std::sync::Arc;

use crate::{
    datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao},
    domain::identity::schema::{Permission, Role, User},
    infra::registry::Registry,
};

#[derive(Clone)]
pub struct IdentityDomain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl IdentityDomain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }

    pub async fn get_user(&self, id: i64) -> Option<User> {
        self.db_dao.get_user(id).await.ok().map(|record| User {
            id: record.id,
            username: record.username,
            email: record.email,
            password: record.password,
            role: record.role,
            enabled: Some(record.enabled),
            created_at: record.created_at,
        })
    }

    pub async fn search_users(
        &self,
        username: Option<String>,
        email: Option<String>,
        role: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<User>, i64) {
        match self
            .db_dao
            .search_users(username, email, role, page, limit)
            .await
        {
            Ok(records) => {
                let users: Vec<User> = records
                    .into_iter()
                    .map(|record| User {
                        id: record.id,
                        username: record.username,
                        email: record.email,
                        password: record.password,
                        role: record.role,
                        enabled: Some(record.enabled),
                        created_at: record.created_at,
                    })
                    .collect();
                let total = users.len() as i64;
                (users, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_user(
        &self,
        username: String,
        email: String,
        password: Option<String>,
        role: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_user(username.as_str(), email.as_str(), password, role)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_user(
        &self,
        id: i64,
        username: Option<String>,
        email: Option<String>,
        password: Option<String>,
        role: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_user(id, username, email, password, role)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn get_role(&self, id: i64) -> Option<Role> {
        self.db_dao.get_role(id).await.ok().map(|record| Role {
            id: record.id,
            name: record.name,
            description: record.description,
            created_at: record.created_at,
        })
    }

    pub async fn search_roles(
        &self,
        name: Option<String>,
        description: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<Role>, i64) {
        match self
            .db_dao
            .search_roles(name, description, page, limit)
            .await
        {
            Ok(records) => {
                let roles: Vec<Role> = records
                    .into_iter()
                    .map(|record| Role {
                        id: record.id,
                        name: record.name,
                        description: record.description,
                        created_at: record.created_at,
                    })
                    .collect();
                let total = roles.len() as i64;
                (roles, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_role(&self, name: String, description: Option<String>) -> Result<i64, String> {
        self.db_dao
            .add_role(name.as_str(), description)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_role(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_role(id, name, description)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn get_permission(&self, id: i64) -> Option<Permission> {
        self.db_dao
            .get_permission(id)
            .await
            .ok()
            .map(|record| Permission {
                id: record.id,
                name: record.name,
                description: record.description,
                created_at: record.created_at,
            })
    }

    pub async fn search_permissions(
        &self,
        name: Option<String>,
        description: Option<String>,
        page: i64,
        limit: i64,
    ) -> (Vec<Permission>, i64) {
        match self
            .db_dao
            .search_permissions(name, description, page, limit)
            .await
        {
            Ok(records) => {
                let permissions: Vec<Permission> = records
                    .into_iter()
                    .map(|record| Permission {
                        id: record.id,
                        name: record.name,
                        description: record.description,
                        created_at: record.created_at,
                    })
                    .collect();
                let total = permissions.len() as i64;
                (permissions, total)
            }
            Err(_) => (Vec::new(), 0),
        }
    }

    pub async fn add_permission(
        &self,
        name: String,
        description: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .add_permission(name.as_str(), description)
            .await
            .map_err(|e| e.to_string())
    }

    pub async fn update_permission(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
    ) -> Result<i64, String> {
        self.db_dao
            .update_permission(id, name, description)
            .await
            .map_err(|e| e.to_string())
    }
}
