use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::schema::{
    AddProviderConfigReq, AddProviderConfigRes, AddSystemConfigReq, AddSystemConfigRes,
    GetProviderConfigReq, GetSystemConfigReq, SearchProviderConfigsReq, SearchProviderConfigsRes,
    SearchSystemConfigsReq, SearchSystemConfigsRes, UpdateProviderConfigReq,
    UpdateProviderConfigRes, UpdateSystemConfigReq, UpdateSystemConfigRes,
};

pub async fn api_get_system_config(
    State(state): State<AppState>,
    Json(req): Json<GetSystemConfigReq>,
) -> impl IntoResponse {
    match state.configuration_domain.get_system_config(req.id).await {
        Some(system_config) => ApiResponse {
            code: 200,
            message: None,
            data: Some(system_config),
        },
        None => ApiResponse {
            code: 404,
            message: Some("system config not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_system_configs(
    State(state): State<AppState>,
    Json(req): Json<SearchSystemConfigsReq>,
) -> impl IntoResponse {
    let (system_configs, total) = state
        .configuration_domain
        .search_system_configs(req.key, req.category, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchSystemConfigsRes {
            system_configs,
            total,
        }),
    }
}

pub async fn api_add_system_config(
    State(state): State<AppState>,
    Json(req): Json<AddSystemConfigReq>,
) -> impl IntoResponse {
    match state
        .configuration_domain
        .add_system_config(req.key, req.value, req.description, req.category)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddSystemConfigRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_system_config(
    State(state): State<AppState>,
    Json(req): Json<UpdateSystemConfigReq>,
) -> impl IntoResponse {
    match state
        .configuration_domain
        .update_system_config(req.id, req.key, req.value, req.description, req.category)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateSystemConfigRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_get_provider_config(
    State(state): State<AppState>,
    Json(req): Json<GetProviderConfigReq>,
) -> impl IntoResponse {
    match state.configuration_domain.get_provider_config(req.id).await {
        Some(provider_config) => ApiResponse {
            code: 200,
            message: None,
            data: Some(provider_config),
        },
        None => ApiResponse {
            code: 404,
            message: Some("provider config not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_provider_configs(
    State(state): State<AppState>,
    Json(req): Json<SearchProviderConfigsReq>,
) -> impl IntoResponse {
    let (provider_configs, total) = state
        .configuration_domain
        .search_provider_configs(req.provider_name, req.config_key, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchProviderConfigsRes {
            provider_configs,
            total,
        }),
    }
}

pub async fn api_add_provider_config(
    State(state): State<AppState>,
    Json(req): Json<AddProviderConfigReq>,
) -> impl IntoResponse {
    match state
        .configuration_domain
        .add_provider_config(
            req.provider_name,
            req.config_key,
            req.config_value,
            req.description,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddProviderConfigRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_provider_config(
    State(state): State<AppState>,
    Json(req): Json<UpdateProviderConfigReq>,
) -> impl IntoResponse {
    match state
        .configuration_domain
        .update_provider_config(
            req.id,
            req.provider_name,
            req.config_key,
            req.config_value,
            req.description,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateProviderConfigRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}
