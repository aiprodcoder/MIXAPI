-- 用量日统计汇总表迁移脚本
-- 创建日期: 2025-08-25
-- 功能: 创建用量统计汇总表，用于存储按日期、令牌、模型分组的统计数据

-- MySQL 语法
CREATE TABLE IF NOT EXISTS usage_statistics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(10) NOT NULL COMMENT '统计日期(YYYY-MM-DD)',
    token_id INT NOT NULL COMMENT '令牌ID',
    token_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '令牌名称',
    model_name VARCHAR(255) NOT NULL COMMENT '模型名称',
    total_requests INT NOT NULL DEFAULT 0 COMMENT '总请求次数',
    successful_requests INT NOT NULL DEFAULT 0 COMMENT '成功请求次数',
    failed_requests INT NOT NULL DEFAULT 0 COMMENT '失败请求次数',
    total_tokens INT NOT NULL DEFAULT 0 COMMENT '总Token消耗',
    prompt_tokens INT NOT NULL DEFAULT 0 COMMENT '提示Token数',
    completion_tokens INT NOT NULL DEFAULT 0 COMMENT '完成Token数',
    total_quota INT NOT NULL DEFAULT 0 COMMENT '总额度消耗',
    created_time BIGINT NOT NULL COMMENT '创建时间戳',
    updated_time BIGINT NOT NULL COMMENT '更新时间戳',
    INDEX idx_date (date),
    INDEX idx_token_id (token_id),
    INDEX idx_model_name (model_name),
    INDEX idx_date_token_model (date, token_id, model_name),
    UNIQUE KEY uk_date_token_model (date, token_id, model_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用量统计汇总表';

-- PostgreSQL 语法 (如果使用PostgreSQL)
/*
CREATE TABLE IF NOT EXISTS usage_statistics (
    id SERIAL PRIMARY KEY,
    date VARCHAR(10) NOT NULL,
    token_id INTEGER NOT NULL,
    token_name VARCHAR(255) NOT NULL DEFAULT '',
    model_name VARCHAR(255) NOT NULL,
    total_requests INTEGER NOT NULL DEFAULT 0,
    successful_requests INTEGER NOT NULL DEFAULT 0,
    failed_requests INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    prompt_tokens INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_quota INTEGER NOT NULL DEFAULT 0,
    created_time BIGINT NOT NULL,
    updated_time BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_usage_statistics_date ON usage_statistics(date);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_token_id ON usage_statistics(token_id);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_model_name ON usage_statistics(model_name);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_date_token_model ON usage_statistics(date, token_id, model_name);
CREATE UNIQUE INDEX IF NOT EXISTS uk_usage_statistics_date_token_model ON usage_statistics(date, token_id, model_name);

COMMENT ON TABLE usage_statistics IS '用量统计汇总表';
COMMENT ON COLUMN usage_statistics.date IS '统计日期(YYYY-MM-DD)';
COMMENT ON COLUMN usage_statistics.token_id IS '令牌ID';
COMMENT ON COLUMN usage_statistics.token_name IS '令牌名称';
COMMENT ON COLUMN usage_statistics.model_name IS '模型名称';
COMMENT ON COLUMN usage_statistics.total_requests IS '总请求次数';
COMMENT ON COLUMN usage_statistics.successful_requests IS '成功请求次数';
COMMENT ON COLUMN usage_statistics.failed_requests IS '失败请求次数';
COMMENT ON COLUMN usage_statistics.total_tokens IS '总Token消耗';
COMMENT ON COLUMN usage_statistics.prompt_tokens IS '提示Token数';
COMMENT ON COLUMN usage_statistics.completion_tokens IS '完成Token数';
COMMENT ON COLUMN usage_statistics.total_quota IS '总额度消耗';
COMMENT ON COLUMN usage_statistics.created_time IS '创建时间戳';
COMMENT ON COLUMN usage_statistics.updated_time IS '更新时间戳';
*/

-- SQLite 语法 (如果使用SQLite)
/*
CREATE TABLE IF NOT EXISTS usage_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    token_id INTEGER NOT NULL,
    token_name TEXT NOT NULL DEFAULT '',
    model_name TEXT NOT NULL,
    total_requests INTEGER NOT NULL DEFAULT 0,
    successful_requests INTEGER NOT NULL DEFAULT 0,
    failed_requests INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    prompt_tokens INTEGER NOT NULL DEFAULT 0,
    completion_tokens INTEGER NOT NULL DEFAULT 0,
    total_quota INTEGER NOT NULL DEFAULT 0,
    created_time INTEGER NOT NULL,
    updated_time INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_usage_statistics_date ON usage_statistics(date);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_token_id ON usage_statistics(token_id);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_model_name ON usage_statistics(model_name);
CREATE INDEX IF NOT EXISTS idx_usage_statistics_date_token_model ON usage_statistics(date, token_id, model_name);
CREATE UNIQUE INDEX IF NOT EXISTS uk_usage_statistics_date_token_model ON usage_statistics(date, token_id, model_name);
*/