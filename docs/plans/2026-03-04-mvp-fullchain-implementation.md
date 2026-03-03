# Fullchain MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在沿用现有后端 POST 接口的前提下，实现内容管理、AI 生成、发布管理的前后端可用 MVP。

**Architecture:** 后端仅做最小补充（内容删除 API）；前端改为单页控制台，按模块调用既有 API。通过统一请求层处理标准响应，降低页面复杂度与重复逻辑。

**Tech Stack:** Rust + Axum, React + TypeScript + Vite

---

### Task 1: 后端补充内容删除能力

**Files:**
- Modify: `backend/src/domain/content/schema.rs`
- Modify: `backend/src/domain/content/api.rs`
- Modify: `backend/src/domain/content/domain.rs`
- Modify: `backend/src/datasource/dbdao/content.rs`
- Modify: `backend/src/main.rs`

**Step 1: 增加删除请求/响应 schema**
- 新增 `DeleteContentReq { id: i64 }`
- 新增 `DeleteContentRes { id: i64 }`

**Step 2: 增加 domain 删除方法**
- `delete_content(id: i64) -> Result<i64, String>`

**Step 3: 增加 dao 删除 SQL**
- 执行 `DELETE FROM content WHERE id = $1 RETURNING id`

**Step 4: 增加 API handler + 路由**
- `api_delete_content`
- 注册 `/api/v1/content/delete`

**Step 5: 编译验证**
- Run: `cd backend && cargo check`

### Task 2: 前端构建统一 API 层

**Files:**
- Create: `frontend/src/api.ts`

**Step 1: 定义通用响应类型**
- `ApiResponse<T>`

**Step 2: 封装 request 函数**
- 统一 POST JSON
- 处理 `code !== 200` 与网络错误

**Step 3: 定义 MVP 业务 API**
- 内容、AI、发布模块方法

**Step 4: 类型检查**
- Run: `cd frontend && npm run build`

### Task 3: 前端实现单页控制台

**Files:**
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/App.css`

**Step 1: 页面结构**
- 三个区块：内容管理、AI 生成、发布管理

**Step 2: 内容管理交互**
- 列表、创建、更新、删除

**Step 3: AI 生成交互**
- 加载生成器、提交生成、展示任务 ID、保存为内容

**Step 4: 发布管理交互**
- 选择发布器+内容发起发布、展示任务列表

**Step 5: 样式与移动端适配**
- 响应式卡片布局

### Task 4: 验证

**Files:**
- N/A

**Step 1: 后端检查**
- Run: `cd backend && cargo check`

**Step 2: 前端构建**
- Run: `cd frontend && npm run build`

**Step 3: 结果记录**
- 汇总可用功能与未覆盖项
