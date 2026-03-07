# Domain Layer Architecture

## Overview

The domain layer in this project follows a Domain-Driven Design (DDD) pattern with clear separation of concerns across domain modules. The project has been migrated from Rust + Axum to **Go + Gin**, and the domain layer now uses Go's standard patterns and popular Go DDD libraries.

## Current State

The domain layer is organized into domain-specific modules, each containing schema definitions and domain logic components. Each domain module follows this pattern:
- `model.go` - Domain entity definitions
- `repository.go` - Data access interfaces
- `service.go` - Business logic
- `handler.go` - HTTP handlers (Gin)
- `dto.go` - Data Transfer Objects

## Domain Modules

### 1. News Domain (`domain/news/`)
**Status**: Partial Implementation
**Key Components**:
- `model.go`: `News` struct with GORM tags
- `repository.go`: News repository interface and implementation
- `service.go`: News business logic
- `handler.go`: API handlers for news CRUD operations
- `dto.go`: News DTOs (GetNewsReq/Res, SearchNewsReq/Res, AddNewsReq/Res)

### 2. Financing Domain (`domain/financing/`)
**Status**: Partial Implementation
**Key Components**:
- `model.go`: Business plan related models
- `repository.go`: Financing repository interface
- `service.go`: Business plan service
- `handler.go`: Business plan API endpoints

### 3. Content Domain (`domain/content/`)
**Status**: Schema Defined
**Key Components**:
- `model.go`: Content, ContentVersion, ContentType, Tag models
- `repository.go`: Content repository interface
- `service.go`: Content service for CRUD operations
- `handler.go`: Content CRUD API endpoints

### 4. AI Generation Domain (`domain/ai_generation/`)
**Status**: Partial Implementation
**Key Components**:
- `model.go`: AI generation related models
- `repository.go`: AI configuration repository
- `service.go`: AI generation service with multi-provider support
- `handler.go`: AI generation API endpoints

### 5. Publishing Domain (`domain/publishing/`)
**Status**: Schema Defined
**Key Components**:
- `model.go`: Publisher, PublishTask, PublishPlatform models
- `repository.go`: Publishing repository interface
- `service.go`: Publishing service for multi-platform publishing
- `handler.go`: Publishing API endpoints

### 6. Scheduling Domain (`domain/scheduling/`)
**Status**: Schema Defined
**Key Components**:
- `model.go`: Scheduler, CronExpression, ScheduledTask models
- `repository.go`: Scheduling repository interface
- `service.go`: Scheduling service with cron support
- `handler.go`: Scheduling API endpoints

### 7. Configuration Domain (`domain/configuration/`)
**Status**: Schema Defined
**Key Components**:
- `model.go`: SystemConfig, ProviderConfig models
- `repository.go`: Configuration repository interface
- `service.go`: Configuration service
- `handler.go`: Configuration API endpoints

### 8. Identity Domain (`domain/identity/`)
**Status**: Partial Implementation
**Key Components**:
- `model.go`: User, Role, Permission models
- `repository.go`: Identity repository interface
- `service.go`: Authentication and authorization service

## Implementation Patterns

### Domain Module Structure (Go)
```go
package content

import (
    "context"
    "time"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// Model definitions
type Content struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"size:255" json:"title"`
    Content     string         `gorm:"type:text" json:"content"`
    ContentType string         `gorm:"size:20" json:"content_type"`
    Author      string         `gorm:"size:100" json:"author"`
    Status      string         `gorm:"size:20" json:"status"`
    Tags        string         `gorm:"size:500" json:"tags"`
    Metadata    string         `gorm:"type:text" json:"metadata"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
}

// Repository interface
type ContentRepository interface {
    Create(ctx context.Context, content *Content) error
    GetByID(ctx context.Context, id uint) (*Content, error)
    List(ctx context.Context, limit, offset int) ([]*Content, error)
    Update(ctx context.Context, content *Content) error
    Delete(ctx context.Context, id uint) error
}

// Service implementation
type ContentService struct {
    db *gorm.DB
}

func NewContentService(db *gorm.DB) *ContentService {
    return &ContentService{db: db}
}

func (s *ContentService) Create(ctx context.Context, content *Content) error {
    return s.db.WithContext(ctx).Create(content).Error
}

// Handler (Gin)
type ContentHandler struct {
    service *ContentService
}

func NewContentHandler(service *ContentService) *ContentHandler {
    return &ContentHandler{service: service}
}

func (h *ContentHandler) Create(c *gin.Context) {
    var req CreateContentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // Business logic
}
```

### DTO Pattern
```go
// Request DTOs
type CreateContentRequest struct {
    Title       string `json:"title" binding:"required"`
    Content     string `json:"content" binding:"required"`
    ContentType string `json:"content_type" binding:"required"`
    Author      string `json:"author"`
    Tags        string `json:"tags"`
}

// Response DTOs
type ContentResponse struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    ContentType string    `json:"content_type"`
    Author      string    `json:"author"`
    Status      string    `json:"status"`
    Tags        string    `json:"tags"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### API Handler Pattern (Gin)
```go
func (h *ContentHandler) RegisterRoutes(r *gin.RouterGroup) {
    r.GET("/contents", h.List)
    r.GET("/contents/:id", h.Get)
    r.POST("/contents", h.Create)
    r.PUT("/contents/:id", h.Update)
    r.DELETE("/contents/:id", h.Delete)
}

func (h *ContentHandler) Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid id"})
        return
    }
    
    content, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    
    c.JSON(200, content)
}
```

## Domain Relationships

```
ContentDomain ──┬─> AIGenerationDomain
                ├─> PublishingDomain
                └─> SchedulingDomain
FinancingDomain ──> ContentDomain
PublishingDomain ──> SchedulingDomain
IdentityDomain ────> All domains (authentication/authorization)
```

## Data Flow

```
Client Request
    ↓
Gin Router (handler.go)
    ↓
Service Layer (service.go)
    ↓
Repository Layer (repository.go)
    ↓
Database (GORM)
    ↓
Response Data
    ↓
JSON Response
    ↓
Client Response
```

## Technology Stack for Domain Layer

| Layer | Technology |
|-------|------------|
| **Web Framework** | Gin |
| **ORM** | GORM |
| **Database** | PostgreSQL / MySQL / SQLite |
| **Validation** | go-playground/validator |
| **Configuration** | spf13/viper |
| **Dependency Injection** | google/wire (optional) |
| **Testing** | testify |

## Migration Notes (Rust to Go)

| Rust (Axum) | Go (Gin) |
|-------------|----------|
| `#[derive(Serialize, Deserialize)]` | `json:"field"` tags |
| `async fn` | Goroutines + context |
| `Result<T, E>` | `error` returns |
| `Arc<T>` | Pointer receivers |
| `tokio` runtime | Native Go concurrency |
| `sqlx` | GORM |
| `axum::extract::State` | `gin.Context` |
| `axum::response::IntoResponse` | `gin.Context.JSON()` |

## Future Implementation Plans

1. Implement full domain logic in service.go files
2. Add domain events for async processing
3. Implement domain validation logic
4. Add domain-specific error handling
5. Add unit and integration tests for each domain module
6. Document domain-specific business rules
7. Add domain-driven design patterns (aggregates, value objects)

## Development Notes

- This is a **Go + Gin** backend project (migrated from Rust + Axum)
- Uses GORM as the primary ORM
- Follows standard Go project layout (not strictly, but idiomatic Go)
- Each domain has its own module with model.go, repository.go, service.go, handler.go, dto.go structure
- Most domains are currently in placeholder state awaiting full implementation
- Domain logic focuses on clean separation between HTTP handlers and business logic
- GORM is used for database operations with proper transaction support
