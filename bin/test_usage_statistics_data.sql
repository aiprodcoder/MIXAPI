-- 测试用量统计功能的SQL脚本
-- 执行此脚本可以插入一些测试数据来验证用量统计功能

-- 插入测试数据到usage_statistics表
INSERT INTO usage_statistics (
    date, token_id, token_name, model_name, 
    total_requests, successful_requests, failed_requests,
    total_tokens, prompt_tokens, completion_tokens, total_quota,
    created_time, updated_time
) VALUES 
-- 今天的数据
('2025-08-25', 1, 'test-token-1', 'gpt-3.5-turbo', 10, 8, 2, 5000, 3000, 2000, 1000, 1724568000, 1724568000),
('2025-08-25', 1, 'test-token-1', 'gpt-4', 5, 5, 0, 8000, 5000, 3000, 2000, 1724568000, 1724568000),
('2025-08-25', 2, 'test-token-2', 'gpt-3.5-turbo', 15, 12, 3, 7500, 4500, 3000, 1500, 1724568000, 1724568000),

-- 昨天的数据  
('2025-08-24', 1, 'test-token-1', 'gpt-3.5-turbo', 20, 18, 2, 10000, 6000, 4000, 2000, 1724481600, 1724481600),
('2025-08-24', 1, 'test-token-1', 'gpt-4', 8, 7, 1, 12000, 7000, 5000, 3000, 1724481600, 1724481600),
('2025-08-24', 2, 'test-token-2', 'gpt-3.5-turbo', 12, 10, 2, 6000, 3600, 2400, 1200, 1724481600, 1724481600),

-- 前天的数据
('2025-08-23', 1, 'test-token-1', 'gpt-3.5-turbo', 25, 22, 3, 12500, 7500, 5000, 2500, 1724395200, 1724395200),
('2025-08-23', 2, 'test-token-2', 'gpt-4', 6, 5, 1, 9000, 5400, 3600, 1800, 1724395200, 1724395200),

-- 一周前的数据
('2025-08-18', 1, 'test-token-1', 'gpt-3.5-turbo', 30, 28, 2, 15000, 9000, 6000, 3000, 1723939200, 1723939200),
('2025-08-18', 2, 'test-token-2', 'gpt-3.5-turbo', 18, 15, 3, 9000, 5400, 3600, 1800, 1723939200, 1723939200);

-- 查询插入的数据验证
SELECT 
    date,
    token_name,
    model_name,
    total_requests,
    successful_requests,
    failed_requests,
    total_tokens,
    total_quota
FROM usage_statistics 
ORDER BY date DESC, token_id ASC, model_name ASC;

-- 查询汇总统计
SELECT 
    SUM(total_requests) as total_requests,
    SUM(successful_requests) as successful_requests,
    SUM(failed_requests) as failed_requests,
    ROUND(SUM(successful_requests) * 100.0 / SUM(total_requests), 2) as success_rate,
    SUM(total_tokens) as total_tokens,
    SUM(total_quota) as total_quota
FROM usage_statistics 
WHERE date >= '2025-08-18';