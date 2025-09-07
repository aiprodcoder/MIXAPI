-- PostgreSQL数据库Token使用次数统计字段迁移脚本

-- 添加今日使用次数字段
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS daily_usage_count INTEGER NOT NULL DEFAULT 0;

-- 添加总使用次数字段
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS total_usage_count INTEGER NOT NULL DEFAULT 0;

-- 添加最后使用日期字段
ALTER TABLE tokens ADD COLUMN IF NOT EXISTS last_usage_date VARCHAR(10) NOT NULL DEFAULT '';

-- 添加字段注释
COMMENT ON COLUMN tokens.daily_usage_count IS '今日使用次数';
COMMENT ON COLUMN tokens.total_usage_count IS '总使用次数';
COMMENT ON COLUMN tokens.last_usage_date IS '最后使用日期(YYYY-MM-DD)';