# Docker PostgreSQL Deployment Guide

## Overview

This guide covers PostgreSQL deployment using Docker for the Magic project.

## Files Created

- `docker-compose.yml` - Docker Compose configuration
- `.env.example` - Environment variable template

## Quick Start

### 1. Start PostgreSQL

```bash
docker-compose up -d
```

### 2. Verify PostgreSQL is Running

```bash
docker-compose ps
docker-compose logs postgres
```

### 3. Connect to PostgreSQL

```bash
# Using psql client
docker exec -it postgres_db psql -U root -d vc_backend

# Using pgAdmin (optional - add pgAdmin service to docker-compose.yml)
# Access at http://localhost:8080
```

### 4. Stop PostgreSQL

```bash
docker-compose down

# To remove volumes (data will be lost)
docker-compose down -v
```

## Database Configuration

| Setting | Value |
|---------|-------|
| Host | localhost or postgres (when using Docker networking) |
| Port | 5432 |
| User | root |
| Password | pass |
| Database | vc_backend |

## Production Considerations

### Environment Variables

For production, use environment files:

```bash
docker run -d \
  --name postgres_db \
  -e POSTGRES_USER=${POSTGRES_USER} \
  -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
  -e POSTGRES_DB=${POSTGRES_DB} \
  -p 5432:5432 \
  -v $(pwd)/data:/var/lib/postgresql/data \
  -v $(pwd)/pg_hba.conf:/etc/postgresql/pg_hba.conf \
  -v $(pwd)/postgresql.conf:/etc/postgresql/postgresql.conf \
  postgres:16-alpine
```

### Health Checks

```bash
# Container health status
docker inspect --format='{{.State.Health.Status}}' postgres_db
```

## Troubleshooting

### Common Issues

1. **Port already in use

```bash
docker stop $(docker ps -q)  # Stop all containers
lsof -i :5432  # Check what's using port 5432
```

2. **Connection refused

```bash
docker-compose logs postgres
docker-compose restart postgres
```

3. **Permission issues

```bash
# Fix volume permissions
docker-compose down
rm -rf postgres_data
docker-compose up -d
```

### Database Migration

```bash
# From your Rust backend
cargo run -- run-migrations

# Or using psql directly
psql postgresql://root:pass@localhost:5432/vc_backend -f schema.sql
```

## Backup and Restore

### Backup

```bash
# Using docker exec with pg_dump
docker exec postgres_db pg_dump -U root vc_backend > backup.sql

# With gzip compression
docker exec postgres_db pg_dump -U root vc_backend | gzip > backup.sql.gz
```

### Restore

```bash
# Using docker exec with psql
docker exec -i postgres_db psql -U root vc_backup < backup.sql

# From compressed backup
gunzip -c backup.sql.gz | docker exec -i postgres_db psql -U root vc_backup
```

## Update PostgreSQL Version

```bash
# Edit docker-compose.yml and change image tag:
# image: postgres:16-alpine  # Update version
# Then:
docker-compose pull
docker-compose up -d
```

## Performance Tuning

Basic settings in `postgresql.conf:

```
shared_buffers = 256MB
work_mem = 16MB
maintenance_work_mem = 64MB
effective_cache_size = 1GB
```

## Integration with Rust Backend

Update `backend/src/config/config.rs:

```rust
impl AppConfig {
    pub fn load() -> Result<Self> {
        Ok(AppConfig {
            database: DatabaseConfig {
                host: "localhost".to_string(),
                user: "root".to_string(),
                pass: "pass".to_string(),
                db_name: "vc_backend".to_string(),
            },
        })
    }
}
```

Or use environment variables for production:

```rust
let database_host = std::env::var("POSTGRES_HOST").unwrap_or_else(|_| "localhost".to_string());
```

## Environment-Specific Configuration

### Development (docker-compose.yml)

```yaml
database:
  host: localhost
  port: 5432
```

### Production (environment variables)

```
POSTGRES_HOST=postgres.production.com
POSTGRES_USER=production_user
POSTGRES_PASSWORD=secure_password
POSTGRES_DB=production_db
```

## Monitoring

### Container Metrics

```bash
# CPU and memory usage
docker stats postgres_db

# Resource limits in docker-compose.yml
resources:
  limits:
    cpus: '1.0'
    memory: 1G
  reservations:
    cpus: '0.5'
    memory: 512M
```

### PostgreSQL Logs

```bash
# Container logs
docker-compose logs -f postgres

# PostgreSQL-specific logs
docker exec postgres_db tail -f /var/lib/postgresql/data/log/postgresql.log
```

## Security Considerations

1. Use environment files for sensitive data
2. Enable SSL connections
3. Use Docker secrets for production secrets
4. Limit container resources
5. Keep PostgreSQL updated

## References

- [Docker Hub PostgreSQL Image](https://hub.docker.com/_/postgres)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
