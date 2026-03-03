use axum::extract::{Json, State};
use axum::response::IntoResponse;

use crate::state::state::AppState;
use crate::utils::http::ApiResponse;

use super::schema::{
    AddGeneratorReq, AddGeneratorRes, AddPromptTemplateReq, AddPromptTemplateRes,
    GenerateContentReq, GenerateContentRes, GetGeneratorReq, GetPromptTemplateReq,
    SearchGeneratorsReq, SearchGeneratorsRes, SearchPromptTemplatesReq, SearchPromptTemplatesRes,
    UpdateGeneratorReq, UpdateGeneratorRes, UpdatePromptTemplateReq, UpdatePromptTemplateRes,
};

pub async fn api_get_generator(
    State(state): State<AppState>,
    Json(req): Json<GetGeneratorReq>,
) -> impl IntoResponse {
    match state.ai_generation_domain.get_generator(req.id).await {
        Some(generator) => ApiResponse {
            code: 200,
            message: None,
            data: Some(generator),
        },
        None => ApiResponse {
            code: 404,
            message: Some("generator not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_generators(
    State(state): State<AppState>,
    Json(req): Json<SearchGeneratorsReq>,
) -> impl IntoResponse {
    let (generators, total) = state
        .ai_generation_domain
        .search_generators(req.provider, req.model, req.enabled, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchGeneratorsRes { generators, total }),
    }
}

pub async fn api_add_generator(
    State(state): State<AppState>,
    Json(req): Json<AddGeneratorReq>,
) -> impl IntoResponse {
    match state
        .ai_generation_domain
        .add_generator(
            req.name,
            req.provider,
            req.model,
            req.api_key,
            req.api_endpoint,
            req.enabled.unwrap_or(true),
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddGeneratorRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_generator(
    State(state): State<AppState>,
    Json(req): Json<UpdateGeneratorReq>,
) -> impl IntoResponse {
    match state
        .ai_generation_domain
        .update_generator(
            req.id,
            req.name,
            req.provider,
            req.model,
            req.api_key,
            req.api_endpoint,
            req.enabled,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdateGeneratorRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_generate_content(
    State(state): State<AppState>,
    Json(req): Json<GenerateContentReq>,
) -> impl IntoResponse {
    match state
        .ai_generation_domain
        .generate_content(req.generator_id, req.input)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(GenerateContentRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_get_prompt_template(
    State(state): State<AppState>,
    Json(req): Json<GetPromptTemplateReq>,
) -> impl IntoResponse {
    match state.ai_generation_domain.get_prompt_template(req.id).await {
        Some(prompt_template) => ApiResponse {
            code: 200,
            message: None,
            data: Some(prompt_template),
        },
        None => ApiResponse {
            code: 404,
            message: Some("prompt template not found".to_string()),
            data: None,
        },
    }
}

pub async fn api_search_prompt_templates(
    State(state): State<AppState>,
    Json(req): Json<SearchPromptTemplatesReq>,
) -> impl IntoResponse {
    let (prompt_templates, total) = state
        .ai_generation_domain
        .search_prompt_templates(req.name, req.provider, req.page, req.limit)
        .await;

    ApiResponse {
        code: 200,
        message: None,
        data: Some(SearchPromptTemplatesRes {
            prompt_templates,
            total,
        }),
    }
}

pub async fn api_add_prompt_template(
    State(state): State<AppState>,
    Json(req): Json<AddPromptTemplateReq>,
) -> impl IntoResponse {
    match state
        .ai_generation_domain
        .add_prompt_template(req.name, req.description, req.template, req.input_variables)
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(AddPromptTemplateRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}

pub async fn api_update_prompt_template(
    State(state): State<AppState>,
    Json(req): Json<UpdatePromptTemplateReq>,
) -> impl IntoResponse {
    match state
        .ai_generation_domain
        .update_prompt_template(
            req.id,
            req.name,
            req.description,
            req.template,
            req.input_variables,
        )
        .await
    {
        Ok(id) => ApiResponse {
            code: 200,
            message: None,
            data: Some(UpdatePromptTemplateRes { id }),
        },
        Err(e) => ApiResponse {
            code: 500,
            message: Some(e),
            data: None,
        },
    }
}
