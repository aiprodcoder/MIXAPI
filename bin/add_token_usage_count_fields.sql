-- Token使用次数统计字段迁移脚本
-- 添加日期: 2025-08-25
-- 功能: 为tokens表添加今日次数、总次数和最后使用日期字段

-- MySQL 语法
-- 添加今日使用次数字段
ALTER TABLE tokens ADD COLUMN daily_usage_count INT NOT NULL DEFAULT 0 COMMENT '今日使用次数';

-- 添加总使用次数字段  
ALTER TABLE tokens ADD COLUMN total_usage_count INT NOT NULL DEFAULT 0 COMMENT '总使用次数';

-- 添加最后使用日期字段
ALTER TABLE tokens ADD COLUMN last_usage_date VARCHAR(10) NOT NULL DEFAULT '' COMMENT '最后使用日期(YYYY-MM-DD)';

-- PostgreSQL 语法 (如果使用PostgreSQL数据库)
-- ALTER TABLE tokens ADD COLUMN daily_usage_count INTEGER NOT NULL DEFAULT 0;
-- ALTER TABLE tokens ADD COLUMN total_usage_count INTEGER NOT NULL DEFAULT 0; 
-- ALTER TABLE tokens ADD COLUMN last_usage_date VARCHAR(10) NOT NULL DEFAULT '';

-- COMMENT ON COLUMN tokens.daily_usage_count IS '今日使用次数';
-- COMMENT ON COLUMN tokens.total_usage_count IS '总使用次数';
-- COMMENT ON COLUMN tokens.last_usage_date IS '最后使用日期(YYYY-MM-DD)';

-- SQLite 语法 (如果使用SQLite数据库)
-- ALTER TABLE tokens ADD COLUMN daily_usage_count INTEGER NOT NULL DEFAULT 0;
-- ALTER TABLE tokens ADD COLUMN total_usage_count INTEGER NOT NULL DEFAULT 0;
-- ALTER TABLE tokens ADD COLUMN last_usage_date TEXT NOT NULL DEFAULT '';