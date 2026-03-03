use super::schema::CreateBusinessPlanReq;
use super::schema::DeleteBusinessPlanReq;
use super::schema::GetBusinessPlansReq;
use super::schema::UpdateBusinessPlanReq;
use super::schema::{GetBusinessPlanReq, GetBusinessPlanRes};
use crate::state::state::AppState;
use crate::utils::http::ApiResponse;
use axum::extract::{Json, State};
use axum::response::IntoResponse;

/// Get business plan by ID
pub async fn api_get_business_plan(
    State(state): State<AppState>,
    Json(req): Json<GetBusinessPlanReq>,
) -> impl IntoResponse {
    match state.financing_domain.get_business_plan(req.id).await {
        Ok(result) => {
            let plan = GetBusinessPlanRes {
                id: result.id,
                title: result.title,
                content: result.content,
                industry: result.industry,
                region: result.region,
                financing_amount: result.financing_amount,
                company_size: result.company_size,
                created_at: result.created_at,
            };
            ApiResponse {
                code: 200,
                message: None,
                data: Some(plan),
            }
        }
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}

/// Get multiple business plans by IDs
pub async fn api_get_business_plans(
    State(state): State<AppState>,
    Json(req): Json<GetBusinessPlansReq>,
) -> impl IntoResponse {
    match state.financing_domain.get_business_plans(req.ids).await {
        Ok(result) => {
            let plans: Vec<GetBusinessPlanRes> = result
                .into_iter()
                .map(|r| GetBusinessPlanRes {
                    id: r.id,
                    title: r.title,
                    content: r.content,
                    industry: r.industry,
                    region: r.region,
                    financing_amount: r.financing_amount,
                    company_size: r.company_size,
                    created_at: r.created_at,
                })
                .collect();
            ApiResponse {
                code: 200,
                message: None,
                data: Some(plans),
            }
        }
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}

/// Create a new business plan
pub async fn api_create_business_plan(
    State(state): State<AppState>,
    Json(req): Json<CreateBusinessPlanReq>,
) -> impl IntoResponse {
    match state
        .financing_domain
        .create_business_plan(
            &req.title,
            &req.content,
            &req.industry,
            &req.region,
            req.financing_amount,
            &req.company_size,
        )
        .await
    {
        Ok(result) => ApiResponse {
            code: 200,
            message: None,
            data: Some(result),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}

/// Update business plan
pub async fn api_update_business_plan(
    State(state): State<AppState>,
    Json(req): Json<UpdateBusinessPlanReq>,
) -> impl IntoResponse {
    match state
        .financing_domain
        .update_business_plan(
            req.id,
            req.title,
            req.content,
            req.industry,
            req.region,
            req.financing_amount,
            req.company_size,
        )
        .await
    {
        Ok(result) => ApiResponse {
            code: 200,
            message: None,
            data: Some(result),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}

/// Delete business plan
pub async fn api_delete_business_plan(
    State(state): State<AppState>,
    Json(req): Json<DeleteBusinessPlanReq>,
) -> impl IntoResponse {
    match state.financing_domain.delete_business_plan(req.id).await {
        Ok(result) => ApiResponse {
            code: 200,
            message: None,
            data: Some(result),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e.to_string()),
            data: None,
        },
    }
}
