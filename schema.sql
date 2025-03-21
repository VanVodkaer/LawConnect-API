-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS LawConnectDB 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;
USE LawConnectDB;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱，必须唯一',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    role TINYINT NOT NULL DEFAULT 1 COMMENT '用户权限：1-普通用户，2-管理员'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建分类表，支持父分类
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    parent_id INT DEFAULT NULL COMMENT '父分类ID',
    CONSTRAINT fk_parent_category FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建文章表（包含分类ID属性），增加评论数量字段 comment_count
CREATE TABLE IF NOT EXISTS articles (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '文章ID',
    is_visible TINYINT NOT NULL DEFAULT 1 COMMENT '是否可见：0-不可见，1-可见',
    title VARCHAR(255) NOT NULL COMMENT '文章标题',
    content TEXT NOT NULL COMMENT '文章内容',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后编辑时间',
    likes INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    comment_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数量',
    user_id INT NOT NULL COMMENT '发布用户ID',
    category_id INT NOT NULL COMMENT '文章分类ID',
    INDEX idx_user_id (user_id),
    INDEX idx_category_id (category_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建评论表（增加了user_id字段）
CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '评论ID',
    article_id INT NOT NULL COMMENT '评论关联文章ID',
    content TEXT NOT NULL COMMENT '评论内容',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '评论发表时间',
    user_id INT NOT NULL COMMENT '发表评论的用户ID',
    is_visible TINYINT NOT NULL DEFAULT 0 COMMENT '是否可见：0-审核中，1-可见，2-审核未通过',
    likes INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    INDEX idx_article_id (article_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建文章点赞记录表（记录用户对文章的点赞，防止重复点赞）
CREATE TABLE IF NOT EXISTS article_likes (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '点赞记录ID',
    article_id INT NOT NULL COMMENT '被点赞的文章ID',
    user_id INT NOT NULL COMMENT '点赞用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    INDEX idx_article_user (article_id, user_id), -- 用于快速查询用户是否已点赞
    UNIQUE KEY uk_article_user (article_id, user_id), -- 确保一个用户只能给同一文章点赞一次
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建评论点赞记录表（记录用户对评论的点赞，防止重复点赞）
CREATE TABLE IF NOT EXISTS comment_likes (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '点赞记录ID',
    comment_id INT NOT NULL COMMENT '被点赞的评论ID',
    user_id INT NOT NULL COMMENT '点赞用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    INDEX idx_comment_user (comment_id, user_id), -- 用于快速查询用户是否已点赞
    UNIQUE KEY uk_comment_user (comment_id, user_id), -- 确保一个用户只能给同一评论点赞一次
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;