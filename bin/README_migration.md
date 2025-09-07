# Token使用次数统计字段 - SQL迁移脚本使用说明

## 文件列表

本目录包含了以下SQL迁移脚本文件：

### 1. 完整迁移脚本
- **`migration_token_usage_count.sql`** - 完整的迁移脚本，包含所有数据库类型的语法和回滚操作

### 2. 数据库专用脚本
- **`add_token_usage_fields_mysql.sql`** - MySQL专用简化脚本
- **`add_token_usage_fields_sqlite.sql`** - SQLite专用简化脚本  
- **`add_token_usage_fields_postgresql.sql`** - PostgreSQL专用简化脚本

### 3. 通用脚本
- **`add_token_usage_count_fields.sql`** - 包含所有数据库语法的通用脚本

## 新增字段说明

| 字段名 | 数据类型 | 默认值 | 说明 |
|--------|----------|--------|------|
| `daily_usage_count` | INT/INTEGER | 0 | 今日使用次数，每日自动重置 |
| `total_usage_count` | INT/INTEGER | 0 | 总使用次数，只增不减 |
| `last_usage_date` | VARCHAR(10)/TEXT | '' | 最后使用日期(YYYY-MM-DD格式) |

## 使用方法

### 自动迁移（推荐）
如果您的项目使用了GORM自动迁移，则无需手动执行SQL脚本：
```bash
# 启动程序时会自动迁移
go run main.go
```

### 手动迁移
如果需要手动执行SQL迁移，请根据您的数据库类型选择对应脚本：

#### MySQL数据库
```bash
mysql -u root -p your_database < add_token_usage_fields_mysql.sql
```

#### SQLite数据库  
```bash
sqlite3 your_database.db < add_token_usage_fields_sqlite.sql
```

#### PostgreSQL数据库
```bash
psql -U username -d your_database -f add_token_usage_fields_postgresql.sql
```

## 迁移前注意事项

1. **备份数据库**: 执行迁移前请务必备份您的数据库
2. **检查权限**: 确保数据库用户有ALTER TABLE权限
3. **停止服务**: 建议在迁移期间停止API服务
4. **测试环境**: 建议先在测试环境验证迁移脚本

## 验证迁移结果

### MySQL验证
```sql
DESCRIBE tokens;
-- 或
SHOW COLUMNS FROM tokens LIKE '%usage%';
```

### SQLite验证  
```sql
PRAGMA table_info(tokens);
```

### PostgreSQL验证
```sql
\d tokens
-- 或
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'tokens' 
  AND column_name LIKE '%usage%';
```

## 回滚操作

如果需要回滚迁移，请参考 `migration_token_usage_count.sql` 文件中的回滚部分，或执行以下SQL：

### MySQL/PostgreSQL回滚
```sql
ALTER TABLE tokens DROP COLUMN daily_usage_count;
ALTER TABLE tokens DROP COLUMN total_usage_count;  
ALTER TABLE tokens DROP COLUMN last_usage_date;
```

### SQLite回滚
SQLite不支持DROP COLUMN，需要重建表：
```sql
-- 1. 创建临时表包含原有字段
CREATE TABLE tokens_temp AS 
SELECT id, user_id, key, status, name, created_time, accessed_time, 
       expired_time, remain_quota, unlimited_quota, model_limits_enabled, 
       model_limits, allow_ips, used_quota, `group`, deleted_at 
FROM tokens;

-- 2. 删除原表
DROP TABLE tokens;

-- 3. 重命名临时表
ALTER TABLE tokens_temp RENAME TO tokens;

-- 4. 重建索引和约束 (根据实际情况调整)
```

## 故障排除

### 常见错误

1. **字段已存在错误**
   - 错误信息: `Duplicate column name` 或 `column already exists`
   - 解决方案: 字段已存在，无需重复添加

2. **权限不足错误**
   - 错误信息: `Access denied` 或 `permission denied`
   - 解决方案: 检查数据库用户权限

3. **表不存在错误**
   - 错误信息: `Table doesn't exist`
   - 解决方案: 确认表名正确，检查数据库选择

### 检查脚本执行状态
```sql
-- 检查字段是否添加成功
SELECT COUNT(*) as field_count
FROM information_schema.columns 
WHERE table_name = 'tokens' 
  AND column_name IN ('daily_usage_count', 'total_usage_count', 'last_usage_date');
-- 应该返回 3
```

## 技术支持

如果在迁移过程中遇到问题：
1. 查看数据库错误日志
2. 确认SQL语法与数据库版本兼容
3. 验证数据库连接和权限
4. 参考项目文档或提交issue