# 技术设计文档（第一阶段）

> 本文档详细描述第一阶段核心功能的技术实现方案，包括数据库设计、API定义、前后端组件设计。

## 目录

1. [核心数据模型](#核心数据模型)
2. [后端API设计](#后端api设计)
3. [前端组件结构](#前端组件结构)
4. [核心功能技术方案](#核心功能技术方案)
5. [开发规范和最佳实践](#开发规范和最佳实践)

---

## 核心数据模型

### 1. Draft（内容草稿）

```go
type Draft struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    UserID        uint           `gorm:"index" json:"user_id"`
    Title         string         `gorm:"size:255" json:"title"`
    Content       string         `gorm:"type:longtext" json:"content"` // HTML或Markdown
    ContentType   string         `gorm:"size:20" json:"content_type"`   // text/code/image等
    Status        string         `gorm:"size:20" json:"status"`         // draft/archived
    Tags          string         `gorm:"type:json" json:"tags"`         // JSON数组
    Metadata      string         `gorm:"type:json" json:"metadata"`     // 自定义元数据
    LastEditedAt  time.Time      `json:"last_edited_at"`                // 用于排序
    SavedVersionID *uint         `json:"saved_version_id"`              // 保存的版本
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
}

// 索引
// - UNIQUE (user_id, title)
// - INDEX (user_id, last_edited_at)
// - INDEX (user_id, status)
```

### 2. Content（正式发布的内容）

```go
type Content struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"index" json:"user_id"`
    Title       string    `gorm:"size:255" json:"title"`
    Body        string    `gorm:"type:longtext" json:"body"` // HTML
    ContentType string    `gorm:"size:20" json:"content_type"`
    Status      string    `gorm:"size:20" json:"status"`
    Tags        string    `gorm:"type:json" json:"tags"`
    Metadata    string    `gorm:"type:json" json:"metadata"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    PublishedAt *time.Time `json:"published_at"`
}

// 索引
// - INDEX (user_id, created_at)
// - INDEX (user_id, status)
// - FULLTEXT (title, body) -- 全文搜索（MySQL）
```

### 3. ContentVersion（版本历史）

```go
type ContentVersion struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    ContentID   uint      `gorm:"index" json:"content_id"`
    DraftID     *uint     `json:"draft_id"`
    VersionNum  int       `json:"version_num"` // 1, 2, 3...
    Title       string    `gorm:"size:255" json:"title"`
    Content     string    `gorm:"type:longtext" json:"content"`
    ChangeSummary string  `gorm:"size:500" json:"change_summary"` // 本次修改的摘要
    CreatedBy   uint      `json:"created_by"` // 修改人
    CreatedAt   time.Time `json:"created_at"`
}

// 索引
// - INDEX (content_id, version_num)
// - INDEX (draft_id)
```

### 4. Tag（标签）

```go
type Tag struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"index" json:"user_id"`
    Name      string    `gorm:"size:50" json:"name"`
    Color     string    `gorm:"size:20" json:"color"` // 配色
    Count     int       `json:"count"` // 使用次数
    CreatedAt time.Time `json:"created_at"`
}

// 索引
// - UNIQUE (user_id, name)
```

### 5. Category（分类）

```go
type Category struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"index" json:"user_id"`
    Name      string    `gorm:"size:100" json:"name"`
    Icon      string    `gorm:"size:50" json:"icon"`
    Color     string    `gorm:"size:20" json:"color"`
    Order     int       `json:"order"` // 排序
    CreatedAt time.Time `json:"created_at"`
}

// 索引
// - INDEX (user_id, order)
```

### 6. PublishTask（发布任务）

```go
type PublishTask struct {
    ID                 uint      `gorm:"primaryKey" json:"id"`
    UserID             uint      `gorm:"index" json:"user_id"`
    ContentID          uint      `gorm:"index" json:"content_id"`
    DraftID            *uint     `json:"draft_id"`
    TargetPlatforms    string    `gorm:"type:json" json:"target_platforms"` // ["juejin", "csdn"]
    Status             string    `gorm:"size:20" json:"status"` // pending/publishing/published/failed
    ScheduledAt        *time.Time `json:"scheduled_at"` // 计划发布时间
    PublishedAt        *time.Time `json:"published_at"` // 实际发布时间
    RetryCount         int       `json:"retry_count"`
    MaxRetries         int       `json:"max_retries"`
    ErrorMessage       string    `gorm:"type:text" json:"error_message"`
    PlatformResults    string    `gorm:"type:json" json:"platform_results"` // 各平台结果
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}

// PlatformResult 结构
type PlatformResult struct {
    Platform   string    `json:"platform"` // juejin/csdn等
    Status     string    `json:"status"` // published/failed
    PostID     string    `json:"post_id"`
    PostURL    string    `json:"post_url"`
    Error      string    `json:"error,omitempty"`
    PublishedAt time.Time `json:"published_at,omitempty"`
    Metrics    Metrics   `json:"metrics,omitempty"`
}

type Metrics struct {
    Views    int `json:"views"`
    Likes    int `json:"likes"`
    Comments int `json:"comments"`
    Shares   int `json:"shares"`
    Follows  int `json:"follows"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 索引
// - INDEX (user_id, scheduled_at)
// - INDEX (user_id, status)
// - INDEX (content_id)
```

### 7. PublishAnalytic（发布数据统计）

```go
type PublishAnalytic struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    PublishTaskID   uint      `gorm:"index" json:"publish_task_id"`
    ContentID       uint      `gorm:"index" json:"content_id"`
    UserID          uint      `gorm:"index" json:"user_id"`
    Platform        string    `gorm:"size:50" json:"platform"`
    PostID          string    `gorm:"size:255" json:"post_id"` // 各平台的文章ID
    Title           string    `gorm:"size:255" json:"title"`
    PostURL         string    `gorm:"type:text" json:"post_url"`
    Metrics         string    `gorm:"type:json" json:"metrics"` // JSON: {views, likes, comments...}
    LastSyncedAt    time.Time `json:"last_synced_at"` // 最后同步数据的时间
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 索引
// - INDEX (user_id, created_at)
// - INDEX (user_id, platform)
// - INDEX (publish_task_id)
```

### 8. PromptTemplate（AI提示词模板）

```go
type PromptTemplate struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    UserID         *uint     `json:"user_id"` // NULL表示系统模板
    Category       string    `gorm:"size:100" json:"category"` // text_generation/social_media等
    SubCategory    string    `gorm:"size:100" json:"sub_category"` // xiaohongshu/weibo等
    Name           string    `gorm:"size:255" json:"name"`
    Description    string    `gorm:"type:text" json:"description"`
    PromptText     string    `gorm:"type:longtext" json:"prompt_text"` // 实际提示词
    Variables      string    `gorm:"type:json" json:"variables"` // ["{keyword}", "{tone}"]
    ExampleInput   string    `gorm:"type:json" json:"example_input"` // 示例输入
    ExampleOutput  string    `gorm:"type:text" json:"example_output"` // 示例输出
    UsageCount     int       `json:"usage_count"`
    Rating         float32   `json:"rating"` // 1-5星
    Tags           string    `gorm:"type:json" json:"tags"`
    Status         string    `gorm:"size:20" json:"status"` // active/archived
    Version        int       `json:"version"`
    CreatedBy      uint      `json:"created_by"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// 索引
// - INDEX (category, sub_category)
// - INDEX (user_id)
// - FULLTEXT (name, description)
```

---

## 后端API设计

### 1. Draft API

```
▶ 草稿管理相关API

POST /api/v1/drafts
创建新草稿
请求：
{
    "title": "我的新文章",
    "content": "<p>内容</p>",
    "content_type": "text",
    "tags": ["go", "性能"],
    "category_id": 1
}
响应：
{
    "id": 1,
    "user_id": 10,
    "title": "我的新文章",
    "content": "<p>内容</p>",
    "created_at": "2026-03-07T10:00:00Z"
}

─────────────────────────

GET /api/v1/drafts
获取草稿列表（支持分页、筛选）
查询参数：
    page: 1
    limit: 20
    status: draft/archived
    keyword: 搜索标题
    tags: go,performance (逗号分割)
响应：
{
    "total": 150,
    "data": [
        {
            "id": 1,
            "title": "...",
            "last_edited_at": "2026-03-07T10:00:00Z"
        }
    ]
}

─────────────────────────

GET /api/v1/drafts/{id}
获取草稿详情
响应：
{
    "id": 1,
    "title": "...",
    "content": "...",
    "tags": [...],
    "versions": [
        {
            "version_num": 1,
            "change_summary": "初始版本",
            "created_at": "..."
        }
    ]
}

─────────────────────────

PUT /api/v1/drafts/{id}
更新草稿（自动保存）
请求：
{
    "title": "...",
    "content": "...",
    "tags": [...],
    "change_summary": "修改了第3段" // 可选
}
响应：
{
    "id": 1,
    "updated_at": "2026-03-07T10:30:00Z",
    "version_num": 5
}

─────────────────────────

POST /api/v1/drafts/{id}/publish
将草稿发布为正式内容
请求：
{
    "status": "published" // 或 "archived"
}
响应：
{
    "content_id": 100,
    "status": "published",
    "published_at": "2026-03-07T11:00:00Z"
}

─────────────────────────

DELETE /api/v1/drafts/{id}
删除草稿
响应：
{
    "success": true,
    "message": "草稿已删除"
}

─────────────────────────

GET /api/v1/drafts/{id}/versions
获取草稿的版本历史
响应：
{
    "data": [
        {
            "version_num": 5,
            "change_summary": "修改了第3段",
            "created_by_user": "张三",
            "created_at": "2026-03-07T10:30:00Z"
        },
        {
            "version_num": 4,
            "change_summary": "修复了链接",
            "created_at": "2026-03-07T10:15:00Z"
        }
    ]
}

─────────────────────────

POST /api/v1/drafts/{id}/revert/{version_num}
恢复到某个版本
响应：
{
    "id": 1,
    "version_num": 6, // 恢复操作本身也会创建新版本
    "reverted_from": 5,
    "change_summary": "恢复到版本5"
}

─────────────────────────

GET /api/v1/drafts/{id}/diff
版本对比（两个版本之间的差异）
查询参数：
    from_version: 4
    to_version: 5
响应：
{
    "from_version": 4,
    "to_version": 5,
    "diff": {
        "title": {
            "old": "旧标题",
            "new": "新标题",
            "changed": true
        },
        "content": {
            "old_html": "...",
            "new_html": "...",
            "changes": [
                {
                    "type": "modified",
                    "position": "第3段",
                    "old": "...",
                    "new": "..."
                }
            ]
        }
    }
}
```

### 2. Content API

```
▶ 正式内容相关API

GET /api/v1/contents
获取内容列表（已发布的）
查询参数：
    page: 1
    limit: 20
    status: published/archived
    content_type: text/code/image
    tag: "go"
    from_date: "2026-03-01"
    to_date: "2026-03-07"
响应：
{
    "total": 240,
    "data": [
        {
            "id": 100,
            "title": "Go性能优化",
            "content_type": "text",
            "status": "published",
            "published_at": "2026-03-07T11:00:00Z",
            "publish_tasks_count": 2 // 发布到几个平台
        }
    ]
}

─────────────────────────

GET /api/v1/contents/{id}
获取内容详情
响应：
{
    "id": 100,
    "title": "Go性能优化",
    "body": "...",
    "published_at": "2026-03-07T11:00:00Z",
    "tags": ["go", "性能"],
    "category": "技术",
    "publish_tasks": [
        {
            "id": 200,
            "platform": "juejin",
            "status": "published",
            "post_url": "..."
        }
    ],
    "analytics": {
        "total_views": 1200,
        "total_likes": 45,
        "total_comments": 8
    }
}

─────────────────────────

DELETE /api/v1/contents/{id}
删除内容（软删除）
响应：
{
    "success": true
}
```

### 3. PublishTask API

```
▶ 发布任务相关API

POST /api/v1/publish-tasks
创建发布任务
请求：
{
    "content_id": 100, // 或 draft_id: 1
    "target_platforms": ["juejin", "csdn", "weibo"],
    "scheduled_at": "2026-03-08T14:30:00Z", // 可选，为空表示立即发布
    "max_retries": 3
}
响应：
{
    "id": 200,
    "status": "pending",
    "scheduled_at": "2026-03-08T14:30:00Z",
    "platform_results": {
        "juejin": {"status": "pending"},
        "csdn": {"status": "pending"},
        "weibo": {"status": "pending"}
    }
}

─────────────────────────

GET /api/v1/publish-tasks
获取发布任务列表
查询参数：
    page: 1
    limit: 20
    status: pending/publishing/published/failed
    from_date: "2026-03-01"
    to_date: "2026-03-07"
响应：
{
    "data": [
        {
            "id": 200,
            "content_id": 100,
            "content_title": "Go性能优化",
            "status": "published",
            "scheduled_at": "2026-03-08T14:30:00Z",
            "published_at": "2026-03-08T14:30:45Z",
            "platforms": [
                {
                    "platform": "juejin",
                    "status": "published",
                    "post_id": "abc123",
                    "post_url": "https://juejin.cn/post/abc123",
                    "metrics": {
                        "views": 450,
                        "likes": 12,
                        "comments": 3
                    }
                }
            ]
        }
    ]
}

─────────────────────────

GET /api/v1/publish-tasks/{id}
获取发布任务详情
响应：
{
    "id": 200,
    "content_id": 100,
    "status": "publishing", // 实时更新
    "platform_results": {
        "juejin": {
            "status": "published",
            "post_id": "abc123",
            "post_url": "...",
            "metrics": {...}
        },
        "csdn": {
            "status": "publishing",
            "start_time": "2026-03-08T14:30:45Z"
        },
        "weibo": {
            "status": "failed",
            "error": "API rate limit exceeded",
            "next_retry_at": "2026-03-08T14:35:00Z"
        }
    }
}

─────────────────────────

PUT /api/v1/publish-tasks/{id}
修改发布任务（改期等）
请求：
{
    "scheduled_at": "2026-03-09T10:00:00Z"
}
响应：
{
    "id": 200,
    "scheduled_at": "2026-03-09T10:00:00Z"
}

─────────────────────────

POST /api/v1/publish-tasks/{id}/execute
立即发起发布（不等待scheduled_at）
响应：
{
    "id": 200,
    "status": "publishing"
}

─────────────────────────

DELETE /api/v1/publish-tasks/{id}
取消发布任务（只能在pending状态）
响应：
{
    "success": true
}

─────────────────────────

GET /api/v1/publish-tasks/{id}/stream
实时推送发布状态（WebSocket或Server-Sent Events）
响应（持续推送）：
{
    "event": "platform_published",
    "platform": "juejin",
    "status": "published",
    "post_url": "...",
    "metrics": {...}
}
```

### 4. Tag & Category API

```
▶ 标签和分类API

GET /api/v1/tags
获取所有标签
响应：
{
    "data": [
        {"id": 1, "name": "go", "count": 25},
        {"id": 2, "name": "性能", "count": 18}
    ]
}

─────────────────────────

POST /api/v1/tags
创建标签
请求：
{
    "name": "rust",
    "color": "#FF5733"
}
响应：
{
    "id": 3,
    "name": "rust",
    "color": "#FF5733"
}

─────────────────────────

DELETE /api/v1/tags/{id}
删除标签
响应：
{
    "success": true
}

─────────────────────────

GET /api/v1/categories
获取所有分类
响应：
{
    "data": [
        {"id": 1, "name": "技术", "icon": "code", "order": 1},
        {"id": 2, "name": "生活", "icon": "heart", "order": 2}
    ]
}

─────────────────────────

POST /api/v1/categories
创建分类
请求：
{
    "name": "技术",
    "icon": "code",
    "color": "#0066FF"
}
响应：
{
    "id": 1,
    "name": "技术"
}

─────────────────────────

PUT /api/v1/categories/{id}/order
修改分类顺序
请求：
{
    "order": 3
}
```

### 5. PromptTemplate API

```
▶ AI提示词模板API

GET /api/v1/prompt-templates
获取提示词模板（支持分类筛选）
查询参数：
    category: text_generation
    sub_category: xiaohongshu
    user_templates: true/false (个人模板/系统模板)
响应：
{
    "data": [
        {
            "id": 1,
            "category": "text_generation",
            "sub_category": "xiaohongshu",
            "name": "小红书种草笔记",
            "description": "为产品写吸引人的小红书笔记",
            "variables": ["{product_name}", "{features}", "{tone}"],
            "usage_count": 125,
            "rating": 4.5
        }
    ]
}

─────────────────────────

GET /api/v1/prompt-templates/{id}
获取模板详情
响应：
{
    "id": 1,
    "name": "小红书种草笔记",
    "prompt_text": "作为内容创作者，为{product_name}写一篇小红书笔记...",
    "variables": ["{product_name}", "{features}", "{tone}"],
    "example_input": {
        "product_name": "iPhone 15",
        "features": "钛金属、A17芯片",
        "tone": "可爱"
    },
    "example_output": "从iPhone 15开箱那一刻我就..."
}

─────────────────────────

POST /api/v1/prompt-templates
创建自定义模板
请求：
{
    "name": "我的模板",
    "category": "text_generation",
    "prompt_text": "...",
    "variables": ["{keyword}"],
    "example_input": {...},
    "example_output": "..."
}
响应：
{
    "id": 999,
    "user_id": 10,
    "name": "我的模板"
}
```

### 6. Analytics API

```
▶ 数据分析API

GET /api/v1/analytics/summary
获取统计摘要
响应：
{
    "period": "2026-03-01 ~ 2026-03-07",
    "total_published": 8, // 本周发布数
    "total_views": 5240,
    "total_likes": 245,
    "total_comments": 82,
    "avg_likes_per_post": 30.6,
    "total_new_followers": 45
}

─────────────────────────

GET /api/v1/analytics/content-ranking
获取内容性能排行
查询参数：
    metric: views/likes/comments
    days: 30
    limit: 10
响应：
{
    "data": [
        {
            "rank": 1,
            "content_id": 100,
            "title": "Go性能优化",
            "views": 1200,
            "likes": 45,
            "comments": 8,
            "shares": 5
        }
    ]
}

─────────────────────────

GET /api/v1/analytics/platform-comparison
平台对比分析
响应：
{
    "data": [
        {
            "platform": "juejin",
            "posts": 8,
            "avg_views": 450,
            "avg_likes": 15,
            "success_rate": "87%"
        },
        {
            "platform": "csdn",
            "posts": 7,
            "avg_views": 320,
            "avg_likes": 12,
            "success_rate": "71%"
        }
    ]
}

─────────────────────────

GET /api/v1/analytics/best-publish-time
获取最佳发布时间建议
响应：
{
    "recommendations": [
        {
            "time": "周二 14:00-16:00",
            "reason": "该时段平均点赞提升 45%",
            "performance": {
                "avg_likes": 35,
                "avg_views": 520
            }
        },
        {
            "time": "周四 20:00-22:00",
            "reason": "该时段转发数最多",
            "performance": {
                "avg_shares": 8,
                "avg_views": 450
            }
        }
    ]
}
```

---

## 前端组件结构

### 目录结构

```
frontend/src/
├── pages/
│   ├── Editor/              # 编辑页面
│   │   ├── index.tsx       # 编辑器主页面
│   │   ├── RichEditor.tsx  # 富文本编辑器组件
│   │   ├── Sidebar.tsx     # 左侧面板（标签、分类等）
│   │   ├── RightPanel.tsx  # 右侧面板（SEO、字数等）
│   │   ├── DraftAutoSave.tsx # 自动保存逻辑
│   │   └── EditorStyles.css
│   │
│   ├── DraftCenter/         # 草稿中心
│   │   ├── index.tsx       # 草稿列表
│   │   ├── DraftList.tsx   # 列表组件
│   │   ├── DraftItem.tsx   # 单个草稿项
│   │   ├── DraftSearch.tsx # 搜索筛选
│   │   └── DraftStyles.css
│   │
│   ├── PublishManager/      # 发布管理
│   │   ├── index.tsx
│   │   ├── PublishCalendar.tsx  # 发布日历
│   │   ├── PublishQueue.tsx     # 发布队列
│   │   ├── PublishTask.tsx      # 单个发布任务卡片
│   │   ├── PublishMonitor.tsx   # 发布监控（实时状态）
│   │   └── PublishStyles.css
│   │
│   ├── Analytics/           # 数据分析
│   │   ├── index.tsx
│   │   ├── Dashboard.tsx    # 看板
│   │   ├── Ranking.tsx      # 排行榜
│   │   ├── PlatformCompare.tsx # 平台对比
│   │   └── AnalyticsStyles.css
│   │
│   └── Content/             # 已发布内容
│       ├── index.tsx
│       └── ContentList.tsx
│
├── components/
│   ├── common/
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   ├── Tabs.tsx
│   │   ├── Card.tsx
│   │   ├── Table.tsx
│   │   └── Loading.tsx
│   │
│   ├── editor/
│   │   ├── RichTextEditor.tsx   # 富文本编辑器封装
│   │   ├── Toolbar.tsx          # 编辑器工具栏
│   │   ├── AutoSave.tsx         # 自动保存组件
│   │   ├── VersionHistory.tsx   # 版本历史面板
│   │   ├── VersionDiff.tsx      # 版本对比
│   │   └── ContentPreview.tsx   # 多平台预览
│   │
│   ├── publish/
│   │   ├── PlatformSelector.tsx # 平台选择器
│   │   ├── ScheduleDatePicker.tsx # 定时选择器
│   │   ├── PublishProgress.tsx  # 发布进度条
│   │   ├── PlatformStatus.tsx   # 各平台发布状态
│   │   └── PublishRecommendation.tsx # 发布建议
│   │
│   ├── tags/
│   │   ├── TagInput.tsx         # 标签输入框
│   │   ├── TagCloud.tsx         # 标签云
│   │   ├── CategorySelector.tsx # 分类选择器
│   │   └── TagManager.tsx       # 标签管理
│   │
│   ├── analytics/
│   │   ├── Chart.tsx            # 图表基础组件
│   │   ├── LineChart.tsx        # 折线图
│   │   ├── BarChart.tsx         # 柱状图
│   │   ├── PieChart.tsx         # 饼图
│   │   └── HeatMap.tsx          # 热力图
│   │
│   └── draft/
│       ├── DraftList.tsx
│       └── DraftCard.tsx
│
├── hooks/
│   ├── useDraft.ts              # 草稿相关hook
│   ├── usePublish.ts            # 发布相关hook
│   ├── useAnalytics.ts          # 分析相关hook
│   ├── useAutoSave.ts           # 自动保存逻辑
│   ├── useVersionHistory.ts     # 版本历史逻辑
│   ├── useContentPreview.ts     # 多平台预览逻辑
│   └── useWebSocket.ts          # WebSocket实时推送
│
├── utils/
│   ├── api.ts                   # API请求封装
│   ├── format.ts                # 格式化工具（时间、字数等）
│   ├── editor.ts                # 编辑器相关工具
│   ├── diff.ts                  # 差异对比工具
│   ├── storage.ts               # 本地存储工具
│   └── notification.ts          # 通知工具
│
├── stores/
│   ├── draft.ts                 # 草稿状态管理
│   ├── publish.ts               # 发布状态管理
│   ├── editor.ts                # 编辑器状态
│   └── ui.ts                    # UI状态（modals, menus等）
│
├── types/
│   ├── draft.ts
│   ├── publish.ts
│   ├── content.ts
│   ├── analytics.ts
│   └── common.ts
│
├── styles/
│   ├── index.css                # 全局样式
│   ├── variables.css            # CSS变量
│   ├── editor.css               # 编辑器样式
│   ├── layout.css               # 布局样式
│   └── components.css           # 组件样式
│
└── App.tsx                      # 主应用
```

### 核心组件详设

#### RichTextEditor.tsx
```tsx
interface RichTextEditorProps {
    value: string
    onChange: (content: string) => void
    onAutoSave?: (content: string) => void
    editorOptions?: {
        height?: number
        toolbar?: string[]
        modules?: any
    }
}

export function RichTextEditor({
    value,
    onChange,
    onAutoSave,
    editorOptions
}: RichTextEditorProps) {
    // 使用 react-quill 或 slate
    // 支持：
    // - 基础格式（加粗、斜体等）
    // - 标题
    // - 列表
    // - 代码块
    // - 图片上传
    // - 链接
    // - 视频嵌入
    // - 表格
    // - 撤销/重做
    // - 自动保存钩子
}
```

#### DraftAutoSave.tsx
```tsx
interface DraftAutoSaveProps {
    draftId: number
    content: string
    title: string
    interval?: number // 自动保存间隔（ms），默认30秒
    onSaving?: () => void
    onSaved?: (version: number) => void
    onError?: (error: Error) => void
}

export function DraftAutoSave({
    draftId,
    content,
    title,
    interval = 30000
}: DraftAutoSaveProps) {
    // 逻辑：
    // 1. 监听 content 和 title 变化
    // 2. 延迟 interval 毫秒后自动保存
    // 3. 显示保存状态（保存中... → 已保存）
    // 4. 如果保存失败，提示重试
    // 5. 定期和后端同步以处理并发编辑
}
```

#### PublishCalendar.tsx
```tsx
interface PublishCalendarProps {
    onDateSelect: (date: Date, tasks?: PublishTask[]) => void
    tasks: PublishTask[]
}

export function PublishCalendar({
    onDateSelect,
    tasks
}: PublishCalendarProps) {
    // 显示一个日历
    // 在有发布任务的日期标记
    // 点击日期显示该日期的任务列表
}
```

#### VersionHistory.tsx
```tsx
interface VersionHistoryProps {
    contentId: number
    currentVersion: number
    onVersionSelect: (version: number) => void
    onRevert: (version: number) => Promise<void>
}

export function VersionHistory({
    contentId,
    currentVersion,
    onVersionSelect,
    onRevert
}: VersionHistoryProps) {
    // 显示版本时间线
    // 支持点击查看某版本
    // 支持版本对比
    // 支持恢复到某版本
}
```

---

## 核心功能技术方案

### 1. 富文本编辑器

**选型方案**：
- 推荐：react-quill（简单易用，生态好）
- 高级：Slate（更灵活，可定制性强）

**实现清单**：
```
□ 集成编辑器库
□ 实现工具栏
□ 支持基础格式
□ 图片上传集成（OSS/S3）
□ 代码块和语法高亮
□ 撤销/重做
□ 自动保存逻辑（30秒或改动时）
□ 实时字数统计
```

**关键点**：
- 内容存储为HTML格式（兼容性好）
- 自动保存不阻塞编辑
- 支持本地草稿（localStorage）作为备份

---

### 2. 自动保存机制

**实现逻辑**：
```
用户编辑内容
    ↓
监听 onChange 事件，更新本地状态
    ↓
设置防抖定时器（300ms，收集所有改动）
    ↓
定时器触发后，发送 PUT /api/v1/drafts/{id}
    ↓
显示"保存中..." 状态
    ↓
保存成功 → 显示"已保存" + 时间戳
    ↓
保存失败 → 显示错误提示 + 重试按钮
    ↓
同时保存到 localStorage 作为本地备份
```

**后端处理**：
- 每次保存时自动创建新版本
- 版本号递增
- 记录 change_summary（可选，用户提供）

---

### 3. 版本历史和对比

**存储方案**：
- 每次保存自动创建 ContentVersion 记录
- 版本号：1, 2, 3...（自增）
- 保存完整的标题和内容副本

**对比算法**：
```go
// 使用 diff 算法（如 google/diff-match-patch）
// 计算两个版本之间的差异
// 返回：
// - 添加的部分（绿色高亮）
// - 删除的部分（红色高亮）
// - 修改的部分（黄色高亮）

type DiffResult struct {
    Type    string // added/removed/modified
    Content string
    Position string // "第3段" 或 行号
}
```

---

### 4. 发布管理

**定时发布实现**：

后端：
```go
// 使用定时任务库（如 robfig/cron）
// 每分钟检查一次是否有要发布的任务
// 任务流程：
// 1. 查询 status=pending 且 scheduled_at <= now 的任务
// 2. 更新状态为 publishing
// 3. 为每个 target_platform 创建发布子任务
// 4. 调用各平台 API 发布
// 5. 更新 platform_results
// 6. 发布失败时设置重试

type PublishTaskProcessor struct {
    db Database
    publishers map[string]PublisherInterface
    maxRetries int
}

func (p *PublishTaskProcessor) ProcessPendingTasks(ctx context.Context) error {
    // 查询待发布任务
    // 对于每个任务的每个平台：
    //   - 调用 publisher.Publish()
    //   - 更新结果
    //   - 失败则记录错误和重试信息
}
```

前端：
```tsx
// 显示日历
// 选择日期和时间
// 或点击"立即发布"按钮
// 发布后显示实时进度

// WebSocket 推送：
ws.on('message', (data) => {
    if (data.event === 'platform_published') {
        // 更新某平台的发布状态
        updatePlatformStatus(data.platform, 'published')
    }
})
```

---

### 5. 数据实时同步

**Web Socket 实时推送**：
```
客户端订阅发布任务的实时状态
    ↓
后端发布任务有进度更新时
    ↓
向所有订阅客户端推送消息
    ↓
客户端更新UI

消息格式：
{
    "type": "publish_status_update",
    "task_id": 200,
    "platform": "juejin",
    "status": "published",
    "post_url": "...",
    "metrics": {...}
}
```

**轮询备选方案**（如果不用WebSocket）：
```
GET /api/v1/publish-tasks/{id} 每2秒轮询一次
显示最新状态
```

---

## 开发规范和最佳实践

### 后端代码规范

```go
// 1. 错误处理
type ErrorCode string
const (
    ErrorCodeBadRequest ErrorCode = "BAD_REQUEST"
    ErrorCodeNotFound ErrorCode = "NOT_FOUND"
    ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
    ErrorCodeServerError ErrorCode = "SERVER_ERROR"
)

type ErrorResponse struct {
    Code    ErrorCode   `json:"code"`
    Message string      `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// 2. API 响应格式
type Response[T any] struct {
    Code    int         `json:"code"`         // 200, 400, 500等
    Message string      `json:"message"`
    Data    T           `json:"data,omitempty"`
    Error   *ErrorResponse `json:"error,omitempty"`
}

// 3. Repository 接口定义
type DraftRepository interface {
    Create(ctx context.Context, draft *Draft) error
    GetByID(ctx context.Context, id uint) (*Draft, error)
    List(ctx context.Context, filter DraftFilter) ([]*Draft, int64, error)
    Update(ctx context.Context, draft *Draft) error
    Delete(ctx context.Context, id uint) error
}

// 4. Service 业务逻辑
type DraftService struct {
    repo DraftRepository
    // 其他依赖
}

func (s *DraftService) SaveDraft(ctx context.Context, draft *Draft) error {
    // 验证
    if draft.Title == "" {
        return fmt.Errorf("title cannot be empty")
    }
    // 业务逻辑
    // 持久化
    return s.repo.Update(ctx, draft)
}

// 5. Handler HTTP处理
func (h *DraftHandler) Update(c *gin.Context) {
    var req UpdateDraftRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{
            Code: ErrorCodeBadRequest,
            Message: err.Error(),
        })
        return
    }

    draft, err := h.service.SaveDraft(c.Request.Context(), ...)
    if err != nil {
        c.JSON(500, ErrorResponse{...})
        return
    }

    c.JSON(200, Response{
        Code: 200,
        Data: draft,
    })
}
```

### 前端代码规范

```tsx
// 1. 组件Props类型定义
interface DraftEditorProps {
    draftId: number
    onSave?: (draft: Draft) => void
    onError?: (error: Error) => void
}

// 2. Hooks 使用
export function useDraft(draftId: number) {
    const [draft, setDraft] = useState<Draft | null>(null)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<Error | null>(null)

    useEffect(() => {
        fetchDraft()
    }, [draftId])

    const fetchDraft = async () => {
        setLoading(true)
        try {
            const data = await api.getDraft(draftId)
            setDraft(data)
        } catch (err) {
            setError(err as Error)
        } finally {
            setLoading(false)
        }
    }

    return { draft, loading, error, refetch: fetchDraft }
}

// 3. 错误边界
class ErrorBoundary extends React.Component<
    { children: React.ReactNode },
    { hasError: boolean; error: Error | null }
> {
    state = { hasError: false, error: null }

    static getDerivedStateFromError(error: Error) {
        return { hasError: true, error }
    }

    render() {
        if (this.state.hasError) {
            return <div>出错了: {this.state.error?.message}</div>
        }
        return this.props.children
    }
}

// 4. API 请求封装
async function api<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    url: string,
    data?: any
): Promise<T> {
    const response = await fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getToken()}`
        },
        body: data ? JSON.stringify(data) : undefined
    })

    if (!response.ok) {
        const error = await response.json()
        throw new Error(error.message)
    }

    return response.json()
}

// 5. 状态管理
type DraftState = {
    drafts: Draft[]
    currentDraft: Draft | null
    loading: boolean
    error: Error | null
}

const draftStore = create<DraftState & DraftActions>((set) => ({
    drafts: [],
    currentDraft: null,
    loading: false,
    error: null,

    fetchDrafts: async () => {
        set({ loading: true })
        try {
            const data = await api<Draft[]>('GET', '/api/v1/drafts')
            set({ drafts: data, error: null })
        } catch (error) {
            set({ error: error as Error })
        } finally {
            set({ loading: false })
        }
    }
}))
```

### 数据库优化

```sql
-- 索引规划
CREATE INDEX idx_drafts_user_id_edited ON drafts(user_id, last_edited_at DESC);
CREATE INDEX idx_drafts_user_status ON drafts(user_id, status);
CREATE INDEX idx_contents_user_published ON contents(user_id, published_at DESC);
CREATE INDEX idx_publish_tasks_scheduled ON publish_tasks(user_id, scheduled_at);
CREATE FULLTEXT INDEX idx_contents_search ON contents(title, body);

-- 查询优化
-- 避免 N+1：使用 GORM 的预加载
db.Preload("PublishTasks").Find(&contents)

-- 分页查询
offset := (page - 1) * limit
db.Limit(limit).Offset(offset).Find(&drafts)
```

---

## 下一步行动

本文档定义了：
1. ✅ 完整的数据库schema
2. ✅ 所有必需的API接口
3. ✅ 前端目录结构和核心组件
4. ✅ 技术实现方案

**立即行动**：
1. 在后端创建 models（Draft, Content, Version等）
2. 创建前端目录结构
3. 实现前端富文本编辑器组件
4. 实现后端API endpoints
5. 连接前后端

详见下一个文档：TASKS.md（具体的任务分配）
