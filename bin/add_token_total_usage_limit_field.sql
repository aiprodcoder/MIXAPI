-- Token总使用次数限制字段迁移脚本
-- 添加日期: 2025-09-11
-- 功能: 为tokens表添加总使用次数限制字段

-- MySQL 语法
ALTER TABLE tokens ADD COLUMN total_usage_limit INT DEFAULT NULL COMMENT '总使用次数限制，NULL表示不限制';

-- PostgreSQL 语法 (如果使用PostgreSQL数据库)
-- ALTER TABLE tokens ADD COLUMN total_usage_limit INTEGER DEFAULT NULL;
-- COMMENT ON COLUMN tokens.total_usage_limit IS '总使用次数限制，NULL表示不限制';

-- SQLite 语法 (如果使用SQLite数据库)
-- ALTER TABLE tokens ADD COLUMN total_usage_limit INTEGER DEFAULT NULL;