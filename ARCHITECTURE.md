# AI Content Creation Engine - 架构设计文档

## 1. 项目概述

这是一个基于 Rust + React 的 AI 内容创作引擎，支持：
- **多内容类型**：文章、代码、图片、音频、视频脚本等
- **多AI提供商**：OpenAI、Ollama、通义千问等
- **自动化发布**：自动发布到微信公众号、CSDN、掘金等平台

## 2. 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户界面 (React)                        │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │
│  │ 内容生成器   │ │ 内容编辑器   │ │ 发布管理器   │
│  │ (Generator) │ │ (Editor) │ │ (Publisher) │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
│         │               │               │
│         └───────────────┴───────────────┘
└─────────┴───────────────┴───────────────┘
          │               │               │
          ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      业务服务层 (Rust)                         │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │
│  │ 内容服务     │ │ AI 服务      │ │ 发布服务     │
│  │ Content      │ │ AI Service   │ │ Publish      │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
│         │               │               │
│  ┌──────┴─────────┴──────┐               ▼                       │
│  │  内容类型处理器       │                     │
│  │ Text/Code/Img/Audio  │
│  └──────┬────────────┬───────┘
│         │            │
│  ┌──────┴─────┐  │  ┌────┴─────────────────┐
│  │  提供商适配器 │  │  │   提供商适配器      │
│  │ OpenAI      │  │  │  Ollama       │
│  └──────┬───────┘  │  └──────┬─────────┘
│         │          │         │
│         └──────────┴─────────┘
└─────────────────────────────────────────────────────────────────┘
          │               │
          ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      数据层                                                   │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │
│  │ 内容存储     │ │ AI 配置     │ │ 发布记录   │
│  │ Content      │ │ AI Config    │ │ Publish Logs │
│  └──────────────┘ └──────────────┘ └──────────────┘
└─────────────────────────────────────────────────────────────────┘
```

## 3. 模块设计

### 3.1 核心模块

| 模块 | 职责 | 主要组件 |
|------|------|----------|
| **内容生成** | 生成各类内容 | TextGenerator, CodeGenerator, ImageGenerator |
| **AI 适配** | 多AI提供商适配 | OpenAIAdapter, OllamaAdapter |
| **内容发布** | 发布到平台 | WeChatPublisher, CSDNPublisher |
| **内容管理** | 内容 CRUD | ContentService |
| **任务调度** | 定时任务 | Scheduler |

### 3.2 数据模型

```rust
// 内容类型枚举
enum ContentType {
    Text,          // 文本
    Code,          // 代码
    Image,         // 图片
    Audio,         // 音频
    Video,         // 视频
    Mixed,         // 混合
}

// AI 提供商枚举
enum AIProvider {
    OpenAI,
    Ollama,
    TongyiQianwen,
    Local,
    Custom(String),
}

// 平台枚举
enum PublishPlatform {
    WeChatOfficialAccount,
    WeChatMoments,
    CSDN,
    JueJin,
    ZhiHu,
    Jianshu,
    Custom(String),
}

// 内容实体
struct Content {
    id: Uuid,
    title: String,
    content: String,
    content_type: ContentType,
    author: String,
    status: ContentStatus,
    tags: Vec<String>,
    metadata: HashMap<String, String>,
    created_at: DateTime<Utc>,
    updated_at: DateTime<Utc>,
}

// 发布任务实体
struct PublishTask {
    id: Uuid,
    content_id: Uuid,
    platform: PublishPlatform,
    status: PublishStatus,
    settings: HashMap<String, String>,
    created_at: DateTime<Utc>,
    published_at: Option<DateTime<Utc>>,
}

// AI 配置实体
struct AIConfig {
    id: Uuid,
    provider: AIProvider,
    api_key: String,
    model: String,
    base_url: String,
    settings: HashMap<String, String>,
    created_at: DateTime<Utc>,
    updated_at: DateTime<Utc>,
    is_default: bool,
}
```

### 3.3 核心接口

```rust
// AI 适配器接口
trait AIAdapter: Send + Sync {
    async fn generate_text(&self, prompt: &str, options: GenerateOptions) -> Result<String>;
    async fn generate_image(&self, prompt: &str, options: GenerateOptions) -> Result<Vec<u8>>;
    async fn generate_audio(&self, prompt: &str, options: GenerateOptions) -> Result<Vec<u8>>;
    async fn generate_code(&self, prompt: &str, options: GenerateOptions) -> Result<String>;
}

// 发布适配器接口
trait Publisher: Send + Sync {
    async fn publish(&self, content: &Content, settings: &HashMap<String, String>) -> Result<String>;
    async fn draft(&self, content: &Content, settings: &HashMap<String, String>) -> Result<String>;
    async fn delete(&self, content: &Content) -> Result<()>;
    async fn get_user_info(&self) -> Result<UserInfo>;
}

// 内容处理器接口
trait ContentHandler: Send + Sync {
    fn can_handle(&self, content_type: &ContentType) -> bool;
    async fn process(&self, content: Content) -> Result<Content>;
}
```

## 4. 后端 API 设计

### 4.1 内容相关 API

```
GET    /api/v1/contents              - 获取内容列表
GET    /api/v1/contents/{id}          - 获取内容详情
POST   /api/v1/contents              - 创建内容
PUT    /api/v1/contents/{id}          - 更新内容
DELETE /api/v1/contents/{id}          - 删除内容
POST   /api/v1/contents/{id}/clone   - 克隆内容
GET    /api/v1/contents/{id}/publish - 获取内容发布记录
```

### 4.2 AI 服务 API

```
POST   /api/v1/ai/generate          - 生成内容
POST   /api/v1/ai/config             - AI 配置管理
GET    /api/v1/ai/config             - 获取 AI 配置
```

### 4.3 发布服务 API

```
POST   /api/v1/publish             - 发布内容
GET    /api/v1/publish/task/{id}   - 获取发布任务状态
GET    /api/v1/publish/platforms   - 获取支持的平台列表
```

### 4.4 任务调度 API

```
GET    /api/v1/scheduler/tasks     - 获取任务列表
POST   /api/v1/scheduler/tasks     - 创建定时任务
PUT    /api/v1/scheduler/tasks/{id} - 更新任务
DELETE /api/v1/scheduler/tasks/{id} - 删除任务
```

## 5. 技术栈

| 层级 | 技术 |
|------|------|
| **后端** | Rust, Actix Web / Tide |
| **前端** | React, TypeScript, Vite, Ant Design |
| **数据库** | PostgreSQL / SQLite |
| **AI** | OpenAI API, Ollama API, Tongyi API |
| **部署** | Docker, Docker Compose |
| **CI/CD** | GitHub Actions |

## 6. 部署架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              用户                                   │
└─────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           反向代理 (Nginx)                           │
└─────────────────────────────────────────────────────────────────────────┘
                                     │
        ┌──────────────────────────────┼──────────────────────────────────────┐
        │                              │                              │
        ▼                              ▼                              ▼
┌───────────────────────┐    ┌───────────────────────┐    ┌───────────────────────┐
│                       │    │                       │    │                       │
│                      Docker                         │    │                      Docker                         │
│                    ┌───────────────┐  │    │                     ┌───────────────┐  │
│                    │  Web App     │  │    │                     │   DB          │
│                    │  (Rust)     │  │    │                     │  (PostgreSQL) │
│                    └───────────────┘  │    │                     └───────────────┘  │
│                                      │    │                                          │
│                    ┌───────────────┐  │    │                     ┌───────────────┐  │
│                    │  Redis      │  │    │                     │  Volume      │
│                    │  (Cache)    │  │    │                     │  (Data)     │
│                    └───────────────┘  │    │                     └───────────────┘  │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## 7. 关键特性

### 7.1 内容生成
- 多内容类型支持
- 多AI提供商适配
- 内容模板管理
- 内容版本控制

### 7.2 自动化发布
- 多平台发布支持
- 定时发布
- 发布队列
- 发布重试

### 7.3 内容管理
- 内容分类管理
- 标签管理
- 内容审核
- 内容统计

## 8. 安全设计

- API 鉴权 (JWT)
- 数据加密存储
- SQL 注入防护
- XSS 防护
- CSRF 防护

## 9. 扩展性设计

- 插件化架构
- 适配器模式
- 事件驱动
- 消息队列

## 10. 未来规划

- 移动端 App
- 更多 AI 模型支持
- 更多发布平台支持
- 数据分析与统计
- 团队协作
