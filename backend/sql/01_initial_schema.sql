-- Draft 内容草稿表
CREATE TABLE IF NOT EXISTS drafts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content LONGTEXT NOT NULL,
    content_type VARCHAR(20) DEFAULT 'text',
    status VARCHAR(20) DEFAULT 'draft', -- draft, archived
    tags JSON,
    metadata JSON,
    last_edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    saved_version_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_edited (user_id, last_edited_at DESC),
    INDEX idx_user_status (user_id, status),
    UNIQUE KEY uk_user_title (user_id, title)
);

-- Content 正式发布的内容表
CREATE TABLE IF NOT EXISTS contents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    body LONGTEXT NOT NULL,
    content_type VARCHAR(20),
    status VARCHAR(20),
    tags JSON,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_created (user_id, created_at DESC),
    INDEX idx_user_status (user_id, status),
    FULLTEXT INDEX ft_search (title, body)
);

-- ContentVersion 版本历史表
CREATE TABLE IF NOT EXISTS content_versions (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL,
    draft_id BIGINT,
    version_num INT NOT NULL,
    title VARCHAR(255),
    content LONGTEXT,
    change_summary VARCHAR(500),
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
    FOREIGN KEY (draft_id) REFERENCES drafts(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_content_version (content_id, version_num)
);

-- Tag 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(20),
    count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_name (user_id, name),
    INDEX idx_user (user_id)
);

-- Category 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    order_num INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_order (user_id, order_num)
);

-- PublishTask 发布任务表
CREATE TABLE IF NOT EXISTS publish_tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    content_id BIGINT NOT NULL,
    draft_id BIGINT,
    target_platforms JSON NOT NULL, -- ["juejin", "csdn"]
    status VARCHAR(20) DEFAULT 'pending', -- pending, publishing, published, failed
    scheduled_at TIMESTAMP,
    published_at TIMESTAMP,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    error_message TEXT,
    platform_results JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
    FOREIGN KEY (draft_id) REFERENCES drafts(id) ON DELETE SET NULL,
    INDEX idx_user_scheduled (user_id, scheduled_at),
    INDEX idx_user_status (user_id, status),
    INDEX idx_content (content_id)
);

-- PublishAnalytic 发布数据统计表
CREATE TABLE IF NOT EXISTS publish_analytics (
    id BIGSERIAL PRIMARY KEY,
    publish_task_id BIGINT NOT NULL,
    content_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    platform VARCHAR(50),
    post_id VARCHAR(255),
    title VARCHAR(255),
    post_url TEXT,
    metrics JSON,
    last_synced_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (publish_task_id) REFERENCES publish_tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_created (user_id, created_at DESC),
    INDEX idx_user_platform (user_id, platform),
    INDEX idx_publish_task (publish_task_id)
);

-- PromptTemplate AI提示词模板表
CREATE TABLE IF NOT EXISTS prompt_templates (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT, -- NULL表示系统模板
    category VARCHAR(100),
    sub_category VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    prompt_text LONGTEXT,
    variables JSON,
    example_input JSON,
    example_output TEXT,
    usage_count INT DEFAULT 0,
    rating FLOAT DEFAULT 0,
    tags JSON,
    status VARCHAR(20) DEFAULT 'active', -- active, archived
    version INT DEFAULT 1,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_category (category, sub_category),
    INDEX idx_user (user_id),
    FULLTEXT INDEX ft_search (name, description)
);

-- 假设 users 表已存在，如果不存在需要创建
-- CREATE TABLE IF NOT EXISTS users (
--     id BIGSERIAL PRIMARY KEY,
--     username VARCHAR(255) UNIQUE NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
