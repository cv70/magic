# 后端集成检查表

## 📋 手动集成步骤（预计 15-30 分钟）

### Step 1: 数据库迁移
```bash
cd backend
mysql -u root -p magic < sql/007_add_square_tables.sql
```
**验证**:
```bash
mysql -u magic -p magic -e "SHOW TABLES LIKE 'square_%';" magic
# 应显示 4 个表: square_posts, square_comments, square_likes, square_collects
```

---

### Step 2: 提取用户认证信息

**文件**: `backend/domain/square/service.go`

需要修改的位置（3 处）:

#### 位置 1: ListPosts Handler (第 ~380 行)
```go
// 当前代码
userID := int64(1)  // ❌ 硬编码

// 改为
userID := service.extractUserIDFromContext(c)
if userID == 0 {
    c.JSON(400, gin.H{"error": "user not found"})
    return
}
```

#### 位置 2: Like/Unlike/Collect/Uncollect Handlers (第 ~420-460 行)
```go
// 当前代码
userID := int64(1)  // ❌ 硬编码

// 改为
userID := service.extractUserIDFromContext(c)
if userID == 0 {
    c.JSON(400, gin.H{"error": "user not found"})
    return
}
```

#### 位置 3: AddComment Handler (第 ~490 行)
```go
// 当前代码
userID := int64(1)  // ❌ 硬编码

// 改为
userID := service.extractUserIDFromContext(c)
if userID == 0 {
    c.JSON(400, gin.H{"error": "user not found"})
    return
}
```

**添加辅助方法** (在 service.go 底部):
```go
func (s *SquareService) extractUserIDFromContext(c *gin.Context) int64 {
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        return 0
    }
    userID, ok := userIDInterface.(int64)
    if !ok {
        return 0
    }
    return userID
}
```

---

### Step 3: 集成 Content Domain

**文件**: `backend/domain/square/service.go`

**位置**: `ApiPublishToSquare` 方法 (第 ~320 行)

当前是 TODO，需要实现从 Draft 加载内容:

```go
// TODO: 获取 Draft 内容
// 需要依赖注入 ContentService

// 改为
draft, err := s.contentService.GetDraft(c.Context(), req.DraftID)
if err != nil {
    c.JSON(400, gin.H{"error": "draft not found"})
    return
}

// 提取 domain 字段
domain := draft.Domain
if domain == "" {
    c.JSON(400, gin.H{"error": "draft domain not set"})
    return
}
```

**在 NewSquareDomainWithDB 中注入**:
```go
type SquareDomain struct {
    db             *gorm.DB
    squareService  *SquareService
    contentService ContentServiceInterface  // ← 新增
}

// 工厂方法改为
func NewSquareDomainWithDB(db *gorm.DB, contentService ContentServiceInterface) *SquareDomain {
    return &SquareDomain{
        db:             db,
        contentService: contentService,
        squareService:  NewSquareService(db, contentService),
    }
}
```

---

### Step 4: 加载用户信息

**文件**: `backend/domain/square/service.go`

**位置**: `ListPosts` 和 `GetPost` 中的 DTOs 转换

当前每个 SquarePostDTO 有 UserID 但 User 对象为空:

```go
// 在构建 DTO 时添加用户信息加载
for _, post := range posts {
    user, _ := s.identityService.GetUser(c.Context(), post.UserID)
    dto.User = &UserInfoDTO{
        ID:       user.ID,
        Username: user.Username,
        Avatar:   user.Avatar,
    }
}
```

**在 NewSquareDomainWithDB 中注入**:
```go
type SquareDomain struct {
    db              *gorm.DB
    contentService  ContentServiceInterface
    identityService IdentityServiceInterface  // ← 新增
}
```

---

### Step 5: 验证编译

```bash
cd backend
go mod tidy
go build -o backend

# 如果有编译错误，根据错误信息修改上述实现
```

---

### Step 6: 本地测试

```bash
# 启动后端
./backend

# 在另一个终端测试 API（需要替换 <token> 为有效的 JWT）
curl -X POST http://localhost:8888/api/v1/square/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "page": 1,
    "page_size": 20,
    "domain": "film-drama"
  }'

# 预期响应
{
  "code": 200,
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 📝 集成前端和后端

### 前端接口已准备就绪
- ✅ `squareApi.list()` - 列表 API 调用
- ✅ `squareApi.get()` - 详情 API 调用
- ✅ `squareApi.like/unlike/collect/uncollect()` - 互动 API
- ✅ `squareApi.comment()` - 评论 API

### 端到端测试步骤

1. **启动前端** (`npm run dev`)
2. **启动后端** (`./backend`)
3. **访问内容广场**: `http://localhost:5173/square`
4. **执行测试**:
   - [ ] 加载广场列表（应为空）
   - [ ] 创建草稿并设置 domain
   - [ ] 发布到广场
   - [ ] 列表显示新帖
   - [ ] 点赞/收藏/评论功能正常
   - [ ] 领域过滤生效
   - [ ] 搜索功能生效

---

## 🔍 疑难排查

### 编译错误: `ContentServiceInterface not found`
**解决**: 在 `domain/content/` 中查找 ContentService 的接口定义，导入该接口

### 编译错误: `extractUserIDFromContext` 不存在
**解决**: 确认已添加辅助方法到 SquareService 结构体

### API 返回 401 Unauthorized
**解决**: 确认 JWT token 有效，检查 Authorization header 格式为 `Bearer <token>`

### 数据库操作失败
**解决**:
1. 确认 007_add_square_tables.sql 已执行
2. 检查 MySQL 连接字符串正确
3. 检查表和索引已创建: `SHOW CREATE TABLE square_posts;`

---

## ✅ 完成标志

当以下条件全部满足时，集成完成：

- [ ] 数据库迁移成功（4 个表已创建）
- [ ] 后端编译成功（`go build` 无错误）
- [ ] 后端启动成功（输出 "Listening on :8888"）
- [ ] API 测试通过（POST /api/v1/square/posts 返回 200）
- [ ] 前端加载广场页面无错误
- [ ] 点赞/收藏等互动功能正常

---

**预计时间**: 15-30 分钟
**难度**: 低（主要是配置和集成，逻辑已完成）
**支持**: 如遇问题，参考各文件的详细文档
