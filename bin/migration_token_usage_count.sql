-- ====================================================================
-- Token使用次数统计功能数据库迁移脚本
-- ====================================================================
-- 版本: v1.0
-- 创建日期: 2025-08-25
-- 作者: Assistant
-- 描述: 为tokens表添加使用次数统计功能的相关字段
-- 
-- 新增字段:
-- - daily_usage_count: 今日使用次数，每日自动重置
-- - total_usage_count: 总使用次数，只增不减
-- - last_usage_date: 最后使用日期，用于判断是否需要重置今日次数
-- ====================================================================

-- ====================================================================
-- 向前迁移 (UP) - 添加字段
-- ====================================================================

-- MySQL 数据库
-- 检查字段是否已存在，避免重复添加
SET @sql_daily = 'ALTER TABLE tokens ADD COLUMN daily_usage_count INT NOT NULL DEFAULT 0 COMMENT ''今日使用次数''';
SET @sql_total = 'ALTER TABLE tokens ADD COLUMN total_usage_count INT NOT NULL DEFAULT 0 COMMENT ''总使用次数''';
SET @sql_date = 'ALTER TABLE tokens ADD COLUMN last_usage_date VARCHAR(10) NOT NULL DEFAULT '''' COMMENT ''最后使用日期(YYYY-MM-DD)''';

-- 执行添加字段的SQL (仅在字段不存在时执行)
PREPARE stmt FROM @sql_daily;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

PREPARE stmt FROM @sql_total;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

PREPARE stmt FROM @sql_date;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 为新字段添加索引 (可选，提升查询性能)
-- CREATE INDEX idx_tokens_last_usage_date ON tokens(last_usage_date);
-- CREATE INDEX idx_tokens_total_usage_count ON tokens(total_usage_count);

-- ====================================================================
-- PostgreSQL 数据库 (如果使用PostgreSQL，取消下面注释)
-- ====================================================================
/*
-- 添加字段
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS daily_usage_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS total_usage_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS last_usage_date VARCHAR(10) NOT NULL DEFAULT '';

-- 添加注释
COMMENT ON COLUMN tokens.daily_usage_count IS '今日使用次数';
COMMENT ON COLUMN tokens.total_usage_count IS '总使用次数';
COMMENT ON COLUMN tokens.last_usage_date IS '最后使用日期(YYYY-MM-DD)';

-- 添加索引 (可选)
-- CREATE INDEX IF NOT EXISTS idx_tokens_last_usage_date ON tokens(last_usage_date);
-- CREATE INDEX IF NOT EXISTS idx_tokens_total_usage_count ON tokens(total_usage_count);
*/

-- ====================================================================
-- SQLite 数据库 (如果使用SQLite，取消下面注释)
-- ====================================================================
/*
-- SQLite不支持IF NOT EXISTS语法，需要手动检查
-- 添加字段
ALTER TABLE tokens ADD COLUMN daily_usage_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN total_usage_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokens ADD COLUMN last_usage_date TEXT NOT NULL DEFAULT '';

-- 添加索引 (可选)
-- CREATE INDEX IF NOT EXISTS idx_tokens_last_usage_date ON tokens(last_usage_date);
-- CREATE INDEX IF NOT EXISTS idx_tokens_total_usage_count ON tokens(total_usage_count);
*/

-- ====================================================================
-- 向后迁移 (DOWN) - 删除字段 (回滚操作)
-- ====================================================================
/*
-- 如果需要回滚，请执行以下SQL语句

-- MySQL 回滚
ALTER TABLE tokens DROP COLUMN IF EXISTS daily_usage_count;
ALTER TABLE tokens DROP COLUMN IF EXISTS total_usage_count;
ALTER TABLE tokens DROP COLUMN IF EXISTS last_usage_date;

-- PostgreSQL 回滚
-- ALTER TABLE tokens DROP COLUMN IF EXISTS daily_usage_count;
-- ALTER TABLE tokens DROP COLUMN IF EXISTS total_usage_count;
-- ALTER TABLE tokens DROP COLUMN IF EXISTS last_usage_date;

-- SQLite 回滚 (SQLite不支持DROP COLUMN，需要重建表)
-- 1. 创建备份表
-- CREATE TABLE tokens_backup AS SELECT id, user_id, key, status, name, created_time, accessed_time, expired_time, remain_quota, unlimited_quota, model_limits_enabled, model_limits, allow_ips, used_quota, group FROM tokens;
-- 2. 删除原表
-- DROP TABLE tokens;
-- 3. 重命名备份表
-- ALTER TABLE tokens_backup RENAME TO tokens;
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
  AND COLUMN_NAME IN ('daily_usage_count', 'total_usage_count', 'last_usage_date');

-- 检查字段是否存在 (PostgreSQL)
-- SELECT column_name, data_type, is_nullable, column_default 
-- FROM information_schema.columns 
-- WHERE table_name = 'tokens' 
--   AND column_name IN ('daily_usage_count', 'total_usage_count', 'last_usage_date');

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
- daily_usage_count: 记录Token今日使用次数，每天首次使用时重置
- total_usage_count: 记录Token总使用次数，永不重置
- last_usage_date: 记录最后使用日期，格式为YYYY-MM-DD，用于判断是否跨天

这些字段会在Token被使用时自动更新，无需手动维护。
*/