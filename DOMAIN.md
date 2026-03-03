# Domain Layer Architecture

## Overview

The domain layer in this project follows a Domain-Driven Design (DDD) pattern with clear separation of concerns across domain modules. The layer is organized into domain-specific modules, each containing schema definitions and domain logic components.

## Current State

The domain layer is currently under active development. Most domains have their module structure in place with schema definitions, but the domain logic implementations are primarily placeholders awaiting full implementation. Each domain module follows this pattern:
- `mod.rs` - Module declarations
- `domain.rs` - Domain struct definitions (currently placeholder implementations)
- `schema.rs` - Data schemas and DTOs
- `api.rs` - API endpoint handlers

## Domain Modules

### 1. News Domain (`domain/news/`)
**Status**: Partial Implementation
**Key Components**:
- `domain.rs`: `NewsDomain` struct with database accessors
- `schema.rs`: News, GetNewsReq/Res, SearchNewsReq/Res, AddNewsReq/Res
- `api.rs`: API endpoints for news CRUD operations

### 2. Financing Domain (`domain/financing/`)
**Status**: Partial Implementation
**Key Components**:
- `domain.rs`: `FinancingDomain` struct with database accessors
- `domain.rs`: Business plan related structures
- `schema.rs`: Business plan schemas
- `api.rs`: Business plan API endpoints

### 3. Content Domain (`domain/content/`)
**Status**: Schema Defined
**Key Components**:
- `domain.rs`: `ContentDomain` struct with database accessors
- `schema.rs`: Content, ContentVersion, ContentType, Tag schemas
- `api.rs`: Content CRUD API endpoints

### 4. AI Generation Domain (`domain/ai_generation/`)
**Status**: Partial Implementation
**Key Components**:
- `domain.rs`: `AIGenerationDomain` struct with database accessors
- `schema.rs`: AI-related schemas
- `api.rs`: AI generation API endpoints

### 5. Publishing Domain (`domain/publishing/`)
**Status**: Schema Defined
**Key Components**:
- `domain.rs`: `PublishingDomain` struct with database accessors
- `schema.rs`: Publisher, PublishTask, PublishPlatform schemas
- `api.rs`: Publishing API endpoints

### 6. Scheduling Domain (`domain/scheduling/`)
**Status**: Schema Defined
**Key Components**:
- `mod.rs`: Module structure with mod.rs, schema.rs, task.rs
- `domain.rs`: `SchedulingDomain` struct with database accessors
- `schema.rs`: Scheduler, CronExpression, ScheduledTask schemas
- `api.rs`: Scheduling API endpoints

### 7. Configuration Domain (`domain/configuration/`)
**Status**: Schema Defined
**Key Components**:
- `domain.rs`: `ConfigurationDomain` struct with database accessors
- `schema.rs`: SystemConfig, ProviderConfig schemas
- `api.rs`: Configuration API endpoints

### 8. Identity Domain (`domain/identity/`)
**Status**: Partial Implementation
**Key Components**:
- `domain.rs`: `IdentityDomain` struct with database accessors
- `schema.rs`: User, Role, Permission schemas

## Implementation Patterns

### Domain Module Structure
```rust
use std::sync::Arc;

use crate::datasource::{dbdao::dao::DBDao, scylladao::dao::ScyllaDao, vectordao::dao::VectorDao};
use crate::infra::registry::Registry;

#[derive(Clone)]
pub struct [DomainName]Domain {
    pub db_dao: Arc<DBDao>,
    pub scylla_dao: Arc<ScyllaDao>,
    pub vector_dao: Arc<VectorDao>,
}

impl [DomainName]Domain {
    pub fn new(r: &Registry) -> Self {
        Self {
            db_dao: r.db_dao.clone(),
            scylla_dao: r.scylla_dao.clone(),
            vector_dao: r.vector_dao.clone(),
        }
    }
}
```

### Schema Pattern
All schemas use `serde` for serialization/deserialization:
```rust
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct [EntityName] {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    // ... other fields
}
```

### API Pattern
API endpoints use Axum extractors:
```rust
use axum::extract::State;
use axum::response::IntoResponse;
use crate::state::state::AppState;

pub async fn api_get_[entity](
    State(state): State<AppState>,
    Json(req): Json<Get[Entity]Req>
) -> impl IntoResponse {
    // API implementation
}
```

## Domain Relationships

```
ContentDomain ──┬─> AIGenerationDomain
                ├─> PublishingDomain
                └─> SchedulingDomain
FinancingDomain ──> ContentDomain
PublishingDomain ──> SchedulingDomain
IdentityDomain ────> All domains (authentication)
```

## Data Flow

```
Client Request
    ↓
API Layer (api.rs)
    ↓
Domain Layer (domain.rs)
    ↓
Database Access (DBDao, ScyllaDao, VectorDao)
    ↓
Response Data
    ↓
API Response
    ↓
Client Response
```

## Future Implementation Plans

1. Implement domain logic in domain.rs files
2. Add domain services for business logic
3. Implement domain events for async processing
4. Add domain validation logic
5. Implement domain-specific error handling
6. Add domain tests for each domain module
7. Document domain-specific business rules

## Development Notes

- This is a Rust + Axum backend project
- Uses PostgreSQL as primary database
- Uses ScyllaDB for time-series data storage
- Uses Redis for caching
- Uses vector database for AI-related data storage
- Follows DDD pattern with domain module pattern
- Each domain has its own module with mod.rs, schema.rs, domain.rs, api.rs structure
- Most domains are currently in placeholder state awaiting full implementation
- Domain logic is minimal - only basic CRUD operations are implemented
- Most domains need full domain logic implementation for business rules and validation