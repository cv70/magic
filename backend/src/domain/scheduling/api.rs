use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::task::{
    AddTaskReq, AddTaskRes, GetTaskReq, RestartTaskReq, RestartTaskRes, RunTaskReq, RunTaskRes,
    SearchTasksReq, SearchTasksRes, StopTaskReq, StopTaskRes, UpdateTaskReq, UpdateTaskRes,
};

pub async fn api_get_task(
    State(state): State<AppState>,
    Json(req): Json<GetTaskReq>,
) -> impl IntoResponse {
    match state.scheduling_domain.get_task(req.id).await {
        Some(task) => ApiResponse {
            code: 200,
            message: None,
            data: Some(task),
        },
        None => ApiResponse {
            code: 404,
            message: Some("task not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_tasks(
    State(state): State<AppState>,
    Json(req): Json<SearchTasksReq>,
) -> impl IntoResponse {
    let (tasks, total) = state
        .scheduling_domain
        .search_tasks(
            req.name,
            req.task_type,
            req.scheduler_id,
            req.enabled,
            req.page,
            req.limit,
        )
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchTasksRes { tasks, total }),
    }
}

pub async fn api_add_task(
    State(state): State<AppState>,
    Json(req): Json<AddTaskReq>,
) -> impl IntoResponse {
    match state
        .scheduling_domain
        .add_task(
            req.name,
            req.task_type,
            req.scheduler_id,
            req.cron_expression,
            req.enabled.unwrap_or(true),
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddTaskRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_task(
    State(state): State<AppState>,
    Json(req): Json<UpdateTaskReq>,
) -> impl IntoResponse {
    match state
        .scheduling_domain
        .update_task(
            req.id,
            req.name,
            req.task_type,
            req.scheduler_id,
            req.cron_expression,
            req.enabled,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateTaskRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_run_task(
    State(state): State<AppState>,
    Json(req): Json<RunTaskReq>,
) -> impl IntoResponse {
    match state.scheduling_domain.run_task(req.id).await {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(RunTaskRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_stop_task(
    State(state): State<AppState>,
    Json(req): Json<StopTaskReq>,
) -> impl IntoResponse {
    match state.scheduling_domain.stop_task(req.id).await {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(StopTaskRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_restart_task(
    State(state): State<AppState>,
    Json(req): Json<RestartTaskReq>,
) -> impl IntoResponse {
    match state.scheduling_domain.restart_task(req.id).await {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(RestartTaskRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}
