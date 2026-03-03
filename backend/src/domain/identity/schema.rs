#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Identity domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct User {
    pub id: i64,
    pub username: String,
    pub email: String,
    pub password: Option<String>,
    pub role: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Role {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Permission {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub created_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetUserReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetUserRes {
    pub user: User,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchUsersReq {
    pub username: Option<String>,
    pub email: Option<String>,
    pub role: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchUsersRes {
    pub users: Vec<User>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddUserReq {
    pub username: String,
    pub email: String,
    pub password: Option<String>,
    pub role: Option<String>,
}

#[derive(Serialize, Debug, Clone)]
pub struct AddUserRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateUserReq {
    pub id: i64,
    pub username: Option<String>,
    pub email: Option<String>,
    pub password: Option<String>,
    pub role: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateUserRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetRoleReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetRoleRes {
    pub role: Role,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchRolesReq {
    pub name: Option<String>,
    pub description: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchRolesRes {
    pub roles: Vec<Role>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddRoleReq {
    pub name: String,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddRoleRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateRoleReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateRoleRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPermissionReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetPermissionRes {
    pub permission: Permission,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPermissionsReq {
    pub name: Option<String>,
    pub description: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchPermissionsRes {
    pub permissions: Vec<Permission>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPermissionReq {
    pub name: String,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddPermissionRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePermissionReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdatePermissionRes {
    pub id: i64,
}
