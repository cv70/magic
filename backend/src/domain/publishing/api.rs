use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::schema::{
    AddPublisherReq, AddPublisherRes, GetPublishLogReq, GetPublishTaskReq, GetPublisherReq,
    PublishContentReq, PublishContentRes, SearchPublishLogsReq, SearchPublishLogsRes,
    SearchPublishTasksReq, SearchPublishTasksRes, SearchPublishersReq, SearchPublishersRes,
    UpdatePublisherReq, UpdatePublisherRes,
};

pub async fn api_get_publisher(
    State(state): State<AppState>,
    Json(req): Json<GetPublisherReq>,
) -> impl IntoResponse {
    match state.publishing_domain.get_publisher(req.id).await {
        Some(publisher) => ApiResponse {
            code: 200,
            message: None,
            data: Some(publisher),
        },
        None => ApiResponse {
            code: 404,
            message: Some("publisher not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_publishers(
    State(state): State<AppState>,
    Json(req): Json<SearchPublishersReq>,
) -> impl IntoResponse {
    let (publishers, total) = state
        .publishing_domain
        .search_publishers(req.platform, req.enabled, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchPublishersRes { publishers, total }),
    }
}

pub async fn api_add_publisher(
    State(state): State<AppState>,
    Json(req): Json<AddPublisherReq>,
) -> impl IntoResponse {
    match state
        .publishing_domain
        .add_publisher(
            req.name,
            req.platform,
            req.platform_id,
            req.enabled.unwrap_or(true),
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddPublisherRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_publisher(
    State(state): State<AppState>,
    Json(req): Json<UpdatePublisherReq>,
) -> impl IntoResponse {
    match state
        .publishing_domain
        .update_publisher(req.id, req.name, req.platform, req.platform_id, req.enabled)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdatePublisherRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_publish_content(
    State(state): State<AppState>,
    Json(req): Json<PublishContentReq>,
) -> impl IntoResponse {
    match state
        .publishing_domain
        .publish_content(req.publisher_id, req.content_id)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(PublishContentRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_get_publish_task(
    State(state): State<AppState>,
    Json(req): Json<GetPublishTaskReq>,
) -> impl IntoResponse {
    match state.publishing_domain.get_publish_task(req.id).await {
        Some(publish_task) => ApiResponse {
            code: 200,
            message: None,
            data: Some(publish_task),
        },
        None => ApiResponse {
            code: 404,
            message: Some("publish task not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_publish_tasks(
    State(state): State<AppState>,
    Json(req): Json<SearchPublishTasksReq>,
) -> impl IntoResponse {
    let (publish_tasks, total) = state
        .publishing_domain
        .search_publish_tasks(
            req.status,
            req.platform,
            req.created_at,
            req.page,
            req.limit,
        )
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchPublishTasksRes {
            publish_tasks,
            total,
        }),
    }
}

pub async fn api_get_publish_log(
    State(state): State<AppState>,
    Json(req): Json<GetPublishLogReq>,
) -> impl IntoResponse {
    match state.publishing_domain.get_publish_log(req.id).await {
        Some(publish_log) => ApiResponse {
            code: 200,
            message: None,
            data: Some(publish_log),
        },
        None => ApiResponse {
            code: 404,
            message: Some("publish log not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_publish_logs(
    State(state): State<AppState>,
    Json(req): Json<SearchPublishLogsReq>,
) -> impl IntoResponse {
    let (publish_logs, total) = state
        .publishing_domain
        .search_publish_logs(req.publish_task_id, req.log_type, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchPublishLogsRes {
            publish_logs,
            total,
        }),
    }
}
