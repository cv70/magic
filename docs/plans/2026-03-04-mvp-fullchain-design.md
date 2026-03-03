# 全链路 MVP 设计文档（沿用现有接口）

## 背景
基于 `ARCHITECTURE.md`，本次目标是在不重构后端接口风格的前提下，打通“内容管理 + AI 生成 + 发布管理”的前后端最小可用闭环。

## 设计结论
- 后端沿用现有 `/api/v1/**` POST 接口，不新增 REST 风格路由。
- 前端实现单页控制台，不引入鉴权、复杂路由和状态管理库。
- 通过统一的前端 API 客户端处理 `ApiResponse { code, message, data }`。

## 模块设计

### 1. 内容管理
- 列表：`POST /api/v1/content/search`
- 新建：`POST /api/v1/content/add`
- 更新：`POST /api/v1/content/update`
- 删除：补充 `POST /api/v1/content/delete`

### 2. AI 生成
- 生成器列表：`POST /api/v1/ai/generator/search`
- 生成任务：`POST /api/v1/ai/generate`
- MVP 中“生成结果”以任务 ID 展示，并支持将用户输入快速保存为内容。

### 3. 发布管理
- 发布器列表：`POST /api/v1/publishing/publisher/search`
- 发起发布：`POST /api/v1/publishing/content/publish`
- 任务列表：`POST /api/v1/publishing/task/search`

## 数据流
1. 前端初始化并行拉取内容列表、生成器、发布器、发布任务。
2. 用户发起 AI 生成后，返回任务 ID；可选择将文本保存为内容。
3. 用户选择内容与发布器后发起发布，发布任务列表刷新。

## 错误处理
- 前端统一处理：
  - `code !== 200` 显示后端 message。
  - 网络错误显示“连接失败”。
- 后端保持现有响应结构，不做全局错误模型改造。

## 测试与验证
- 后端：补充删除接口对应 schema 与编译验证。
- 前端：执行 TypeScript 构建，确保类型与打包通过。
- 联调：通过页面交互验证三条主链路请求成功。

## 非目标
- 本次不做 REST 兼容层。
- 不做登录鉴权。
- 不引入复杂 UI 框架。
