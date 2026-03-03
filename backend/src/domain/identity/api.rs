use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::schema::{
    AddPermissionReq, AddPermissionRes, AddRoleReq, AddRoleRes, AddUserReq, AddUserRes,
    GetPermissionReq, GetRoleReq, GetUserReq, SearchPermissionsReq, SearchPermissionsRes,
    SearchRolesReq, SearchRolesRes, SearchUsersReq, SearchUsersRes, UpdatePermissionReq,
    UpdatePermissionRes, UpdateRoleReq, UpdateRoleRes, UpdateUserReq, UpdateUserRes,
};

pub async fn api_get_user(
    State(state): State<AppState>,
    Json(req): Json<GetUserReq>,
) -> impl IntoResponse {
    match state.identity_domain.get_user(req.id).await {
        Some(user) => ApiResponse {
            code: 200,
            message: None,
            data: Some(user),
        },
        None => ApiResponse {
            code: 404,
            message: Some("user not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_users(
    State(state): State<AppState>,
    Json(req): Json<SearchUsersReq>,
) -> impl IntoResponse {
    let (users, total) = state
        .identity_domain
        .search_users(req.username, req.email, req.role, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchUsersRes { users, total }),
    }
}

pub async fn api_add_user(
    State(state): State<AppState>,
    Json(req): Json<AddUserReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .add_user(req.username, req.email, req.password, req.role)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddUserRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_user(
    State(state): State<AppState>,
    Json(req): Json<UpdateUserReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .update_user(req.id, req.username, req.email, req.password, req.role)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateUserRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_get_role(
    State(state): State<AppState>,
    Json(req): Json<GetRoleReq>,
) -> impl IntoResponse {
    match state.identity_domain.get_role(req.id).await {
        Some(role) => ApiResponse {
            code: 200,
            message: None,
            data: Some(role),
        },
        None => ApiResponse {
            code: 404,
            message: Some("role not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_roles(
    State(state): State<AppState>,
    Json(req): Json<SearchRolesReq>,
) -> impl IntoResponse {
    let (roles, total) = state
        .identity_domain
        .search_roles(req.name, req.description, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchRolesRes { roles, total }),
    }
}

pub async fn api_add_role(
    State(state): State<AppState>,
    Json(req): Json<AddRoleReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .add_role(req.name, req.description)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddRoleRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_role(
    State(state): State<AppState>,
    Json(req): Json<UpdateRoleReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .update_role(req.id, req.name, req.description)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateRoleRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_get_permission(
    State(state): State<AppState>,
    Json(req): Json<GetPermissionReq>,
) -> impl IntoResponse {
    match state.identity_domain.get_permission(req.id).await {
        Some(permission) => ApiResponse {
            code: 200,
            message: None,
            data: Some(permission),
        },
        None => ApiResponse {
            code: 404,
            message: Some("permission not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_permissions(
    State(state): State<AppState>,
    Json(req): Json<SearchPermissionsReq>,
) -> impl IntoResponse {
    let (permissions, total) = state
        .identity_domain
        .search_permissions(req.name, req.description, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchPermissionsRes { permissions, total }),
    }
}

pub async fn api_add_permission(
    State(state): State<AppState>,
    Json(req): Json<AddPermissionReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .add_permission(req.name, req.description)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddPermissionRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_permission(
    State(state): State<AppState>,
    Json(req): Json<UpdatePermissionReq>,
) -> impl IntoResponse {
    match state
        .identity_domain
        .update_permission(req.id, req.name, req.description)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdatePermissionRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}
