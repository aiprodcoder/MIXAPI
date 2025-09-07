package model

import (
	"errors"
	"one-api/common"
	"time"

	"gorm.io/gorm"
)

type UsageStatistics struct {
	Id                 int    `json:"id" gorm:"primaryKey"`
	Date               string `json:"date" gorm:"type:varchar(10);not null;index:idx_date;uniqueIndex:uk_date_token_model,composite:date"`
	TokenId            int    `json:"token_id" gorm:"not null;index:idx_token_id;uniqueIndex:uk_date_token_model,composite:token_id"`
	TokenName          string `json:"token_name" gorm:"type:varchar(255);not null;default:''"`
	ModelName          string `json:"model_name" gorm:"type:varchar(255);not null;index:idx_model_name;uniqueIndex:uk_date_token_model,composite:model_name"`
	TotalRequests      int    `json:"total_requests" gorm:"not null;default:0"`
	SuccessfulRequests int    `json:"successful_requests" gorm:"not null;default:0"`
	FailedRequests     int    `json:"failed_requests" gorm:"not null;default:0"`
	TotalTokens        int    `json:"total_tokens" gorm:"not null;default:0"`
	PromptTokens       int    `json:"prompt_tokens" gorm:"not null;default:0"`
	CompletionTokens   int    `json:"completion_tokens" gorm:"not null;default:0"`
	TotalQuota         int    `json:"total_quota" gorm:"not null;default:0"`
	CreatedTime        int64  `json:"created_time" gorm:"bigint;not null"`
	UpdatedTime        int64  `json:"updated_time" gorm:"bigint;not null"`
}

func (UsageStatistics) TableName() string {
	return "usage_statistics"
}

// GetUsageStatistics 获取用量统计数据
func GetUsageStatistics(startDate, endDate string, tokenId int, modelName string, page, pageSize int) ([]*UsageStatistics, int64, error) {
	var statistics []*UsageStatistics
	var total int64

	query := DB.Model(&UsageStatistics{})

	// 添加查询条件
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if tokenId > 0 {
		query = query.Where("token_id = ?", tokenId)
	}
	if modelName != "" {
		query = query.Where("model_name LIKE ?", "%"+modelName+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("date DESC, token_id ASC, model_name ASC").
		Offset(offset).Limit(pageSize).Find(&statistics).Error

	return statistics, total, err
}

// GetMonthlyUsageStatistics 获取月度用量统计数据
func GetMonthlyUsageStatistics(startDate, endDate string, tokenId int, modelName string, page, pageSize int) ([]*UsageStatistics, int64, error) {
	var statistics []*UsageStatistics
	var total int64

	// 使用原生SQL查询实现按月分组统计
	db := DB.Model(&UsageStatistics{})

	// 构建查询条件
	conditions := ""
	params := []interface{}{}

	// 添加日期范围条件（按月查询）
	if startDate != "" {
		conditions += " AND date >= ?"
		params = append(params, startDate+"-01")
	}
	if endDate != "" {
		// 获取 endDate 所在月份的最后一天
		if len(endDate) >= 7 {
			year := endDate[0:4]
			month := endDate[5:7]
			conditions += " AND date <= ?"
			params = append(params, year+"-"+month+"-31")
		}
	}
	
	if tokenId > 0 {
		conditions += " AND token_id = ?"
		params = append(params, tokenId)
	}
	if modelName != "" {
		conditions += " AND model_name LIKE ?"
		params = append(params, "%"+modelName+"%")
	}

	// 构建完整的SQL查询
	sql := `
		SELECT 
			MAX(id) as id,
			SUBSTR(date, 1, 7) as date,
			token_id,
			token_name,
			model_name,
			SUM(total_requests) as total_requests,
			SUM(successful_requests) as successful_requests,
			SUM(failed_requests) as failed_requests,
			SUM(total_tokens) as total_tokens,
			SUM(prompt_tokens) as prompt_tokens,
			SUM(completion_tokens) as completion_tokens,
			SUM(total_quota) as total_quota,
			MAX(created_time) as created_time,
			MAX(updated_time) as updated_time
		FROM usage_statistics 
		WHERE 1=1` + conditions + `
		GROUP BY SUBSTR(date, 1, 7), token_id, token_name, model_name
		ORDER BY date DESC, token_id ASC, model_name ASC
	`

	// 获取总数
	countSQL := `
		SELECT COUNT(*) as count FROM (
			SELECT 1
			FROM usage_statistics 
			WHERE 1=1` + conditions + `
			GROUP BY SUBSTR(date, 1, 7), token_id, token_name, model_name
		) as grouped_data
	`

	var countResult struct {
		Count int64 `json:"count"`
	}
	err := db.Raw(countSQL, params...).Scan(&countResult).Error
	if err != nil {
		return nil, 0, err
	}
	total = countResult.Count

	// 分页查询
	offset := (page - 1) * pageSize
	limitSQL := sql + " LIMIT ? OFFSET ?"
	params = append(params, pageSize, offset)

	err = db.Raw(limitSQL, params...).Scan(&statistics).Error
	return statistics, total, err
}

// UpsertUsageStatistics 插入或更新用量统计数据
func UpsertUsageStatistics(date string, tokenId int, tokenName, modelName string,
	totalRequests, successfulRequests, failedRequests int,
	totalTokens, promptTokens, completionTokens, totalQuota int) error {

	now := common.GetTimestamp()

	// 使用 ON DUPLICATE KEY UPDATE (MySQL) 或 ON CONFLICT (PostgreSQL)
	if common.UsingMySQL {
		return DB.Exec(`
			INSERT INTO usage_statistics 
			(date, token_id, token_name, model_name, total_requests, successful_requests, failed_requests,
			 total_tokens, prompt_tokens, completion_tokens, total_quota, created_time, updated_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			token_name = VALUES(token_name),
			total_requests = total_requests + VALUES(total_requests),
			successful_requests = successful_requests + VALUES(successful_requests),
			failed_requests = failed_requests + VALUES(failed_requests),
			total_tokens = total_tokens + VALUES(total_tokens),
			prompt_tokens = prompt_tokens + VALUES(prompt_tokens),
			completion_tokens = completion_tokens + VALUES(completion_tokens),
			total_quota = total_quota + VALUES(total_quota),
			updated_time = VALUES(updated_time)
		`, date, tokenId, tokenName, modelName, totalRequests, successfulRequests, failedRequests,
			totalTokens, promptTokens, completionTokens, totalQuota, now, now).Error
	} else if common.UsingPostgreSQL {
		return DB.Exec(`
			INSERT INTO usage_statistics 
			(date, token_id, token_name, model_name, total_requests, successful_requests, failed_requests,
			 total_tokens, prompt_tokens, completion_tokens, total_quota, created_time, updated_time)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			ON CONFLICT (date, token_id, model_name) DO UPDATE SET
			token_name = EXCLUDED.token_name,
			total_requests = usage_statistics.total_requests + EXCLUDED.total_requests,
			successful_requests = usage_statistics.successful_requests + EXCLUDED.successful_requests,
			failed_requests = usage_statistics.failed_requests + EXCLUDED.failed_requests,
			total_tokens = usage_statistics.total_tokens + EXCLUDED.total_tokens,
			prompt_tokens = usage_statistics.prompt_tokens + EXCLUDED.prompt_tokens,
			completion_tokens = usage_statistics.completion_tokens + EXCLUDED.completion_tokens,
			total_quota = usage_statistics.total_quota + EXCLUDED.total_quota,
			updated_time = EXCLUDED.updated_time
		`, date, tokenId, tokenName, modelName, totalRequests, successfulRequests, failedRequests,
			totalTokens, promptTokens, completionTokens, totalQuota, now, now).Error
	} else {
		// SQLite - 使用 INSERT OR REPLACE
		var existing UsageStatistics
		err := DB.Where("date = ? AND token_id = ? AND model_name = ?", date, tokenId, modelName).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// 记录不存在，直接插入
			newRecord := UsageStatistics{
				Date:               date,
				TokenId:            tokenId,
				TokenName:          tokenName,
				ModelName:          modelName,
				TotalRequests:      totalRequests,
				SuccessfulRequests: successfulRequests,
				FailedRequests:     failedRequests,
				TotalTokens:        totalTokens,
				PromptTokens:       promptTokens,
				CompletionTokens:   completionTokens,
				TotalQuota:         totalQuota,
				CreatedTime:        now,
				UpdatedTime:        now,
			}
			return DB.Create(&newRecord).Error
		} else if err != nil {
			return err
		} else {
			// 记录存在，更新数据
			updates := map[string]interface{}{
				"token_name":          tokenName,
				"total_requests":      existing.TotalRequests + totalRequests,
				"successful_requests": existing.SuccessfulRequests + successfulRequests,
				"failed_requests":     existing.FailedRequests + failedRequests,
				"total_tokens":        existing.TotalTokens + totalTokens,
				"prompt_tokens":       existing.PromptTokens + promptTokens,
				"completion_tokens":   existing.CompletionTokens + completionTokens,
				"total_quota":         existing.TotalQuota + totalQuota,
				"updated_time":        now,
			}
			return DB.Model(&existing).Updates(updates).Error
		}
	}
}

// RecordUsageStatistics 记录用量统计（从日志记录中调用）
func RecordUsageStatistics(tokenId int, tokenName, modelName string,
	promptTokens, completionTokens int, quota int, isSuccess bool) error {

	if tokenId <= 0 || modelName == "" {
		return errors.New("invalid parameters for usage statistics")
	}

	date := time.Now().Format("2006-01-02")
	totalTokens := promptTokens + completionTokens
	totalRequests := 1
	successfulRequests := 0
	failedRequests := 0

	if isSuccess {
		successfulRequests = 1
	} else {
		failedRequests = 1
	}

	return UpsertUsageStatistics(date, tokenId, tokenName, modelName,
		totalRequests, successfulRequests, failedRequests,
		totalTokens, promptTokens, completionTokens, quota)
}

// GetUsageStatisticsSummary 获取用量统计摘要信息
func GetUsageStatisticsSummary(startDate, endDate string, tokenId int, modelName string) (map[string]interface{}, error) {
	query := DB.Model(&UsageStatistics{})

	// 添加查询条件
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}
	if tokenId > 0 {
		query = query.Where("token_id = ?", tokenId)
	}
	if modelName != "" {
		query = query.Where("model_name LIKE ?", "%"+modelName+"%")
	}

	var result struct {
		TotalRequests      int `json:"total_requests"`
		SuccessfulRequests int `json:"successful_requests"`
		FailedRequests     int `json:"failed_requests"`
		TotalTokens        int `json:"total_tokens"`
		TotalQuota         int `json:"total_quota"`
	}

	err := query.Select(
		"SUM(total_requests) as total_requests",
		"SUM(successful_requests) as successful_requests",
		"SUM(failed_requests) as failed_requests",
		"SUM(total_tokens) as total_tokens",
		"SUM(total_quota) as total_quota",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_requests":      result.TotalRequests,
		"successful_requests": result.SuccessfulRequests,
		"failed_requests":     result.FailedRequests,
		"success_rate":        0.0,
		"total_tokens":        result.TotalTokens,
		"total_quota":         result.TotalQuota,
	}

	if result.TotalRequests > 0 {
		summary["success_rate"] = float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100
	}

	return summary, nil
}

// GetMonthlyUsageStatisticsSummary 获取月度用量统计摘要信息
func GetMonthlyUsageStatisticsSummary(startDate, endDate string, tokenId int, modelName string) (map[string]interface{}, error) {
	// 使用原生SQL查询实现按月分组统计
	db := DB.Model(&UsageStatistics{})

	// 构建查询条件
	conditions := ""
	params := []interface{}{}

	// 添加日期范围条件（按月查询）
	if startDate != "" {
		conditions += " AND date >= ?"
		params = append(params, startDate+"-01")
	}
	if endDate != "" {
		// 获取 endDate 所在月份的最后一天
		if len(endDate) >= 7 {
			year := endDate[0:4]
			month := endDate[5:7]
			conditions += " AND date <= ?"
			params = append(params, year+"-"+month+"-31")
		}
	}
	
	if tokenId > 0 {
		conditions += " AND token_id = ?"
		params = append(params, tokenId)
	}
	if modelName != "" {
		conditions += " AND model_name LIKE ?"
		params = append(params, "%"+modelName+"%")
	}

	// 构建完整的SQL查询
	sql := `
		SELECT 
			SUM(total_requests) as total_requests,
			SUM(successful_requests) as successful_requests,
			SUM(failed_requests) as failed_requests,
			SUM(total_tokens) as total_tokens,
			SUM(total_quota) as total_quota
		FROM (
			SELECT 
				SUM(total_requests) as total_requests,
				SUM(successful_requests) as successful_requests,
				SUM(failed_requests) as failed_requests,
				SUM(total_tokens) as total_tokens,
				SUM(total_quota) as total_quota
			FROM usage_statistics 
			WHERE 1=1` + conditions + `
			GROUP BY SUBSTR(date, 1, 7), token_id, token_name, model_name
		) as grouped_data
	`

	var result struct {
		TotalRequests      int `json:"total_requests"`
		SuccessfulRequests int `json:"successful_requests"`
		FailedRequests     int `json:"failed_requests"`
		TotalTokens        int `json:"total_tokens"`
		TotalQuota         int `json:"total_quota"`
	}

	err := db.Raw(sql, params...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_requests":      result.TotalRequests,
		"successful_requests": result.SuccessfulRequests,
		"failed_requests":     result.FailedRequests,
		"success_rate":        0.0,
		"total_tokens":        result.TotalTokens,
		"total_quota":         result.TotalQuota,
	}

	if result.TotalRequests > 0 {
		summary["success_rate"] = float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100
	}

	return summary, nil
}

// GetUserUsageStatistics 获取特定用户的用量统计数据
func GetUserUsageStatistics(userId int, startDate, endDate string, tokenId int, modelName string, page, pageSize int) ([]*UsageStatistics, int64, error) {
	var statistics []*UsageStatistics
	var total int64

	// 基本查询，需要通过token表连接过滤用户
	query := DB.Table("usage_statistics").
		Joins("JOIN tokens ON usage_statistics.token_id = tokens.id").
		Where("tokens.user_id = ?", userId)

	// 添加查询条件
	if startDate != "" {
		query = query.Where("usage_statistics.date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("usage_statistics.date <= ?", endDate)
	}
	if tokenId > 0 {
		query = query.Where("usage_statistics.token_id = ?", tokenId)
	}
	if modelName != "" {
		query = query.Where("usage_statistics.model_name LIKE ?", "%"+modelName+"%")
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Select("usage_statistics.*").
		Order("usage_statistics.date DESC, usage_statistics.token_id ASC, usage_statistics.model_name ASC").
		Offset(offset).Limit(pageSize).Find(&statistics).Error

	return statistics, total, err
}

// GetUserMonthlyUsageStatistics 获取特定用户的月度用量统计数据
func GetUserMonthlyUsageStatistics(userId int, startDate, endDate string, tokenId int, modelName string, page, pageSize int) ([]*UsageStatistics, int64, error) {
	var statistics []*UsageStatistics
	var total int64

	// 使用原生SQL查询实现按月分组统计
	db := DB.Table("usage_statistics").
		Joins("JOIN tokens ON usage_statistics.token_id = tokens.id").
		Where("tokens.user_id = ?", userId)

	// 构建查询条件
	conditions := " AND tokens.user_id = ?"
	params := []interface{}{userId}

	// 添加日期范围条件（按月查询）
	if startDate != "" {
		conditions += " AND usage_statistics.date >= ?"
		params = append(params, startDate+"-01")
	}
	if endDate != "" {
		// 获取 endDate 所在月份的最后一天
		if len(endDate) >= 7 {
			year := endDate[0:4]
			month := endDate[5:7]
			conditions += " AND usage_statistics.date <= ?"
			params = append(params, year+"-"+month+"-31")
		}
	}
	
	if tokenId > 0 {
		conditions += " AND usage_statistics.token_id = ?"
		params = append(params, tokenId)
	}
	if modelName != "" {
		conditions += " AND usage_statistics.model_name LIKE ?"
		params = append(params, "%"+modelName+"%")
	}

	// 构建完整的SQL查询
	sql := `
		SELECT 
			MAX(usage_statistics.id) as id,
			SUBSTR(usage_statistics.date, 1, 7) as date,
			usage_statistics.token_id,
			usage_statistics.token_name,
			usage_statistics.model_name,
			SUM(usage_statistics.total_requests) as total_requests,
			SUM(usage_statistics.successful_requests) as successful_requests,
			SUM(usage_statistics.failed_requests) as failed_requests,
			SUM(usage_statistics.total_tokens) as total_tokens,
			SUM(usage_statistics.prompt_tokens) as prompt_tokens,
			SUM(usage_statistics.completion_tokens) as completion_tokens,
			SUM(usage_statistics.total_quota) as total_quota,
			MAX(usage_statistics.created_time) as created_time,
			MAX(usage_statistics.updated_time) as updated_time
		FROM usage_statistics 
		JOIN tokens ON usage_statistics.token_id = tokens.id
		WHERE 1=1` + conditions + `
		GROUP BY SUBSTR(usage_statistics.date, 1, 7), usage_statistics.token_id, usage_statistics.token_name, usage_statistics.model_name
		ORDER BY date DESC, token_id ASC, model_name ASC
	`

	// 获取总数
	countSQL := `
		SELECT COUNT(*) as count FROM (
			SELECT 1
			FROM usage_statistics 
			JOIN tokens ON usage_statistics.token_id = tokens.id
			WHERE 1=1` + conditions + `
			GROUP BY SUBSTR(usage_statistics.date, 1, 7), usage_statistics.token_id, usage_statistics.token_name, usage_statistics.model_name
		) as grouped_data
	`

	var countResult struct {
		Count int64 `json:"count"`
	}
	err := db.Raw(countSQL, params...).Scan(&countResult).Error
	if err != nil {
		return nil, 0, err
	}
	total = countResult.Count

	// 分页查询
	offset := (page - 1) * pageSize
	limitSQL := sql + " LIMIT ? OFFSET ?"
	params = append(params, pageSize, offset)

	err = db.Raw(limitSQL, params...).Scan(&statistics).Error
	return statistics, total, err
}

// GetUserUsageStatisticsSummary 获取特定用户的用量统计摘要信息
func GetUserUsageStatisticsSummary(userId int, startDate, endDate string, tokenId int, modelName string) (map[string]interface{}, error) {
	query := DB.Table("usage_statistics").
		Joins("JOIN tokens ON usage_statistics.token_id = tokens.id").
		Where("tokens.user_id = ?", userId)

	// 添加查询条件
	if startDate != "" {
		query = query.Where("usage_statistics.date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("usage_statistics.date <= ?", endDate)
	}
	if tokenId > 0 {
		query = query.Where("usage_statistics.token_id = ?", tokenId)
	}
	if modelName != "" {
		query = query.Where("usage_statistics.model_name LIKE ?", "%"+modelName+"%")
	}

	var result struct {
		TotalRequests      int `json:"total_requests"`
		SuccessfulRequests int `json:"successful_requests"`
		FailedRequests     int `json:"failed_requests"`
		TotalTokens        int `json:"total_tokens"`
		TotalQuota         int `json:"total_quota"`
	}

	err := query.Select(
		"SUM(usage_statistics.total_requests) as total_requests",
		"SUM(usage_statistics.successful_requests) as successful_requests",
		"SUM(usage_statistics.failed_requests) as failed_requests",
		"SUM(usage_statistics.total_tokens) as total_tokens",
		"SUM(usage_statistics.total_quota) as total_quota",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_requests":      result.TotalRequests,
		"successful_requests": result.SuccessfulRequests,
		"failed_requests":     result.FailedRequests,
		"success_rate":        0.0,
		"total_tokens":        result.TotalTokens,
		"total_quota":         result.TotalQuota,
	}

	if result.TotalRequests > 0 {
		summary["success_rate"] = float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100
	}

	return summary, nil
}

// GetUserMonthlyUsageStatisticsSummary 获取特定用户的月度用量统计摘要信息
func GetUserMonthlyUsageStatisticsSummary(userId int, startDate, endDate string, tokenId int, modelName string) (map[string]interface{}, error) {
	// 使用原生SQL查询实现按月分组统计
	db := DB.Table("usage_statistics").
		Joins("JOIN tokens ON usage_statistics.token_id = tokens.id").
		Where("tokens.user_id = ?", userId)

	// 构建查询条件
	conditions := " AND tokens.user_id = ?"
	params := []interface{}{userId}

	// 添加日期范围条件（按月查询）
	if startDate != "" {
		conditions += " AND usage_statistics.date >= ?"
		params = append(params, startDate+"-01")
	}
	if endDate != "" {
		// 获取 endDate 所在月份的最后一天
		if len(endDate) >= 7 {
			year := endDate[0:4]
			month := endDate[5:7]
			conditions += " AND usage_statistics.date <= ?"
			params = append(params, year+"-"+month+"-31")
		}
	}
	
	if tokenId > 0 {
		conditions += " AND usage_statistics.token_id = ?"
		params = append(params, tokenId)
	}
	if modelName != "" {
		conditions += " AND usage_statistics.model_name LIKE ?"
		params = append(params, "%"+modelName+"%")
	}

	// 构建完整的SQL查询
	sql := `
		SELECT 
			SUM(total_requests) as total_requests,
			SUM(successful_requests) as successful_requests,
			SUM(failed_requests) as failed_requests,
			SUM(total_tokens) as total_tokens,
			SUM(total_quota) as total_quota
		FROM (
			SELECT 
				SUM(usage_statistics.total_requests) as total_requests,
				SUM(usage_statistics.successful_requests) as successful_requests,
				SUM(usage_statistics.failed_requests) as failed_requests,
				SUM(usage_statistics.total_tokens) as total_tokens,
				SUM(usage_statistics.total_quota) as total_quota
			FROM usage_statistics 
			JOIN tokens ON usage_statistics.token_id = tokens.id
			WHERE 1=1` + conditions + `
			GROUP BY SUBSTR(usage_statistics.date, 1, 7), usage_statistics.token_id, usage_statistics.token_name, usage_statistics.model_name
		) as grouped_data
	`

	var result struct {
		TotalRequests      int `json:"total_requests"`
		SuccessfulRequests int `json:"successful_requests"`
		FailedRequests     int `json:"failed_requests"`
		TotalTokens        int `json:"total_tokens"`
		TotalQuota         int `json:"total_quota"`
	}

	err := db.Raw(sql, params...).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_requests":      result.TotalRequests,
		"successful_requests": result.SuccessfulRequests,
		"failed_requests":     result.FailedRequests,
		"success_rate":        0.0,
		"total_tokens":        result.TotalTokens,
		"total_quota":         result.TotalQuota,
	}

	if result.TotalRequests > 0 {
		summary["success_rate"] = float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100
	}

	return summary, nil
}
