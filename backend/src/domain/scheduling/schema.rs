#![allow(dead_code)]
use serde::{Deserialize, Serialize};

// Scheduling domain schemas

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Scheduler {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub cron_expression: String,
    pub enabled: Option<bool>,
    pub next_run_at: Option<String>,
    pub last_run_at: Option<String>,
    pub last_run_duration: Option<i64>,
    pub last_run_status: Option<String>,
    pub last_run_error: Option<String>,
    pub last_run_result: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct CronExpression {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub expression: String,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct ScheduledTask {
    pub id: i64,
    pub scheduler_id: i64,
    pub task_name: String,
    pub task_type: String,
    pub task_config: Option<String>,
    pub status: Option<String>,
    pub error: Option<String>,
    pub created_at: Option<String>,
    pub started_at: Option<String>,
    pub completed_at: Option<String>,
    pub next_run_at: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetSchedulerReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetSchedulerRes {
    pub scheduler: Scheduler,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchSchedulersReq {
    pub name: Option<String>,
    pub enabled: Option<bool>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchSchedulersRes {
    pub schedulers: Vec<Scheduler>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddSchedulerReq {
    pub name: String,
    pub description: Option<String>,
    pub cron_expression: String,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddSchedulerRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateSchedulerReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
    pub cron_expression: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateSchedulerRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetCronExpressionReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetCronExpressionRes {
    pub cron_expression: CronExpression,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchCronExpressionsReq {
    pub name: Option<String>,
    pub enabled: Option<bool>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchCronExpressionsRes {
    pub cron_expressions: Vec<CronExpression>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddCronExpressionReq {
    pub name: String,
    pub description: Option<String>,
    pub expression: String,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddCronExpressionRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateCronExpressionReq {
    pub id: i64,
    pub name: Option<String>,
    pub description: Option<String>,
    pub expression: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateCronExpressionRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetScheduledTaskReq {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct GetScheduledTaskRes {
    pub scheduled_task: ScheduledTask,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchScheduledTasksReq {
    pub scheduler_id: Option<i64>,
    pub task_type: Option<String>,
    pub status: Option<String>,
    pub page: i64,
    pub limit: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct SearchScheduledTasksRes {
    pub scheduled_tasks: Vec<ScheduledTask>,
    pub total: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddScheduledTaskReq {
    pub scheduler_id: i64,
    pub task_name: String,
    pub task_type: String,
    pub task_config: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AddScheduledTaskRes {
    pub id: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateScheduledTaskReq {
    pub id: i64,
    pub scheduler_id: i64,
    pub task_name: Option<String>,
    pub task_type: Option<String>,
    pub task_config: Option<String>,
    pub enabled: Option<bool>,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UpdateScheduledTaskRes {
    pub id: i64,
}
