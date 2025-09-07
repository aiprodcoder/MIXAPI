-- ====================================================================
-- Token使用日志表迁移脚本 (用于访问频率限制)
-- ====================================================================
-- 版本: v1.0
-- 创建日期: 2025-08-25
-- 描述: 创建token_usage_logs表，用于记录令牌使用情况以支持访问频率限制功能
-- ====================================================================

-- ====================================================================
-- MySQL 语法
-- ====================================================================
CREATE TABLE IF NOT EXISTS token_usage_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    token_id INT NOT NULL COMMENT '令牌ID',
    created_at BIGINT NOT NULL COMMENT '创建时间戳',
    INDEX idx_token_created (token_id, created_at),
    INDEX idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token使用日志表(用于访问频率限制)';

-- ====================================================================
-- PostgreSQL 语法 (如果使用PostgreSQL，取消下面注释)
-- ====================================================================
/*
CREATE TABLE IF NOT EXISTS token_usage_logs (
    id SERIAL PRIMARY KEY,
    token_id INTEGER NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_token_usage_logs_token_created ON token_usage_logs(token_id, created_at);
CREATE INDEX IF NOT EXISTS idx_token_usage_logs_created ON token_usage_logs(created_at);

COMMENT ON TABLE token_usage_logs IS 'Token使用日志表(用于访问频率限制)';
COMMENT ON COLUMN token_usage_logs.token_id IS '令牌ID';
COMMENT ON COLUMN token_usage_logs.created_at IS '创建时间戳';
*/

-- ====================================================================
-- SQLite 语法 (如果使用SQLite，取消下面注释)
-- ====================================================================
/*
CREATE TABLE IF NOT EXISTS token_usage_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_token_usage_logs_token_created ON token_usage_logs(token_id, created_at);
CREATE INDEX IF NOT EXISTS idx_token_usage_logs_created ON token_usage_logs(created_at);
*/

-- ====================================================================
-- 清理旧数据的脚本 (可选，定期执行以减少存储空间)
-- ====================================================================
/*
-- 删除30天前的记录 (建议设置定时任务执行)
DELETE FROM token_usage_logs WHERE created_at < UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY));

-- PostgreSQL版本
-- DELETE FROM token_usage_logs WHERE created_at < EXTRACT(EPOCH FROM NOW() - INTERVAL '30 days');

-- SQLite版本
-- DELETE FROM token_usage_logs WHERE created_at < strftime('%s', 'now', '-30 days');
*/

-- ====================================================================
-- 使用说明
-- ====================================================================
/*
此表用于记录令牌的每次使用，支持访问频率限制功能：

1. 表结构：
   - id: 主键，自增
   - token_id: 令牌ID，关联tokens表
   - created_at: 创建时间戳

2. 索引：
   - idx_token_created: 复合索引，用于快速查询某个令牌在特定时间范围内的使用次数
   - idx_created: 时间索引，用于定期清理旧数据

3. 使用场景：
   - 分钟级限制：查询当前分钟内的记录数
   - 日级限制：查询当天内的记录数
   - 定期清理：删除过期记录以节省存储空间

4. 性能考虑：
   - 该表会频繁插入，但查询相对较少
   - 索引设计优化了查询性能
   - 建议定期清理旧数据
*/