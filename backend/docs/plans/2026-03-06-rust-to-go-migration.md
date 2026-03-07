# Rust 业务接口迁移到 Go - 实现计划

> **执行方式：** 使用 superpowers:executing-plans 按步骤执行此计划

**目标：** 将 Rust 项目中的 8 个业务 domain（financing, identity, content, configuration, scheduling, publishing, ai_generation, news）完整迁移到 Go，包括 schema、domain 逻辑和 API 接口

**架构：** 按照 Go 项目现有的分层结构（schema → domain → datasource → api），先完整实现 financing domain 作为示例，再按同样模式迁移其他 7 个 domain，最后统一注册到路由

**技术栈：** Go + Gin + GORM + 现有 datasource（dbdao、scylladao、vectordao）

---

## 第一阶段：Financing Domain 完整实现（示例模块）

### Task 1: 创建 financing schema.go

**文件：**
- Create: `/home/x/space/magic/backend/domain/financing/schema.go`

**内容：**

```go
package financing

import "time"

// BusinessPlan 商业计划实体
type BusinessPlan struct {
	ID                int64      `json:"id"`
	Title             string     `json:"title"`
	Content           string     `json:"content"`
	Industry          string     `json:"industry"`
	Region            string     `json:"region"`
	FinancingAmount   float64    `json:"financing_amount"`
	CompanySize       string     `json:"company_size"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// 请求/响应结构体

type GetBusinessPlanReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetBusinessPlanRes struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Industry        string     `json:"industry"`
	Region          string     `json:"region"`
	FinancingAmount float64    `json:"financing_amount"`
	CompanySize     string     `json:"company_size"`
	CreatedAt       *time.Time `json:"created_at"`
}

type GetBusinessPlansReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type CreateBusinessPlanReq struct {
	Title           string  `json:"title" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	Industry        string  `json:"industry" binding:"required"`
	Region          string  `json:"region" binding:"required"`
	FinancingAmount float64 `json:"financing_amount" binding:"required"`
	CompanySize     string  `json:"company_size" binding:"required"`
}

type UpdateBusinessPlanReq struct {
	ID              int64    `json:"id" binding:"required"`
	Title           *string  `json:"title"`
	Content         *string  `json:"content"`
	Industry        *string  `json:"industry"`
	Region          *string  `json:"region"`
	FinancingAmount *float64 `json:"financing_amount"`
	CompanySize     *string  `json:"company_size"`
}

type DeleteBusinessPlanReq struct {
	ID int64 `json:"id" binding:"required"`
}
```

**步骤 1：编写并检查文件**

运行：
```bash
cat /home/x/space/magic/backend/domain/financing/schema.go
```

预期：显示完整的 schema 定义

---

### Task 2: 创建 financing domain.go

**文件：**
- Create: `/home/x/space/magic/backend/domain/financing/domain.go`

**内容：**

```go
package financing

import (
	"backend/datasource/dbdao"
	"context"
	"errors"
)

type FinancingDomain struct {
	DB *dbdao.DB
}

func NewFinancingDomain(db *dbdao.DB) *FinancingDomain {
	return &FinancingDomain{DB: db}
}

// GetBusinessPlan 获取商业计划
func (d *FinancingDomain) GetBusinessPlan(ctx context.Context, id int64) (*BusinessPlan, error) {
	return d.DB.GetBusinessPlan(ctx, id)
}

// GetBusinessPlans 批量获取商业计划
func (d *FinancingDomain) GetBusinessPlans(ctx context.Context, ids []int64) ([]*BusinessPlan, error) {
	if len(ids) == 0 {
		return []*BusinessPlan{}, nil
	}
	return d.DB.GetBusinessPlans(ctx, ids)
}

// CreateBusinessPlan 创建商业计划
func (d *FinancingDomain) CreateBusinessPlan(
	ctx context.Context,
	title, content, industry, region, companySize string,
	financingAmount float64,
) (*BusinessPlan, error) {
	if title == "" || content == "" || industry == "" || region == "" {
		return nil, errors.New("missing required fields")
	}

	plan := &BusinessPlan{
		Title:           title,
		Content:         content,
		Industry:        industry,
		Region:          region,
		FinancingAmount: financingAmount,
		CompanySize:     companySize,
	}

	return d.DB.CreateBusinessPlan(ctx, plan)
}

// UpdateBusinessPlan 更新商业计划
func (d *FinancingDomain) UpdateBusinessPlan(
	ctx context.Context,
	id int64,
	title, content, industry, region, companySize *string,
	financingAmount *float64,
) (*BusinessPlan, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	plan := &BusinessPlan{
		ID:              id,
		Title:           *title,
		Content:         *content,
		Industry:        *industry,
		Region:          *region,
		FinancingAmount: *financingAmount,
		CompanySize:     *companySize,
	}

	return d.DB.UpdateBusinessPlan(ctx, plan)
}

// DeleteBusinessPlan 删除商业计划
func (d *FinancingDomain) DeleteBusinessPlan(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return d.DB.DeleteBusinessPlan(ctx, id)
}
```

**步骤 1：编写并检查文件**

运行：
```bash
cat /home/x/space/magic/backend/domain/financing/domain.go
```

预期：显示完整的 domain 实现

---

### Task 3: 创建 financing datasource (dbdao)

**文件：**
- Create: `/home/x/space/magic/backend/datasource/dbdao/financing.go`

**内容：**

```go
package dbdao

import (
	"backend/domain/financing"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// BusinessPlanModel 数据库模型
type BusinessPlanModel struct {
	ID              int64     `gorm:"primarykey;column:id"`
	Title           string    `gorm:"column:title;type:varchar(255)"`
	Content         string    `gorm:"column:content;type:longtext"`
	Industry        string    `gorm:"column:industry;type:varchar(100)"`
	Region          string    `gorm:"column:region;type:varchar(100)"`
	FinancingAmount float64   `gorm:"column:financing_amount;type:decimal(15,2)"`
	CompanySize     string    `gorm:"column:company_size;type:varchar(50)"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime:milli"`
}

// TableName 指定表名
func (BusinessPlanModel) TableName() string {
	return "business_plans"
}

// GetBusinessPlan 获取商业计划
func (d *DB) GetBusinessPlan(ctx context.Context, id int64) (*financing.BusinessPlan, error) {
	var model BusinessPlanModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("business plan not found")
		}
		return nil, result.Error
	}

	return &financing.BusinessPlan{
		ID:              model.ID,
		Title:           model.Title,
		Content:         model.Content,
		Industry:        model.Industry,
		Region:          model.Region,
		FinancingAmount: model.FinancingAmount,
		CompanySize:     model.CompanySize,
		CreatedAt:       &model.CreatedAt,
		UpdatedAt:       &model.UpdatedAt,
	}, nil
}

// GetBusinessPlans 批量获取商业计划
func (d *DB) GetBusinessPlans(ctx context.Context, ids []int64) ([]*financing.BusinessPlan, error) {
	var models []BusinessPlanModel
	result := d.DB().WithContext(ctx).Where("id IN ?", ids).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	plans := make([]*financing.BusinessPlan, len(models))
	for i, model := range models {
		plans[i] = &financing.BusinessPlan{
			ID:              model.ID,
			Title:           model.Title,
			Content:         model.Content,
			Industry:        model.Industry,
			Region:          model.Region,
			FinancingAmount: model.FinancingAmount,
			CompanySize:     model.CompanySize,
			CreatedAt:       &model.CreatedAt,
			UpdatedAt:       &model.UpdatedAt,
		}
	}

	return plans, nil
}

// CreateBusinessPlan 创建商业计划
func (d *DB) CreateBusinessPlan(ctx context.Context, plan *financing.BusinessPlan) (*financing.BusinessPlan, error) {
	model := &BusinessPlanModel{
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
	}

	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &financing.BusinessPlan{
		ID:              model.ID,
		Title:           model.Title,
		Content:         model.Content,
		Industry:        model.Industry,
		Region:          model.Region,
		FinancingAmount: model.FinancingAmount,
		CompanySize:     model.CompanySize,
		CreatedAt:       &model.CreatedAt,
		UpdatedAt:       &model.UpdatedAt,
	}, nil
}

// UpdateBusinessPlan 更新商业计划
func (d *DB) UpdateBusinessPlan(ctx context.Context, plan *financing.BusinessPlan) (*financing.BusinessPlan, error) {
	model := &BusinessPlanModel{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
	}

	result := d.DB().WithContext(ctx).Model(&model).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("business plan not found")
	}

	return &financing.BusinessPlan{
		ID:              model.ID,
		Title:           model.Title,
		Content:         model.Content,
		Industry:        model.Industry,
		Region:          model.Region,
		FinancingAmount: model.FinancingAmount,
		CompanySize:     model.CompanySize,
		CreatedAt:       &model.CreatedAt,
		UpdatedAt:       &model.UpdatedAt,
	}, nil
}

// DeleteBusinessPlan 删除商业计划
func (d *DB) DeleteBusinessPlan(ctx context.Context, id int64) error {
	result := d.DB().WithContext(ctx).Where("id = ?", id).Delete(&BusinessPlanModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("business plan not found")
	}

	return nil
}
```

**步骤 1：编写并检查文件**

运行：
```bash
cat /home/x/space/magic/backend/datasource/dbdao/financing.go
```

预期：显示完整的 datasource 实现

---

### Task 4: 创建 financing api.go

**文件：**
- Create: `/home/x/space/magic/backend/domain/financing/api.go`

**内容：**

```go
package financing

import (
	"github.com/cv70/pkgo/ghttp"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiGetBusinessPlan 获取商业计划详情
func (d *FinancingDomain) ApiGetBusinessPlan(c *gin.Context) {
	var req GetBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.GetBusinessPlan(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to get business plan: %v", err)
		ghttp.RespError(c, 500, "failed to get business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	ghttp.RespSuccess(c, res)
}

// ApiGetBusinessPlans 批量获取商业计划
func (d *FinancingDomain) ApiGetBusinessPlans(c *gin.Context) {
	var req GetBusinessPlansReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	plans, err := d.GetBusinessPlans(c.Request.Context(), req.IDs)
	if err != nil {
		klog.Errorf("failed to get business plans: %v", err)
		ghttp.RespError(c, 500, "failed to get business plans")
		return
	}

	responses := make([]*GetBusinessPlanRes, len(plans))
	for i, plan := range plans {
		responses[i] = &GetBusinessPlanRes{
			ID:              plan.ID,
			Title:           plan.Title,
			Content:         plan.Content,
			Industry:        plan.Industry,
			Region:          plan.Region,
			FinancingAmount: plan.FinancingAmount,
			CompanySize:     plan.CompanySize,
			CreatedAt:       plan.CreatedAt,
		}
	}

	ghttp.RespSuccess(c, responses)
}

// ApiCreateBusinessPlan 创建商业计划
func (d *FinancingDomain) ApiCreateBusinessPlan(c *gin.Context) {
	var req CreateBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.CreateBusinessPlan(
		c.Request.Context(),
		req.Title,
		req.Content,
		req.Industry,
		req.Region,
		req.CompanySize,
		req.FinancingAmount,
	)
	if err != nil {
		klog.Errorf("failed to create business plan: %v", err)
		ghttp.RespError(c, 500, "failed to create business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	ghttp.RespSuccess(c, res)
}

// ApiUpdateBusinessPlan 更新商业计划
func (d *FinancingDomain) ApiUpdateBusinessPlan(c *gin.Context) {
	var req UpdateBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	plan, err := d.UpdateBusinessPlan(
		c.Request.Context(),
		req.ID,
		req.Title,
		req.Content,
		req.Industry,
		req.Region,
		req.CompanySize,
		req.FinancingAmount,
	)
	if err != nil {
		klog.Errorf("failed to update business plan: %v", err)
		ghttp.RespError(c, 500, "failed to update business plan")
		return
	}

	res := &GetBusinessPlanRes{
		ID:              plan.ID,
		Title:           plan.Title,
		Content:         plan.Content,
		Industry:        plan.Industry,
		Region:          plan.Region,
		FinancingAmount: plan.FinancingAmount,
		CompanySize:     plan.CompanySize,
		CreatedAt:       plan.CreatedAt,
	}

	ghttp.RespSuccess(c, res)
}

// ApiDeleteBusinessPlan 删除商业计划
func (d *FinancingDomain) ApiDeleteBusinessPlan(c *gin.Context) {
	var req DeleteBusinessPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	err := d.DeleteBusinessPlan(c.Request.Context(), req.ID)
	if err != nil {
		klog.Errorf("failed to delete business plan: %v", err)
		ghttp.RespError(c, 500, "failed to delete business plan")
		return
	}

	ghttp.RespSuccess(c, gin.H{"id": req.ID})
}
```

**步骤 1：编写并检查文件**

运行：
```bash
cat /home/x/space/magic/backend/domain/financing/api.go
```

预期：显示完整的 API 实现

---

### Task 5: 在 main.go 中注册 financing 路由

**文件修改：**
- Modify: `/home/x/space/magic/backend/main.go`

在 `func main()` 中添加 financing domain 的注册代码：

```go
// financing 融资模块
{
    financingDomain := financing.NewFinancingDomain(registry.DB)
    v1.POST("/financing/business-plans/get", financingDomain.ApiGetBusinessPlan)
    v1.POST("/financing/business-plans/batch-get", financingDomain.ApiGetBusinessPlans)
    v1.POST("/financing/business-plans/create", financingDomain.ApiCreateBusinessPlan)
    v1.POST("/financing/business-plans/update", financingDomain.ApiUpdateBusinessPlan)
    v1.POST("/financing/business-plans/delete", financingDomain.ApiDeleteBusinessPlan)
}
```

**步骤 1：修改并检查**

运行：
```bash
grep -A 5 "financing" /home/x/space/magic/backend/main.go
```

预期：显示 financing 相关的路由注册代码

---

### Task 6: 编译并测试 financing domain

**步骤 1：构建项目**

运行：
```bash
cd /home/x/space/magic/backend && go build -v
```

预期：编译成功，没有错误

**步骤 2：提交 financing domain 完整实现**

运行：
```bash
cd /home/x/space/magic/backend && git add domain/financing/ datasource/dbdao/financing.go main.go && git commit -m "feat: implement financing domain (schema, domain, datasource, api)"
```

预期：成功提交

---

## 第二阶段：迁移 Identity Domain

参考第一阶段的 financing 模式，按以下顺序实现：

1. 创建 `domain/identity/schema.go` - 参考 Rust `src/domain/identity/schema.rs`
2. 创建 `domain/identity/domain.go` - 参考 Rust `src/domain/identity/domain.rs`
3. 创建 `datasource/dbdao/identity.go` - 实现 User、Role、Permission 相关的数据库操作
4. 创建 `domain/identity/api.go` - 参考 Rust `src/domain/identity/api.rs`
5. 在 `main.go` 中注册 identity 路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- User（id, username, email, password, role, enabled, created_at）
- Role（id, name, description, created_at）
- Permission（id, name, description, created_at）

---

## 第三阶段：迁移 Content Domain

参考 financing 模式实现：

1. `domain/content/schema.go`
2. `domain/content/domain.go`
3. `datasource/dbdao/content.go`
4. `domain/content/api.go`
5. 在 `main.go` 中注册路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- Content（id, title, body, content_type, status, tag, created_at）

---

## 第四阶段：迁移 Configuration Domain

参考 financing 模式实现：

1. `domain/configuration/schema.go`
2. `domain/configuration/domain.go`
3. `datasource/dbdao/configuration.go`
4. `domain/configuration/api.go`
5. 在 `main.go` 中注册路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- Configuration（id, key, value, created_at, updated_at）

---

## 第五阶段：迁移 Scheduling Domain

参考 financing 模式实现：

1. `domain/scheduling/schema.go`
2. `domain/scheduling/domain.go`
3. `datasource/dbdao/scheduling.go`
4. `domain/scheduling/api.go`
5. 在 `main.go` 中注册路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- Task（id, title, description, status, scheduled_time, created_at）

---

## 第六阶段：迁移 Publishing Domain

参考 financing 模式实现：

1. `domain/publishing/schema.go`
2. `domain/publishing/domain.go`
3. `datasource/dbdao/publishing.go`
4. `domain/publishing/api.go`
5. 在 `main.go` 中注册路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- Publication（id, title, content, status, published_at, created_at）

---

## 第七阶段：迁移 AI Generation Domain

参考 financing 模式实现：

1. `domain/ai_generation/schema.go`
2. `domain/ai_generation/domain.go`
3. `datasource/dbdao/ai_generation.go`
4. `domain/ai_generation/api.go`
5. 在 `main.go` 中注册路由
6. 编译测试并提交

**关键数据结构（参考 Rust）：**
- GeneratedContent（id, prompt, content, model, created_at）

---

## 第八阶段：补充迁移 News Domain

虽然 news domain 在 Go 中已部分实现，但需要完整的 datasource 支持：

1. 检查并完善 `datasource/dbdao/news.go`
2. 如果使用 scylladao，添加相关实现
3. 如果使用 vectordao，添加向量搜索实现
4. 确保 `domain/news/domain.go` 和 `domain/news/api.go` 完整

---

## 第九阶段：最终集成和验证

1. **验证所有路由** - 检查 main.go 中所有 domain 都已正确注册
2. **创建 migration 脚本** - 生成创建所有数据表的 SQL 脚本
3. **测试所有端点** - 使用 Postman 或 curl 测试每个 domain 的 API
4. **性能测试** - 运行负载测试确保 Go 版本性能良好
5. **最终提交** - 总结迁移完成，提交 final commit

---

## 注意事项

1. **数据库初始化** - 在执行任何 API 调用前，需要先创建对应的数据表
2. **错误处理** - 统一使用 `ghttp.RespError` 返回错误信息
3. **日志记录** - 使用 `klog` 进行日志记录，记录所有关键操作
4. **上下文传递** - 所有 datasource 操作都应该传递 `context.Context`
5. **时间戳** - 所有时间戳字段使用 `time.Time`，并配置为数据库自动更新

---

## 执行方式选择

**计划已完成并保存到 `docs/plans/2026-03-06-rust-to-go-migration.md`**

两种执行选项：

**1. 当前会话执行（推荐）** - 我将在当前会话中逐步实现，每个关键步骤后都会验证和提交
```bash
使用 superpowers:executing-plans 进行当前会话执行
```

**2. 独立会话执行** - 你可以在单独的 session 中按计划步骤执行
```bash
在新会话中使用 superpowers:executing-plans，参考此计划
```

你的选择是？

