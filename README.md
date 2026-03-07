# AI Content Engine

一个基于 Go + Gin (后端) + React (前端) 的 AI 内容创作引擎。

## 功能特性

- ✅ 多内容类型：文章、代码、图片、音频等
- ✅ 多AI提供商：OpenAI、Ollama、通义千问等
- ✅ 自动化发布：自动发布到微信公众号、CSDN、掘金等平台
- ✅ 内容管理：内容编辑、发布管理、版本控制
- ✅ 定时发布：支持定时发布内容到指定平台
- ✅ 多平台支持：微信公众号、CSDN、掘金、知乎、简书等

## 技术栈

### 后端
- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Go-Redis](https://github.com/redis/go-redis) - Redis 客户端
- [validator](https://github.com/go-playground/validator) - 数据验证
- [viper](https://github.com/spf13/viper) - 配置管理
- [Docker](https://www.docker.com/) - 容器化

### 前端
- [React](https://react.dev/) - 前端框架
- [TypeScript](https://www.typescriptlang.org/) - 类型安全
- [Vite](https://vitejs.dev/) - 构建工具
- [Ant Design](https://ant.design/) - UI 组件库

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+ / Bun
- PostgreSQL 14+ / MySQL 8+ / SQLite
- Docker 24+

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/your-username/ai-content-engine.git
cd ai-content-engine

# 后端开发
cd backend
go mod tidy
go run cmd/server/main.go --config config.yaml

# 前端开发
cd ../frontend
npm install
npm run dev
```

### 使用 Docker

```bash
# 构建镜像
docker build -t ai-content-engine .
docker-compose up -d

# 访问应用
open http://localhost:3000
```

## 项目结构

```
ai-content-engine/
├── backend/              # Go 后端
│   ├── cmd/
│   │   └── server/      # 主程序入口
│   │       └── main.go
│   ├── internal/
│   │   ├── config/      # 配置管理
│   │   ├── handler/     # HTTP 处理器
│   │   ├── middleware/  # 中间件
│   │   ├── model/       # 数据模型
│   │   ├── repository/  # 数据访问层
│   │   ├── service/     # 业务逻辑
│   │   └── router/      # 路由配置
│   ├── go.mod           # Go 依赖
│   ├── go.sum
│   └── Dockerfile       # Docker 构建文件
├── frontend/             # React 前端
│   ├── src/
│   │   ├── App.tsx       # 应用入口
│   │   ├── components/   # 可复用组件
│   │   ├── pages/        # 页面组件
│   │   ├── hooks/        # 自定义 Hook
│   │   └── utils/        # 工具函数
│   ├── index.html
│   ├── vite.config.ts   # Vite 配置
│   └── package.json     # npm 依赖
├── docker-compose.yml   # Docker Compose 配置
├── docker-compose.debug.yml # 调试用配置
├── docker-compose.prod.yml # 生产用配置
├── .gitignore
├── .golangci.yml        # Go linter 配置
├── .prettierrc          # Prettier 配置
├── .eslintrc.js         # ESLint 配置
├── .dockerignore
├── ARCHITECTURE.md      # 架构文档
├── DOMAIN.md            # 领域层文档
└── README.md            # 项目说明
```

## API 文档

```bash
# 访问 API 文档
curl -X GET http://localhost:8080/api/v1/contents
curl -X POST http://localhost:8080/api/v1/contents
curl -X GET http://localhost:8080/api/v1/contents/{id}
curl -X PUT http://localhost:8080/api/v1/contents/{id}
curl -X DELETE http://localhost:8080/api/v1/contents/{id}
curl -X GET http://localhost:8080/api/v1/ai/generate
curl -X POST http://localhost:8080/api/v1/publish
```

## 配置说明

### YAML 配置文件（必需）

```yaml
# backend/config.yaml
server:
  env: dev
  host: 0.0.0.0
  port: 8888

database:
  driver: postgres  # postgres / mysql / sqlite
  host: localhost
  user: root
  pass: pass
  db_name: ai_content_engine
  ssl_mode: disable

redis:
  host: localhost
  port: 6379
  password: ""

llm:
  base_url: http://localhost:11434
  model: qwen2
  timeout: 30
  max_tokens: 2048

text_embedding:
  base_url: http://localhost:11434
```

启动命令（必须传 `--config`）：

```bash
cd backend
go run cmd/server/main.go --config ./config.yaml
```

## 安全建议

- 🔒 使用 HTTPS 加密传输
- 🔒 定期更新依赖包
- 🔒 启用 CORS 策略
- 🔒 验证用户输入 (使用 validator)
- 🔒 使用 JWT 认证
- 🔒 审计日志记录
- 🔒 敏感信息加密
- 🔒 权限控制
- 🔒 速率限制

## 常见问题

**Q: 如何添加新的 AI 提供商？**

A: 在 `backend/internal/adapters/` 目录下创建新的适配器实现 `AIAdapter` 接口

**Q: 如何添加新的发布平台？**

A: 在 `backend/internal/adapters/publishers/` 目录下创建新的发布器实现 `Publisher` 接口

**Q: 如何配置多 AI 提供商？**

A: 在 `backend/config.yaml` 里修改 `llm` 相关配置，并通过 `--config` 指定对应配置文件

**Q: 如何配置定时发布？**

A: 使用 `SchedulerService` 服务配置定时发布任务

**Q: 如何配置内容审核？**

A: 使用 `ContentService` 服务配置内容审核规则

**Q: 如何配置内容版本控制？**

A: 使用 `ContentService` 服务配置版本控制策略

**Q: 如何配置发布队列？**

A: 使用 `PublisherService` 服务配置发布队列策略

**Q: 如何配置发布重试？**

A: 使用 `PublisherService` 服务配置重试策略

**Q: 如何配置发布通知？**

A: 使用 `PublisherService` 服务配置通知策略

## 贡献指南

欢迎提交 Issue 和 PR！

## 许可证

MIT License

## 联系方式

- GitHub: https://github.com/your-username/ai-content-engine
- Email: your-email@example.com
- Twitter: @your-twitter
