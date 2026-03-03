mod config;
mod datasource;
mod domain;
mod infra;
mod state;
mod utils;

use anyhow::Result;
use axum::{Router, routing::post};

use crate::config::config::AppConfig;
use crate::domain::ai_generation::api::{
    api_add_generator, api_add_prompt_template, api_generate_content, api_get_generator,
    api_get_prompt_template, api_search_generators, api_search_prompt_templates,
    api_update_generator, api_update_prompt_template,
};
use crate::domain::configuration::api::{
    api_add_provider_config, api_add_system_config, api_get_provider_config, api_get_system_config,
    api_search_provider_configs, api_search_system_configs, api_update_provider_config,
    api_update_system_config,
};
use crate::domain::content::api::{
    api_add_content, api_delete_content, api_get_content, api_search_contents, api_update_content,
};
use crate::domain::financing::api::{
    api_create_business_plan, api_delete_business_plan, api_get_business_plan,
    api_get_business_plans, api_update_business_plan,
};
use crate::domain::identity::api::{
    api_add_permission, api_add_role, api_add_user, api_get_permission, api_get_role, api_get_user,
    api_search_permissions, api_search_roles, api_search_users, api_update_permission,
    api_update_role, api_update_user,
};
use crate::domain::news::api::{api_add_news, api_get_news, api_search_news};
use crate::domain::publishing::api::{
    api_add_publisher, api_get_publish_log, api_get_publish_task, api_get_publisher,
    api_publish_content, api_search_publish_logs, api_search_publish_tasks, api_search_publishers,
    api_update_publisher,
};
use crate::domain::scheduling::api::{
    api_add_task, api_get_task, api_restart_task, api_run_task, api_search_tasks, api_stop_task,
    api_update_task,
};
use crate::infra::registry::Registry;
use crate::state::state::AppState;

fn parse_config_path_from_iter<I>(args: I) -> Result<String>
where
    I: IntoIterator<Item = String>,
{
    let mut args = args.into_iter();
    while let Some(arg) = args.next() {
        if arg == "--config" {
            let Some(path) = args.next() else {
                return Err(anyhow::anyhow!("--config requires a file path"));
            };
            return Ok(path);
        }
    }
    Err(anyhow::anyhow!("missing required --config <path>"))
}

fn parse_config_path() -> Result<String> {
    parse_config_path_from_iter(std::env::args().skip(1))
}

fn api_v1_route() -> Router<AppState> {
    Router::new()
        .route("/finance/business/plans", post(api_get_business_plans))
        .route("/finance/business/plan", post(api_get_business_plan))
        .route("/finance/business/create", post(api_create_business_plan))
        .route("/finance/business/update", post(api_update_business_plan))
        .route("/finance/business/delete", post(api_delete_business_plan))
        .route("/news/get", post(api_get_news))
        .route("/news/search", post(api_search_news))
        .route("/news/add", post(api_add_news))
        .route("/content/get", post(api_get_content))
        .route("/content/search", post(api_search_contents))
        .route("/content/add", post(api_add_content))
        .route("/content/update", post(api_update_content))
        .route("/content/delete", post(api_delete_content))
        .route("/config/system/get", post(api_get_system_config))
        .route("/config/system/search", post(api_search_system_configs))
        .route("/config/system/add", post(api_add_system_config))
        .route("/config/system/update", post(api_update_system_config))
        .route("/config/provider/get", post(api_get_provider_config))
        .route("/config/provider/search", post(api_search_provider_configs))
        .route("/config/provider/add", post(api_add_provider_config))
        .route("/config/provider/update", post(api_update_provider_config))
        .route("/ai/generator/get", post(api_get_generator))
        .route("/ai/generator/search", post(api_search_generators))
        .route("/ai/generator/add", post(api_add_generator))
        .route("/ai/generator/update", post(api_update_generator))
        .route("/ai/generate", post(api_generate_content))
        .route("/ai/template/get", post(api_get_prompt_template))
        .route("/ai/template/search", post(api_search_prompt_templates))
        .route("/ai/template/add", post(api_add_prompt_template))
        .route("/ai/template/update", post(api_update_prompt_template))
        .route("/publishing/publisher/get", post(api_get_publisher))
        .route("/publishing/publisher/search", post(api_search_publishers))
        .route("/publishing/publisher/add", post(api_add_publisher))
        .route("/publishing/publisher/update", post(api_update_publisher))
        .route("/publishing/content/publish", post(api_publish_content))
        .route("/publishing/task/get", post(api_get_publish_task))
        .route("/publishing/task/search", post(api_search_publish_tasks))
        .route("/publishing/log/get", post(api_get_publish_log))
        .route("/publishing/log/search", post(api_search_publish_logs))
        .route("/scheduling/task/get", post(api_get_task))
        .route("/scheduling/task/search", post(api_search_tasks))
        .route("/scheduling/task/add", post(api_add_task))
        .route("/scheduling/task/update", post(api_update_task))
        .route("/scheduling/task/run", post(api_run_task))
        .route("/scheduling/task/stop", post(api_stop_task))
        .route("/scheduling/task/restart", post(api_restart_task))
        .route("/identity/user/get", post(api_get_user))
        .route("/identity/user/search", post(api_search_users))
        .route("/identity/user/add", post(api_add_user))
        .route("/identity/user/update", post(api_update_user))
        .route("/identity/role/get", post(api_get_role))
        .route("/identity/role/search", post(api_search_roles))
        .route("/identity/role/add", post(api_add_role))
        .route("/identity/role/update", post(api_update_role))
        .route("/identity/permission/get", post(api_get_permission))
        .route("/identity/permission/search", post(api_search_permissions))
        .route("/identity/permission/add", post(api_add_permission))
        .route("/identity/permission/update", post(api_update_permission))
}

fn app_route() -> Router<AppState> {
    Router::new().nest("/api/v1", api_v1_route())
}

#[tokio::main]
async fn main() -> Result<()> {
    let config_path = parse_config_path()?;
    let config = AppConfig::load_from_yaml_file(&config_path)?;
    let registry = Registry::new(&config).await?;
    let app_state = AppState::new(&registry);

    let app = app_route().with_state(app_state);
    let addr = "0.0.0.0:8888";
    println!("Starting server on {}", addr);
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::parse_config_path_from_iter;

    #[test]
    fn parse_config_path_ok() {
        let args = vec!["--config".to_string(), "/tmp/a.yaml".to_string()];
        let path = parse_config_path_from_iter(args).expect("parse args");
        assert_eq!(path, "/tmp/a.yaml".to_string());
    }

    #[test]
    fn parse_config_path_missing_value_error() {
        let args = vec!["--config".to_string()];
        let err = parse_config_path_from_iter(args).expect_err("should fail");
        assert!(
            err.to_string().contains("--config requires a file path"),
            "unexpected error: {err}"
        );
    }

    #[test]
    fn parse_config_path_required_error_when_absent() {
        let args: Vec<String> = vec![];
        let err = parse_config_path_from_iter(args).expect_err("should fail");
        assert!(
            err.to_string().contains("missing required --config <path>"),
            "unexpected error: {err}"
        );
    }
}
