use super::schema::{AddNewsReq, AddNewsRes, GetNewsReq, SearchNewsReq, SearchNewsRes};
use crate::state::state::AppState;
use crate::utils::http::ApiResponse;
use axum::extract::{Json, State};
use axum::response::IntoResponse;

/// Get news by ID
pub async fn api_get_news(
    State(state): State<AppState>,
    Json(req): Json<GetNewsReq>,
) -> impl IntoResponse {
    match state.news_domain.get_news(req.id).await {
        Some(news) => ApiResponse {
            code: 200,
            message: None,
            data: Some(news),
        },
        None => ApiResponse {
            code: 404,
            message: Some("news not found".to_string()),
            data: None,
        },
    }
}

/// Search news by keywords
pub async fn api_search_news(
    State(state): State<AppState>,
    Json(req): Json<SearchNewsReq>,
) -> impl IntoResponse {
    let (news, total) = state
        .news_domain
        .search_news(
            req.query,
            req.category,
            req.region,
            req.industry,
            req.page,
            req.limit,
        )
        .await;
    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchNewsRes { news, total }),
    }
}

/// Add a new news item
pub async fn api_add_news(
    State(state): State<AppState>,
    Json(req): Json<AddNewsReq>,
) -> impl IntoResponse {
    match state
        .news_domain
        .add_news(req.title, req.content, req.category, req.region, req.source)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddNewsRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}
