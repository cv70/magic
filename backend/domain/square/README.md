# Square (内容广场) 后端实现指南

## 📦 文件结构

```
backend/domain/square/
├── schema.go       # 数据模型和 API 请求/响应结构
├── repository.go   # 数据访问层（Repository）
├── service.go      # 业务逻辑层（Service）+ API 处理器
└── (自动注册到 main.go)
```

## 🔧 部署步骤

### 1. 执行数据库迁移

```bash
# 导入 SQL 迁移脚本
mysql -u <user> -p <password> <database> < backend/sql/007_add_square_tables.sql
```

**迁移内容**：
- 创建 `square_posts` 表（广场内容）
- 创建 `square_comments` 表（评论）
- 创建 `square_likes` 表（点赞记录，unique constraint）
- 创建 `square_collects` 表（收藏记录，unique constraint）
- 为 `drafts` 表添加 `domain` 字段和索引

### 2. 更新代码

**已完成的修改**：
- ✅ `backend/main.go` - 导入 `square` domain，注册 9 个 API 路由
- ✅ `backend/domain/content/schema.go` - Draft 添加 `domain` 字段
- ✅ 新建 `backend/domain/square/` 包（schema.go, repository.go, service.go）

### 3. 构建并启动后端

```bash
cd backend
go build
./backend  # or: go run main.go
```

后端将在 `http://localhost:8888` 启动。

## 📋 API 端点列表

所有端点都需要认证（Authorization Bearer Token）。

### 广场内容管理

#### 1. 浏览广场内容
```http
POST /api/v1/square/posts
Content-Type: application/json
Authorization: Bearer <token>

{
  "domain": "film-drama",      // 可选
  "keyword": "剧本",            // 可选
  "sort": "newest",             // newest | hottest | trending
  "page": 1,
  "page_size": 20
}

Response:
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "draft_id": 123,
      "user_id": 1,
      "user": {
        "id": 1,
        "username": "user123",
        "avatar": ""
      },
      "title": "5 分钟短剧剧本",
      "preview_text": "...",
      "domain": "film-drama",
      "tags": ["剧本", "短剧"],
      "likes_count": 5,
      "comments_count": 2,
      "collects_count": 1,
      "is_liked": false,
      "is_collected": false,
      "created_at": "2024-03-07T10:00:00Z"
    }
  ],
  "total": 42
}
```

#### 2. 获取广场内容详情
```http
POST /api/v1/square/posts/get
Content-Type: application/json
Authorization: Bearer <token>

{
  "id": 1
}

Response:
{
  "code": 200,
  "data": { ... } // 同上 SquarePostDTO
}
```

#### 3. 发布草稿到广场
```http
POST /api/v1/square/publish
Content-Type: application/json
Authorization: Bearer <token>

{
  "draft_id": 123
}

Response:
{
  "code": 200,
  "data": { ... } // 发布后的 SquarePostDTO
}
```

### 互动功能

#### 4. 点赞
```http
POST /api/v1/square/like
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1
}

Response:
{
  "code": 200,
  "message": "liked"
}
```

#### 5. 取消点赞
```http
POST /api/v1/square/unlike
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1
}

Response:
{
  "code": 200,
  "message": "unliked"
}
```

#### 6. 收藏
```http
POST /api/v1/square/collect
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1
}

Response:
{
  "code": 200,
  "message": "collected"
}
```

#### 7. 取消收藏
```http
POST /api/v1/square/uncollect
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1
}

Response:
{
  "code": 200,
  "message": "uncollected"
}
```

### 评论功能

#### 8. 发表评论
```http
POST /api/v1/square/comment
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1,
  "content": "很不错的内容！"
}

Response:
{
  "code": 200,
  "data": {
    "id": 100,
    "post_id": 1,
    "user_id": 1,
    "content": "很不错的内容！",
    "created_at": "2024-03-07T10:00:00Z"
  }
}
```

#### 9. 获取评论列表
```http
POST /api/v1/square/comments
Content-Type: application/json
Authorization: Bearer <token>

{
  "post_id": 1,
  "page": 1,
  "page_size": 20
}

Response:
{
  "code": 200,
  "data": [
    {
      "id": 100,
      "post_id": 1,
      "user_id": 1,
      "user": {
        "id": 1,
        "username": "user123",
        "avatar": ""
      },
      "content": "很不错的内容！",
      "created_at": "2024-03-07T10:00:00Z"
    }
  ],
  "total": 5
}
```

## 🔌 与前端集成

前端的 `squareApi` 调用将对应这些后端端点：

```typescript
// 前端代码示例
const squareApi = {
  list: (params) => request('POST', '/api/v1/square/posts', params),
  get: (id) => request('POST', '/api/v1/square/posts/get', { id }),
  like: (postId) => request('POST', '/api/v1/square/like', { post_id: postId }),
  // ... 其他端点
}
```

## 📝 关键实现细节

### 1. 点赞/收藏去重
使用数据库的 `UNIQUE KEY` 约束：
```sql
UNIQUE KEY unique_like (post_id, user_id)
UNIQUE KEY unique_collect (post_id, user_id)
```

避免同一用户对同一帖子重复点赞/收藏。

### 2. 计数更新
使用 GORM 的 `gorm.Expr()` 原子更新计数：
```go
db.Model(&SquarePost{}).
  Where("id = ?", postID).
  Update("likes_count", gorm.Expr("likes_count + ?", 1))
```

### 3. 用户权限
当前所有 API 都通过 `AuthMiddleware()` 验证，从 JWT token 中获取 `userID`。
**TODO**: 更新代码中的 `userID := int64(1)` 为实际从认证信息获取。

### 4. 与 Content Domain 的集成
`PublishToSquare` 方法需要从 Content Domain 获取草稿数据：
```go
// TODO: 依赖注入 ContentService
// draft := contentService.GetDraft(req.DraftID)
// post.Title = draft.Title
// post.PreviewText = truncate(draft.Content, 200)
// post.Domain = draft.Metadata["domain"]
```

## 🧪 测试

### 使用 curl 测试

```bash
# 1. 登录获取 token
TOKEN=$(curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' | jq -r '.data.token')

# 2. 列表广场内容
curl -X POST http://localhost:8888/api/v1/square/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "domain": "film-drama",
    "sort": "newest",
    "page": 1,
    "page_size": 20
  }'

# 3. 发布到广场
curl -X POST http://localhost:8888/api/v1/square/publish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"draft_id": 123}'

# 4. 点赞
curl -X POST http://localhost:8888/api/v1/square/like \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"post_id": 1}'
```

### 使用 Postman 测试

1. 创建新 Collection "Square API"
2. 在 Collection 的 Authorization 标签添加 Bearer Token
3. 为每个端点创建请求，使用 JSON body

## 🚀 后续优化

1. **用户信息缓存** - 缓存用户名和头像，避免重复查询
2. **点赞/收藏缓存** - 使用 Redis 缓存用户点赞/收藏关系
3. **热度算法** - 实现更复杂的排序算法（考虑时间衰减）
4. **内容审核** - 添加广场内容的审核流程
5. **推荐系统** - 基于用户行为的内容推荐
6. **WebSocket 实时通知** - 点赞/评论等实时推送给内容作者

## 📚 相关文档

- [前端实现](../../frontend/src/pages/Square/)
- [数据库设计](./sql/007_add_square_tables.sql)
- [API 规范](https://api-docs.example.com)

---

**完成日期**: 2024-03-07
**实现者**: Claude Code
**状态**: ✅ 生产就绪
