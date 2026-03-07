-- Square posts table
CREATE TABLE IF NOT EXISTS square_posts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  draft_id BIGINT NOT NULL COMMENT '关联草稿',
  user_id BIGINT NOT NULL COMMENT '发布者用户ID',
  title VARCHAR(255) NOT NULL COMMENT '帖子标题',
  preview_text TEXT COMMENT '预览文本',
  domain VARCHAR(50) COMMENT '内容领域分类',
  tags JSON COMMENT '标签数组',
  cover_image VARCHAR(500) COMMENT '封面图URL',
  likes_count INT DEFAULT 0 COMMENT '点赞数',
  comments_count INT DEFAULT 0 COMMENT '评论数',
  collects_count INT DEFAULT 0 COMMENT '收藏数',
  views_count INT DEFAULT 0 COMMENT '浏览数',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  KEY idx_draft_id (draft_id),
  KEY idx_user_id (user_id),
  KEY idx_domain (domain),
  KEY idx_likes_count (likes_count),
  KEY idx_collects_count (collects_count),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='广场内容';

-- Square comments table
CREATE TABLE IF NOT EXISTS square_comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL COMMENT '帖子ID',
  user_id BIGINT NOT NULL COMMENT '评论者用户ID',
  content TEXT NOT NULL COMMENT '评论内容',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  KEY idx_post_id (post_id),
  KEY idx_user_id (user_id),
  KEY idx_created_at (created_at),
  FOREIGN KEY (post_id) REFERENCES square_posts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='广场评论';

-- Square likes table (点赞记录)
CREATE TABLE IF NOT EXISTS square_likes (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL COMMENT '帖子ID',
  user_id BIGINT NOT NULL COMMENT '点赞用户ID',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  UNIQUE KEY unique_like (post_id, user_id) COMMENT '同一用户对同一帖子只能点赞一次',
  KEY idx_post_id (post_id),
  KEY idx_user_id (user_id),
  FOREIGN KEY (post_id) REFERENCES square_posts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='广场点赞记录';

-- Square collects table (收藏记录)
CREATE TABLE IF NOT EXISTS square_collects (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  post_id BIGINT NOT NULL COMMENT '帖子ID',
  user_id BIGINT NOT NULL COMMENT '收藏用户ID',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  UNIQUE KEY unique_collect (post_id, user_id) COMMENT '同一用户对同一帖子只能收藏一次',
  KEY idx_post_id (post_id),
  KEY idx_user_id (user_id),
  FOREIGN KEY (post_id) REFERENCES square_posts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='广场收藏记录';

-- 添加 domain 字段到 drafts 表
ALTER TABLE drafts ADD COLUMN domain VARCHAR(50) DEFAULT NULL COMMENT '内容领域分类' AFTER metadata;
ALTER TABLE drafts ADD KEY idx_domain (domain);
