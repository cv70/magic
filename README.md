# AI Content Engine

一个基于 Rust + Axum (后端) + React (前端) 的 AI 内容创作引擎。

## 功能特性

- ✅ 多内容类型：文章、代码、图片、音频等
- ✅ 多AI提供商：OpenAI、Ollama、通义千问等
- ✅ 自动化发布：自动发布到微信公众号、CSDN、掘金等平台
- ✅ 内容管理：内容编辑、发布管理、版本控制
- ✅ 定时发布：支持定时发布内容到指定平台
- ✅ 多平台支持：微信公众号、CSDN、掘金、知乎、简书等

## 技术栈

### 后端
- [Axum](https://github.com/tokio-rs/axum) - 高性能 Rust Web 框架
- [Tokio](https://tokio.rs/) - 异步运行时
- [SQLx](https://github.com/launchbadge/sqlx) - 异步数据库驱动
- [Redis](https://redis.io/) - 缓存服务
- [Rustls](https://github.com/rustls/rustls) - TLS 库
- [OpenAPI](https://www.openapis.org/) - API 文档
- [Docker](https://www.docker.com/) - 容器化

### 前端
- [React](https://react.dev/) - 前端框架
- [TypeScript](https://www.typescriptlang.org/) - 类型安全
- [Vite](https://vitejs.dev/) - 构建工具
- [Ant Design](https://ant.design/) - UI 组件库

## 快速开始

### 环境要求

- Rust 1.74+ (2024 Edition)
- Node.js 18+ / Bun
- PostgreSQL 14+
- Docker 24+

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/your-username/ai-content-engine.git
cd ai-content-engine

# 安装依赖
cd backend
cargo build --release
cargo run --release -- --config ./config.yaml

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
├── backend/              # Rust 后端
│   ├── src/
│   │   ├── main.rs       # 主程序入口
│   │   ├── api/          # API 路由
│   │   ├── models/       # 数据模型
│   │   ├── services/     # 业务逻辑
│   │   ├── adapters/     # 外部服务适配器
│   │   └── config.rs     # 配置管理
│   ├── Cargo.toml       # Rust 依赖
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
├── .gitignore         # Git 忽略文件
├── .rustfmt.toml      # Rust 格式化配置
├── .prettierrc        # Prettier 配置
├── .eslintrc.js       # ESLint 配置
├── .dockerignore      # Docker 忽略文件
├── .gitignore         # Git 忽略文件
├── ARCHITECTURE.md      # 架构文档
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

```bash
# backend/config.yaml
server:
  env: dev
  host: 0.0.0.0
  port: 8888

database:
  host: localhost
  user: root
  pass: pass
  db_name: vc_backend

redis:
  host: localhost
  port: 6379
  pass: null

vector:
  host: localhost
  api_key: ""

scylla:
  host: localhost
  user: scylla
  pass: scylla

llm:
  base_url: http://localhost:11434
  model: qwen2
  timeout: 30
  max_tokens: 2048

text_embedding:
  base_url: http://localhost:11434

browser:
  url: http://localhost:9222
```

启动命令（必须传 `--config`）：

```bash
cd backend
cargo run --release -- --config ./config.yaml
```

## 安全建议

- 🔒 使用 HTTPS 加密传输
- 🔒 定期更新依赖包
- 🔒 启用 CORS 策略
- 🔒 验证用户输入
- 🔒 使用 JWT 认证
- 🔒 审计日志记录
- 🔒 敏感信息加密
- 🔒 权限控制
- 🔒 速率限制

## 常见问题

**Q: 如何添加新的 AI 提供商？

A: 在 `backend/src/adapters/` 目录下创建新的适配器实现 `AIAdapter` trait

**Q: 如何添加新的发布平台？

A: 在 `backend/src/adapters/publishers/` 目录下创建新的发布器实现 `Publisher` trait

**Q: 如何配置多 AI 提供商？

A: 在 `backend/config.yaml` 里修改 `llm` 相关配置，并通过 `--config` 指定对应配置文件

**Q: 如何配置定时发布？

A: 使用 `AIScheduler` 服务配置定时发布任务

**Q: 如何配置内容审核？

A: 使用 `ContentService` 服务配置内容审核规则

**Q: 如何配置内容版本控制？

A: 使用 `ContentService` 服务配置版本控制策略

**Q: 如何配置发布队列？

A: 使用 `PublisherService` 服务配置发布队列策略

**Q: 如何配置发布重试？

A: 使用 `PublisherService` 服务配置重试策略

**Q: 如何配置发布通知？

A: 使用 `PublisherService` 服务配置通知策略

## 贡献指南

欢迎提交 Issue 和 PR！

## 许可证

MIT License

## 联系方式

- GitHub: https://github.com/your-username/ai-content-engine
- Email: your-email@example.com
- Twitter: @your-twitter
