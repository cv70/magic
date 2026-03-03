use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::schema::{
    AddContentReq, AddContentRes, DeleteContentReq, DeleteContentRes, GetContentReq,
    SearchContentReq, SearchContentRes, UpdateContentReq, UpdateContentRes,
};

pub async fn api_get_content(
    State(state): State<AppState>,
    Json(req): Json<GetContentReq>,
) -> impl IntoResponse {
    match state.content_domain.get_content(req.id).await {
        Some(content) => ApiResponse {
            code: 200,
            message: None,
            data: Some(content),
        },
        None => ApiResponse {
            code: 404,
            message: Some("content not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_contents(
    State(state): State<AppState>,
    Json(req): Json<SearchContentReq>,
) -> impl IntoResponse {
    let (contents, total) = state
        .content_domain
        .search_contents(
            req.query,
            req.content_type,
            req.status,
            req.tag,
            req.page,
            req.limit,
        )
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchContentRes { contents, total }),
    }
}

pub async fn api_add_content(
    State(state): State<AppState>,
    Json(req): Json<AddContentReq>,
) -> impl IntoResponse {
    match state
        .content_domain
        .add_content(req.title, req.content, req.content_type, req.tags)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddContentRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_content(
    State(state): State<AppState>,
    Json(req): Json<UpdateContentReq>,
) -> impl IntoResponse {
    match state
        .content_domain
        .update_content(
            req.id,
            req.title,
            req.content,
            req.content_type,
            req.tags,
            req.status,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateContentRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_delete_content(
    State(state): State<AppState>,
    Json(req): Json<DeleteContentReq>,
) -> impl IntoResponse {
    match state.content_domain.delete_content(req.id).await {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(DeleteContentRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}
