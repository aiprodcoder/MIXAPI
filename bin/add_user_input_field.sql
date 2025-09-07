-- 为logs表添加user_input字段来记录用户输入内容
-- 适用于 MySQL、SQLite、PostgreSQL 数据库

-- MySQL 语法
-- ALTER TABLE logs ADD COLUMN user_input TEXT COMMENT '用户输入内容';

-- SQLite 语法
-- ALTER TABLE logs ADD COLUMN user_input TEXT;

-- PostgreSQL 语法
-- ALTER TABLE logs ADD COLUMN user_input TEXT;

-- 通用语法（兼容多数据库）
ALTER TABLE logs ADD COLUMN user_input TEXT;

-- 为新字段添加索引（可选，如果需要搜索用户输入内容的话）
-- CREATE INDEX idx_logs_user_input ON logs(user_input(100));

-- 说明：
-- 1. user_input 字段用于存储用户通过API发送的实际输入内容（如messages内容）
-- 2. 原有的 content 字段继续用于存储计费相关信息
-- 3. 该字段主要记录消费类型(type=2)和错误类型(type=5)日志的用户输入
-- 4. 字段类型使用TEXT以支持长文本内容