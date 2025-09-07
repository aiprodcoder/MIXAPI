-- 简化版Token使用次数统计字段迁移脚本
-- 适用于MySQL数据库，可直接执行

-- 添加今日使用次数字段
ALTER TABLE tokens ADD COLUMN daily_usage_count INT NOT NULL DEFAULT 0 COMMENT '今日使用次数';

-- 添加总使用次数字段  
ALTER TABLE tokens ADD COLUMN total_usage_count INT NOT NULL DEFAULT 0 COMMENT '总使用次数';

-- 添加最后使用日期字段
ALTER TABLE tokens ADD COLUMN last_usage_date VARCHAR(10) NOT NULL DEFAULT '' COMMENT '最后使用日期(YYYY-MM-DD)';