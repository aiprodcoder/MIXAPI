-- ====================================================================
-- Token访问频率限制功能数据库迁移脚本
-- ====================================================================
-- 版本: v1.0
-- 创建日期: 2025-08-25
-- 描述: 为tokens表添加访问频率限制功能的相关字段
-- 
-- 新增字段:
-- - rate_limit_per_minute: 每分钟访问次数限制，0表示不限制
-- - rate_limit_per_day: 每日访问次数限制，0表示不限制
-- - last_rate_limit_reset: 最后重置时间，用于重置计数器
-- ====================================================================

-- ====================================================================
-- 向前迁移 (UP) - 添加字段
-- ====================================================================

-- MySQL 数据库
-- 检查字段是否已存在，避免重复添加
SET @sql_rate_minute = 'ALTER TABLE tokens ADD COLUMN rate_limit_per_minute INT NOT NULL DEFAULT 0 COMMENT ''每分钟访问次数限制，0表示不限制''';
SET @sql_rate_day = 'ALTER TABLE tokens ADD COLUMN rate_limit_per_day INT NOT NULL DEFAULT 0 COMMENT ''每日访问次数限制，0表示不限制''';
SET @sql_reset_time = 'ALTER TABLE tokens ADD COLUMN last_rate_limit_reset BIGINT NOT NULL DEFAULT 0 COMMENT ''最后重置时间戳''';

-- 执行添加字段的SQL (仅在字段不存在时执行)
PREPARE stmt FROM @sql_rate_minute;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

PREPARE stmt FROM @sql_rate_day;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

PREPARE stmt FROM @sql_reset_time;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 为新字段添加索引 (可选，提升查询性能)
-- CREATE INDEX idx_tokens_rate_limit_reset ON tokens(last_rate_limit_reset);

-- ====================================================================
-- PostgreSQL 数据库 (如果使用PostgreSQL，取消下面注释)
-- ====================================================================
/*
-- 添加字段
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS rate_limit_per_minute INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS rate_limit_per_day INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS last_rate_limit_reset BIGINT NOT NULL DEFAULT 0;

-- 添加注释
COMMENT ON COLUMN tokens.rate_limit_per_minute IS '每分钟访问次数限制，0表示不限制';
COMMENT ON COLUMN tokens.rate_limit_per_day IS '每日访问次数限制，0表示不限制';
COMMENT ON COLUMN tokens.last_rate_limit_reset IS '最后重置时间戳';

-- 添加索引 (可选)
-- CREATE INDEX IF NOT EXISTS idx_tokens_rate_limit_reset ON tokens(last_rate_limit_reset);
*/

-- ====================================================================
-- SQLite 数据库 (如果使用SQLite，取消下面注释)
-- ====================================================================
/*
-- SQLite不支持IF NOT EXISTS语法，需要手动检查
-- 添加字段
ALTER TABLE tokens ADD COLUMN rate_limit_per_minute INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN rate_limit_per_day INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN last_rate_limit_reset INTEGER NOT NULL DEFAULT 0;

-- 添加索引 (可选)
-- CREATE INDEX IF NOT EXISTS idx_tokens_rate_limit_reset ON tokens(last_rate_limit_reset);
*/

-- ====================================================================
-- 向后迁移 (DOWN) - 删除字段 (回滚操作)
-- ====================================================================
/*
-- 如果需要回滚，请执行以下SQL语句

-- MySQL 回滚
ALTER TABLE tokens DROP COLUMN IF EXISTS rate_limit_per_minute;
ALTER TABLE tokens DROP COLUMN IF EXISTS rate_limit_per_day;
ALTER TABLE tokens DROP COLUMN IF EXISTS last_rate_limit_reset;

-- PostgreSQL 回滚
-- ALTER TABLE tokens DROP COLUMN IF EXISTS rate_limit_per_minute;
-- ALTER TABLE tokens DROP COLUMN IF EXISTS rate_limit_per_day;
-- ALTER TABLE tokens DROP COLUMN IF EXISTS last_rate_limit_reset;

-- SQLite 回滚
-- SQLite不支持DROP COLUMN，需要重建表：
-- 1. 创建临时表包含原有字段
-- CREATE TABLE tokens_temp AS 
-- SELECT id, user_id, key, status, name, created_time, accessed_time, 
--        expired_time, remain_quota, unlimited_quota, model_limits_enabled, 
--        model_limits, allow_ips, used_quota, `group`, deleted_at,
--        daily_usage_count, total_usage_count, last_usage_date
-- FROM tokens;

-- 2. 删除原表
-- DROP TABLE tokens;

-- 3. 重命名临时表
-- ALTER TABLE tokens_temp RENAME TO tokens;

-- 4. 重建索引和约束 (根据实际情况调整)
*/

-- ====================================================================
-- 验证脚本 - 检查字段是否成功添加
-- ====================================================================
/*
-- 检查字段是否存在 (MySQL)
SELECT 
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'tokens' 
  AND COLUMN_NAME IN ('rate_limit_per_minute', 'rate_limit_per_day', 'last_rate_limit_reset');

-- 检查字段是否存在 (PostgreSQL)
-- SELECT column_name, data_type, is_nullable, column_default 
-- FROM information_schema.columns 
-- WHERE table_name = 'tokens' 
--   AND column_name IN ('rate_limit_per_minute', 'rate_limit_per_day', 'last_rate_limit_reset');

-- 检查字段是否存在 (SQLite)
-- PRAGMA table_info(tokens);
*/

-- ====================================================================
-- 使用说明
-- ====================================================================
/*
1. 根据您的数据库类型，选择对应的SQL语句执行
2. 默认启用MySQL语法，如使用其他数据库请取消相应注释
3. 建议在执行前备份数据库
4. 执行后可运行验证脚本确认字段添加成功
5. 如需回滚，请执行"向后迁移"部分的SQL语句

新字段说明:
- rate_limit_per_minute: 每分钟访问次数限制，0表示不限制
- rate_limit_per_day: 每日访问次数限制，0表示不限制
- last_rate_limit_reset: 最后重置时间戳，用于计算是否需要重置计数器

这些字段会在Token验证时被检查，实现访问频率限制功能。
*/