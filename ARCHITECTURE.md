# AI Content Engine - 系统架构设计

## 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          用户浏览器 (Frontend)                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │  React 18 + TypeScript + Ant Design                             │   │
│  │  ┌──────────────────────────────────────────────────────────┐   │   │
│  │  │  📄 Pages                    🎯 Components               │   │   │
│  │  ├──────────────────────────────┼──────────────────────────┤   │   │
│  │  │  • Dashboard                 │ • DomainSelector (领域卡)│   │   │
│  │  │  • DraftCenter (✨ 领域筛选) │ • DomainParamsForm (参) │   │   │
│  │  │  • Editor (编辑器)           │ • SquareCard (广场卡)   │   │   │
│  │  │  • AIStudio (✨ 4步向导)    │ • TipTapEditor          │   │   │
│  │  │  • PublishManager (发布)     │ • AutoSaveIndicator     │   │   │
│  │  │  • Analytics (数据分析)      │                         │   │   │
│  │  │  • Square (✨ 内容广场)      │                         │   │   │
│  │  │  • Settings (设置)           │                         │   │   │
│  │  └──────────────────────────────┴──────────────────────────┘   │   │
│  │                                                                  │   │
│  │  🌍 Constants (配置数据)          🔗 API Layer                 │   │
│  │  ├─ domains.ts (9个领域)        ├─ authApi                    │   │
│  │  └─ domainTemplates.ts (20模板) ├─ draftApi                   │   │
│  │                                  ├─ contentApi                 │   │
│  │  📊 State Management              ├─ publisherApi              │   │
│  │  ├─ React Hooks                   ├─ aiApi                     │   │
│  │  ├─ localStorage                  ├─ squareApi (✨ 新增)       │   │
│  │  └─ Zustand (认证/主题)          └─ analyticsApi              │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    ▼
                    ┌────────────────────────────────────┐
                    │   HTTP/HTTPS (Bearer Token Auth)   │
                    │  Content-Type: application/json    │
                    └────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        后端服务器 (Backend)                             │
│                      Go + Gin + GORM + MySQL                          │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  🔐 认证层 (Auth Middleware)                                    │  │
│  │  └─ JWT Token Validation (Bearer Token)                         │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  🎯 Domain Layer (DDD Architecture)                             │  │
│  │  ├─ Identity Domain                                             │  │
│  │  │  ├─ login / register / me                                    │  │
│  │  │  └─ user management                                          │  │
│  │  │                                                              │  │
│  │  ├─ Content Domain                                              │  │
│  │  │  ├─ drafts (CRUD)                                            │  │
│  │  │  ├─ content (CRUD)                                           │  │
│  │  │  ├─ versions (历史)                                          │  │
│  │  │  └─ tags & categories                                        │  │
│  │  │                                                              │  │
│  │  ├─ AI Generation Domain                                        │  │
│  │  │  ├─ generators                                               │  │
│  │  │  ├─ generate (async task)                                    │  │
│  │  │  └─ task polling                                             │  │
│  │  │                                                              │  │
│  │  ├─ Publishing Domain                                           │  │
│  │  │  ├─ publishers                                               │  │
│  │  │  ├─ publish tasks                                            │  │
│  │  │  └─ analytics                                                │  │
│  │  │                                                              │  │
│  │  ├─ Square Domain (✨ 新增)                                    │  │
│  │  │  ├─ posts (列表、详情、发布)                                │  │
│  │  │  ├─ likes & collects                                         │  │
│  │  │  ├─ comments                                                 │  │
│  │  │  └─ filters & sorting                                        │  │
│  │  │                                                              │  │
│  │  ├─ Scheduling Domain                                           │  │
│  │  │  └─ scheduled tasks                                          │  │
│  │  │                                                              │  │
│  │  └─ Configuration Domain                                        │  │
│  │     └─ system & provider configs                                │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  📊 Repository Layer (Data Access)                              │  │
│  │  ├─ SquarePostRepository (新增)                                │  │
│  │  ├─ SquareCommentRepository (新增)                             │  │
│  │  ├─ SquareLikeRepository (新增)                                │  │
│  │  ├─ SquareCollectRepository (新增)                             │  │
│  │  └─ [Other repositories...]                                     │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  🗄️ Database Layer (GORM ORM)                                  │  │
│  │  └─ MySQL Connection Pooling                                    │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         数据库 (MySQL 8.0)                             │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │  📋 Existing Tables (Previous Implementation)                   │  │
│  │  ├─ users (身份管理)                                            │  │
│  │  ├─ drafts (草稿，✨ 新增 domain 字段)                        │  │
│  │  ├─ content (发布内容)                                          │  │
│  │  ├─ content_versions (版本历史)                                 │  │
│  │  ├─ tags & categories (标签和分类)                              │  │
│  │  ├─ publishers & publish_tasks (发布管理)                       │  │
│  │  └─ [Other tables...]                                           │  │
│  │                                                                  │  │
│  │  ✨ New Tables (Square Implementation)                         │  │
│  │  ├─ square_posts                                                │  │
│  │  │  ├─ PK: id                                                   │  │
│  │  │  ├─ FK: draft_id, user_id                                    │  │
│  │  │  ├─ Indexes: domain, created_at, likes_count, ...           │  │
│  │  │  └─ Columns: title, content, preview_text, tags, ...        │  │
│  │  │                                                              │  │
│  │  ├─ square_comments                                             │  │
│  │  │  ├─ PK: id                                                   │  │
│  │  │  ├─ FK: post_id, user_id                                     │  │
│  │  │  └─ Indexes: post_id, created_at                             │  │
│  │  │                                                              │  │
│  │  ├─ square_likes                                                │  │
│  │  │  ├─ PK: id                                                   │  │
│  │  │  ├─ FK: post_id, user_id                                     │  │
│  │  │  └─ UNIQUE(post_id, user_id) - 防重复点赞                  │  │
│  │  │                                                              │  │
│  │  └─ square_collects                                             │  │
│  │     ├─ PK: id                                                   │  │
│  │     ├─ FK: post_id, user_id                                     │  │
│  │     └─ UNIQUE(post_id, user_id) - 防重复收藏                  │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 数据流图

### 1. 内容广场浏览流程

```
User 访问 /square
        ↓
前端加载 Square 页面组件
        ↓
调用 squareApi.list(domain, keyword, sort, page)
        ↓
发送 POST /api/v1/square/posts (JWT Token)
        ↓
后端 ApiListSquarePosts 处理器
  ├─ 从 JWT 提取 userID
  ├─ 验证分页参数
  ├─ 构建查询条件（domain filter, keyword search, sort）
  ├─ 从 square_posts 表查询
  ├─ 为每条帖子添加用户点赞/收藏状态
  └─ 返回 SquarePostDTO 数组
        ↓
前端接收数据，渲染 SquareCard 组件网格
        ↓
用户交互（滚动、过滤、排序）
```

### 2. 点赞流程

```
用户点击心形按钮
        ↓
前端调用 squareApi.like(postId)
        ↓
发送 POST /api/v1/square/like { post_id: 1 }
        ↓
后端 ApiLike 处理器
  ├─ 从 JWT 提取 userID
  ├─ 在 square_likes 表插入 (post_id, user_id)
  │  └─ 如果已存在，返回 "already liked" 错误
  ├─ square_posts.likes_count += 1
  └─ 返回 200 OK
        ↓
前端：
  ├─ 更新本地状态 is_liked = true
  ├─ 更新计数 likes_count += 1
  └─ 重新渲染卡片
```

### 3. 发布草稿到广场流程

```
用户在编辑器完成编辑
        ↓
前端在 Editor 组件添加"发布到广场"按钮
        ↓
用户点击按钮
        ↓
前端调用 squareApi.publish(draftId)
        ↓
发送 POST /api/v1/square/publish { draft_id: 123 }
        ↓
后端 ApiPublishToSquare 处理器
  ├─ 从 JWT 提取 userID
  ├─ 从 content domain 获取草稿
  │  ├─ 验证草稿归属 (draft.user_id == userID)
  │  └─ 提取 domain 字段（从 metadata）
  ├─ 创建 square_posts 记录
  │  ├─ title: draft.title
  │  ├─ preview_text: draft.content[:200]
  │  ├─ domain: draft.metadata["domain"]
  │  └─ user_id: userID
  └─ 返回 SquarePostDTO
        ↓
前端：
  ├─ 显示成功提示
  ├─ 导航到 /square 页面
  └─ 广场列表显示新帖子
```

---

## API 调用时序图

```
Frontend                          Backend                      Database
   │                               │                              │
   ├─ POST /api/v1/square/posts ─>│                              │
   │  {domain, keyword, sort}     │                              │
   │                              ├─ Validate JWT token          │
   │                              │                              │
   │                              ├─ Build query ────────────>│
   │                              │   SELECT * FROM             │
   │                              │   square_posts WHERE         │
   │                              │   domain = ?                 │
   │                              │   AND title LIKE ?           │
   │                              │   ORDER BY created_at DESC   │
   │                              │<─────────── Return rows ──┤
   │                              │                              │
   │                              ├─ For each post:             │
   │                              │   Check if liked/collected   │
   │                              │   ────────────────────>│
   │                              │   SELECT COUNT(*) FROM       │
   │                              │   square_likes WHERE         │
   │                              │   post_id = ? AND user_id= ? │
   │                              │<─────────────────────┤
   │                              │                              │
   │  <─ 200 {posts, total} ──────┤                              │
   │    [SquarePostDTO[]]         │                              │
   │                              │                              │
   ├─ POST /api/v1/square/like ──>│                              │
   │  {post_id: 1}               │                              │
   │                              ├─ Insert like ───────────>│
   │                              │   INSERT INTO               │
   │                              │   square_likes VALUES        │
   │                              │   (1, user_id)               │
   │                              │<────────────────────┤
   │                              │                              │
   │                              ├─ Update count ─────────>│
   │                              │   UPDATE square_posts        │
   │                              │   SET likes_count = ...+1    │
   │                              │<────────────────────┤
   │  <─ 200 {message: "liked"} ──┤                              │
   │                              │                              │
```

---

## 9 个领域分类系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                    领域分类系统 (Domains)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  🎬 film-drama (影视与短剧制作)                               │
│  ├─ Icon: 🎬                                                    │
│  ├─ SubCategories: [剧本, 分镜, 台词, 场景描述]               │
│  └─ Templates:                                                 │
│     ├─ 短剧剧本生成 (base_prompt template)                     │
│     └─ 精准对话生成 (base_prompt template)                     │
│                                                                 │
│  🎭 ai-comics (AI漫剧/动态漫)                                │
│  ├─ Icon: 🎭                                                    │
│  ├─ SubCategories: [故事板, 对话气泡, 人物设定]                │
│  └─ Templates: [漫画故事板生成]                                │
│                                                                 │
│  📱 micro-drama (微短剧)                                        │
│  ├─ Icon: 📱                                                    │
│  └─ Templates: [竖屏脚本生成]                                  │
│                                                                 │
│  🎞️ animation (动画电影/长视频)                               │
│  └─ Templates: [角色设定生成]                                  │
│                                                                 │
│  📣 marketing (营销与电商广告)                                │
│  └─ Templates: [产品文案生成]                                  │
│                                                                 │
│  📚 education (教育与知识出版)                                │
│  └─ Templates: [课程讲义生成]                                  │
│                                                                 │
│  🖼️ interactive (交互式教材与绘本)                            │
│  └─ Templates: [互动故事生成]                                  │
│                                                                 │
│  📰 news-media (新闻媒体与广电)                               │
│  └─ Templates: [新闻稿生成]                                    │
│                                                                 │
│  🎵 music-audio (音乐与音频娱乐)                             │
│  └─ Templates: [歌词创作]                                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    模板系统 (Templates)                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Template = {                                                  │
│    id: string                  // template id                 │
│    domainId: string            // 所属领域                     │
│    name: string                // 模板名称                     │
│    description: string         // 模板描述                     │
│    basePrompt: string          // 基础提示词模板（含参数占位）  │
│    paramFields: [              // 动态参数字段定义              │
│      {                                                         │
│        key: "duration"         // 参数键                       │
│        label: "时长"           // UI 显示标签                  │
│        type: "slider"          // 字段类型                     │
│        min: 5,                 // 最小值                       │
│        max: 60,                // 最大值                       │
│        defaultValue: 20        // 默认值                       │
│      },                                                        │
│      ...                                                       │
│    ],                                                          │
│    exampleOutput: string       // 示例输出                     │
│  }                                                             │
│                                                                 │
│  User选择模板 → 填写参数 → buildPromptFromTemplate()          │
│               → 最终提示词 → 调用AI API生成                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 文件依赖关系

```
Frontend Dependencies:
Square (页面)
├─ SquareCard (组件)
├─ squareApi (API调用)
├─ getDomainById (helper)
└─ types: SquarePost, SquareComment

AIStudio (页面，改造为向导式)
├─ DomainSelector (Step 1)
├─ DomainParamsForm (Step 3)
├─ getTemplatesByDomain (helper)
├─ buildPromptFromTemplate (helper)
└─ squareApi.publish (发布)

DraftCenter (页面，新增领域筛选)
├─ getDomainsList (helper)
├─ getDomainById (helper)
└─ squareApi.list (查询)

Backend Dependencies:
SquareDomain
├─ SquareService
│  ├─ SquarePostRepository
│  ├─ SquareCommentRepository
│  ├─ SquareLikeRepository
│  └─ SquareCollectRepository
└─ [Other services...]

Database Dependencies:
square_posts → drafts (FK: draft_id)
           → users (FK: user_id)
           
square_comments → square_posts (FK: post_id)
               → users (FK: user_id)
               
square_likes → square_posts (FK: post_id)
            → users (FK: user_id)
            
square_collects → square_posts (FK: post_id)
               → users (FK: user_id)
               
drafts → ✨ NEW: domain field
```

---

## 部署架构

```
┌──────────────────────────────────────┐
│    开发环境                          │
├──────────────────────────────────────┤
│ npm run dev  →  localhost:5173       │
│ go run main.go  →  localhost:8888    │
│ MySQL  →  localhost:3306             │
└──────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│                   生产环境                                │
├──────────────────────────────────────────────────────────┤
│  Front End                                              │
│  ├─ Vite Build (dist/)                                 │
│  ├─ CDN / Web Server                                    │
│  └─ HTTPS + 缓存策略                                     │
│                                                         │
│  Backend                                                │
│  ├─ Docker Container                                    │
│  ├─ Kubernetes / VPS                                    │
│  └─ 负载均衡 + 自动伸缩                                   │
│                                                         │
│  Database                                               │
│  ├─ MySQL Cluster (主从/分片)                           │
│  ├─ Redis Cache (点赞/收藏关系)                          │
│  └─ 定期备份                                             │
│                                                         │
│  Monitoring                                             │
│  ├─ Prometheus + Grafana                               │
│  ├─ 应用日志 (ELK Stack)                                 │
│  └─ 告警系统                                             │
└──────────────────────────────────────────────────────────┘
```

---

**✅ 架构设计完成**

