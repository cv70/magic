#![allow(dead_code)]
use axum::response::{IntoResponse, Json};
use serde::{Deserialize, Serialize};

/// HTTP状态码类型别名
type HTTPCode = i32;

pub const HTTP_CODE_OK: HTTPCode = 200;
pub const HTTP_CODE_BAD_REQUEST: HTTPCode = 400;
pub const HTTP_CODE_INTERNAL_SERVER_ERROR: HTTPCode = 500;

/// 标准化的HTTP响应包装对象
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct ApiResponse<T> {
    /// 状态码
    pub code: HTTPCode,
    /// 响应消息
    #[serde(skip_serializing_if = "Option::is_none")]
    pub message: Option<String>,
    /// 响应数据
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<T>,
}

impl<T: Serialize> IntoResponse for ApiResponse<T> {
    fn into_response(self) -> axum::response::Response {
        Json(self).into_response()
    }
}

/// 成功响应，数据可选
pub fn ok<T: Serialize>(data: Option<T>) -> impl IntoResponse {
    Json(ApiResponse {
        code: HTTP_CODE_OK,
        message: None,
        data: data,
    })
}

/// 错误响应（客户端错误）
pub fn error(code: HTTPCode, message: String) -> impl IntoResponse {
    Json(ApiResponse::<()> {
        code: code,
        message: Some(message),
        data: None,
    })
}
