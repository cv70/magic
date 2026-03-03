#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Task domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Task {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub task_type: String,
    pub scheduler_id: i64,
    pub cron_expression: Option<String>,
    pub next_run_at: Option<String>,
    pub last_run_at: Option<String>,
    pub last_run_status: Option<String>,
    pub last_run_error: Option<String>,
    pub last_run_duration: Option<i64>,
    pub last_run_result: Option<String>,
    pub enabled: Option<bool>,
    pub created_at: Option<String>,
    pub updated_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetTaskRes {
    pub task: Task,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchTasksReq {
    pub name: Option<String>,
    pub task_type: Option<String>,
    pub scheduler_id: Option<i64>,
    pub enabled: Option<bool>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchTasksRes {
    pub tasks: Vec<Task>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddTaskReq {
    pub name: String,
    pub description: Option<String>,
    pub task_type: String,
    pub scheduler_id: i64,
    pub cron_expression: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddTaskRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateTaskReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
    pub task_type: Option<String>,
    pub scheduler_id: Option<i64>,
    pub cron_expression: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateTaskRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RunTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RunTaskRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StopTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct StopTaskRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RestartTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct RestartTaskRes {
    pub id: i64,
}
