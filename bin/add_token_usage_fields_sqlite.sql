-- SQLite数据库Token使用次数统计字段迁移脚本

-- 添加今日使用次数字段
ALTER TABLE tokens ADD COLUMN daily_usage_count INTEGER NOT NULL DEFAULT 0;

-- 添加总使用次数字段
ALTER TABLE tokens ADD COLUMN total_usage_count INTEGER NOT NULL DEFAULT 0;

-- 添加最后使用日期字段
ALTER TABLE tokens ADD COLUMN last_usage_date TEXT NOT NULL DEFAULT '';