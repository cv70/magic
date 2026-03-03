#![allow(dead_code)]
use super::dao::DBDao;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use sqlx::Row;

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct UserRecord {
    pub id: i64,
    pub username: String,
    pub email: String,
    pub password: Option<String>,
    pub role: Option<String>,
    pub enabled: bool,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct RoleRecord {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Deserialize, Serialize, Debug, Clone)]
pub struct PermissionRecord {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub created_at: Option<String>,
}

impl DBDao {
    /// Get user by ID
    pub async fn get_user(&self, id: i64) -> Result<UserRecord> {
        let row = sqlx::query(
            "SELECT id, username, email, password, role, enabled, created_at
             FROM users WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let user = UserRecord {
            id: row.get("id"),
            username: row.get("username"),
            email: row.get("email"),
            password: row.get("password"),
            role: row.get("role"),
            enabled: row.get("enabled"),
            created_at: row.get("created_at"),
        };

        Ok(user)
    }

    /// Search users with pagination
    pub async fn search_users(
        &self,
        username: Option<String>,
        _email: Option<String>,
        _role: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<UserRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, username, email, password, role, enabled, created_at
             FROM users
             WHERE ($1 IS NULL OR username = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(username)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let users = rows
            .into_iter()
            .map(|row| UserRecord {
                id: row.get("id"),
                username: row.get("username"),
                email: row.get("email"),
                password: row.get("password"),
                role: row.get("role"),
                enabled: row.get("enabled"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(users)
    }

    /// Add user
    pub async fn add_user(
        &self,
        username: &str,
        email: &str,
        password: Option<String>,
        role: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO users (username, email, password, role)
             VALUES ($1, $2, $3, $4)
             RETURNING id",
        )
        .bind(username)
        .bind(email)
        .bind(password)
        .bind(role)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update user
    pub async fn update_user(
        &self,
        id: i64,
        username: Option<String>,
        email: Option<String>,
        password: Option<String>,
        role: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE users
             SET username = COALESCE($1, username),
                 email = COALESCE($2, email),
                 password = COALESCE($3, password),
                 role = COALESCE($4, role)
             WHERE id = $5
             RETURNING id",
        )
        .bind(username)
        .bind(email)
        .bind(password)
        .bind(role)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Get role by ID
    pub async fn get_role(&self, id: i64) -> Result<RoleRecord> {
        let row = sqlx::query(
            "SELECT id, name, description, created_at
             FROM roles WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let role = RoleRecord {
            id: row.get("id"),
            name: row.get("name"),
            description: row.get("description"),
            created_at: row.get("created_at"),
        };

        Ok(role)
    }

    /// Search roles with pagination
    pub async fn search_roles(
        &self,
        name: Option<String>,
        _description: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<RoleRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, description, created_at
             FROM roles
             WHERE ($1 IS NULL OR name = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(name)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let roles = rows
            .into_iter()
            .map(|row| RoleRecord {
                id: row.get("id"),
                name: row.get("name"),
                description: row.get("description"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(roles)
    }

    /// Add role
    pub async fn add_role(&self, name: &str, description: Option<String>) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO roles (name, description)
             VALUES ($1, $2)
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update role
    pub async fn update_role(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE roles
             SET name = COALESCE($1, name),
                 description = COALESCE($2, description)
             WHERE id = $3
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Get permission by ID
    pub async fn get_permission(&self, id: i64) -> Result<PermissionRecord> {
        let row = sqlx::query(
            "SELECT id, name, description, created_at
             FROM permissions WHERE id = $1",
        )
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        let permission = PermissionRecord {
            id: row.get("id"),
            name: row.get("name"),
            description: row.get("description"),
            created_at: row.get("created_at"),
        };

        Ok(permission)
    }

    /// Search permissions with pagination
    pub async fn search_permissions(
        &self,
        name: Option<String>,
        _description: Option<String>,
        page: i64,
        limit: i64,
    ) -> Result<Vec<PermissionRecord>> {
        let offset = (page - 1) * limit;
        let rows = sqlx::query(
            "SELECT id, name, description, created_at
             FROM permissions
             WHERE ($1 IS NULL OR name = $1)
             OFFSET $2 LIMIT $3",
        )
        .bind(name)
        .bind(offset)
        .bind(limit)
        .fetch_all(&self.db)
        .await?;

        let permissions = rows
            .into_iter()
            .map(|row| PermissionRecord {
                id: row.get("id"),
                name: row.get("name"),
                description: row.get("description"),
                created_at: row.get("created_at"),
            })
            .collect::<Vec<_>>();

        Ok(permissions)
    }

    /// Add permission
    pub async fn add_permission(&self, name: &str, description: Option<String>) -> Result<i64> {
        let row = sqlx::query(
            "INSERT INTO permissions (name, description)
             VALUES ($1, $2)
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }

    /// Update permission
    pub async fn update_permission(
        &self,
        id: i64,
        name: Option<String>,
        description: Option<String>,
    ) -> Result<i64> {
        let row = sqlx::query(
            "UPDATE permissions
             SET name = COALESCE($1, name),
                 description = COALESCE($2, description)
             WHERE id = $3
             RETURNING id",
        )
        .bind(name)
        .bind(description)
        .bind(id)
        .fetch_one(&self.db)
        .await?;

        Ok(row.get("id"))
    }
}
