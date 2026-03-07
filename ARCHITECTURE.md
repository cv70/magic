# AI Content Creation Engine - 架构设计文档

## 1. 项目概述

这是一个基于 Go + Gin 的 AI 内容创作引擎，支持：
- **多内容类型**：文章、代码、图片、音频、视频脚本等
- **多AI提供商**：OpenAI、Ollama、通义千问等
- **自动化发布**：自动发布到微信公众号、CSDN、掘金等平台

## 2. 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户界面 (React)                        │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │
│  │ 内容生成器   │ │ 内容编辑器   │ │ 发布管理器   │
│  │ (Generator) │ │ (Editor)     │ │ (Publisher) │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
│         │               │               │
│         └───────────────┴───────────────┘
└─────────┴───────────────┴───────────────┘
          │               │               │
          ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      业务服务层 (Go + Gin)                     │
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
│  │ OpenAI      │  │  │  Ollama             │
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

```go
// ContentType 内容类型枚举
type ContentType string

const (
    ContentTypeText   ContentType = "text"   // 文本
    ContentTypeCode   ContentType = "code"   // 代码
    ContentTypeImage  ContentType = "image"  // 图片
    ContentTypeAudio  ContentType = "audio"  // 音频
    ContentTypeVideo  ContentType = "video"  // 视频
    ContentTypeMixed  ContentType = "mixed"  // 混合
)

// AIProvider AI 提供商枚举
type AIProvider string

const (
    ProviderOpenAI       AIProvider = "openai"
    ProviderOllama       AIProvider = "ollama"
    ProviderTongyiQianwen AIProvider = "tongyi_qianwen"
    ProviderLocal        AIProvider = "local"
    ProviderCustom       AIProvider = "custom"
)

// PublishPlatform 平台枚举
type PublishPlatform string

const (
    PlatformWeChatOfficialAccount PublishPlatform = "wechat_official_account"
    PlatformWeChatMoments        PublishPlatform = "wechat_moments"
    PlatformCSDN                 PublishPlatform = "csdn"
    PlatformJueJin              PublishPlatform = "juejin"
    PlatformZhiHu               PublishPlatform = "zhihu"
    PlatformJianshu             PublishPlatform = "jianshu"
)

// Content 内容实体
type Content struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"size:255" json:"title"`
    Content     string    `gorm:"type:text" json:"content"`
    ContentType ContentType `gorm:"size:20" json:"content_type"`
    Author      string    `gorm:"size:100" json:"author"`
    Status      string    `gorm:"size:20" json:"status"`
    Tags        string    `gorm:"size:500" json:"tags"` // JSON array stored as string
    Metadata    string    `gorm:"type:text" json:"metadata"` // JSON stored as string
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// PublishTask 发布任务实体
type PublishTask struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    ContentID   uint           `gorm:"index" json:"content_id"`
    Platform    PublishPlatform `gorm:"size:50" json:"platform"`
    Status      string         `gorm:"size:20" json:"status"`
    Settings    string         `gorm:"type:text" json:"settings"` // JSON stored as string
    CreatedAt   time.Time      `json:"created_at"`
    PublishedAt *time.Time     `json:"published_at,omitempty"`
}

// AIConfig AI 配置实体
type AIConfig struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Provider  AIProvider `gorm:"size:50" json:"provider"`
    APIKey    string    `gorm:"size:500" json:"api_key"`
    Model     string    `gorm:"size:100" json:"model"`
    BaseURL   string    `gorm:"size:255" json:"base_url"`
    Settings  string    `gorm:"type:text" json:"settings"` // JSON stored as string
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsDefault bool      `gorm:"default:false" json:"is_default"`
}
```

### 3.3 核心接口

```go
// AIAdapter AI 适配器接口
type AIAdapter interface {
    GenerateText(ctx context.Context, prompt string, options GenerateOptions) (string, error)
    GenerateImage(ctx context.Context, prompt string, options GenerateOptions) ([]byte, error)
    GenerateAudio(ctx context.Context, prompt string, options GenerateOptions) ([]byte, error)
    GenerateCode(ctx context.Context, prompt string, options GenerateOptions) (string, error)
}

// Publisher 发布适配器接口
type Publisher interface {
    Publish(ctx context.Context, content *Content, settings map[string]string) (string, error)
    Draft(ctx context.Context, content *Content, settings map[string]string) (string, error)
    Delete(ctx context.Context, content *Content) error
    GetUserInfo(ctx context.Context) (*UserInfo, error)
}

// ContentHandler 内容处理器接口
type ContentHandler interface {
    CanHandle(contentType ContentType) bool
    Process(ctx context.Context, content Content) (Content, error)
}

// GenerateOptions 生成选项
type GenerateOptions struct {
    Model      string            `json:"model"`
    MaxTokens  int               `json:"max_tokens"`
    Temperature float64          `json:"temperature"`
    TopP       float64           `json:"top_p"`
    Extra      map[string]string `json:"extra"`
}

// UserInfo 用户信息
type UserInfo struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Token string `json:"token"`
}
```

## 4. 后端 API 设计

### 4.1 内容相关 API

```
GET    /api/v1/contents              - 获取内容列表
GET    /api/v1/contents/{id}         - 获取内容详情
POST   /api/v1/contents              - 创建内容
PUT    /api/v1/contents/{id}        - 更新内容
DELETE /api/v1/contents/{id}        - 删除内容
POST   /api/v1/contents/{id}/clone   - 克隆内容
GET    /api/v1/contents/{id}/publish - 获取内容发布记录
```

### 4.2 AI 服务 API

```
POST   /api/v1/ai/generate          - 生成内容
POST   /api/v1/ai/config            - AI 配置管理
GET    /api/v1/ai/config            - 获取 AI 配置
```

### 4.3 发布服务 API

```
POST   /api/v1/publish               - 发布内容
GET    /api/v1/publish/task/{id}     - 获取发布任务状态
GET    /api/v1/publish/platforms     - 获取支持的平台列表
```

### 4.4 任务调度 API

```
GET    /api/v1/scheduler/tasks       - 获取任务列表
POST   /api/v1/scheduler/tasks       - 创建定时任务
PUT    /api/v1/scheduler/tasks/{id}  - 更新任务
DELETE /api/v1/scheduler/tasks/{id}  - 删除任务
```

## 5. 技术栈

| 层级 | 技术 |
|------|------|
| **后端** | Go 1.21+, Gin Web Framework |
| **数据库** | PostgreSQL / MySQL / SQLite, GORM |
| **缓存** | Redis |
| **AI** | OpenAI API, Ollama API, 通义千问 API |
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
│                    │  (Go + Gin) │  │    │                     │  (PostgreSQL) │
│                    └───────────────┘  │    │                     └───────────────┘  │
│                                      │    │                                          │
│                    ┌───────────────┐  │    │                     ┌───────────────┐  │
│                    │  Redis       │  │    │                     │  Volume      │
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
- SQL 注入防护 (GORM 参数化查询)
- XSS 防护
- CSRF 防护

## 9. 扩展性设计

- 插件化架构
- 适配器模式
- 事件驱动
- 消息队列 (RabbitMQ / Kafka)

## 10. 未来规划

- 移动端 App
- 更多 AI 模型支持
- 更多发布平台支持
- 数据分析与统计
- 团队协作
