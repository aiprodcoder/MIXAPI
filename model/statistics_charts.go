package model

import (
	"one-api/common"
	"strconv"
	"time"
)

// ChannelStatistics 按渠道统计的数据结构
type ChannelStatistics struct {
	Time       string `json:"time"`
	ChannelId  int    `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Quota      int    `json:"quota"`
	Count      int    `json:"count"`
}

// TokenStatistics 按令牌统计的数据结构
type TokenStatistics struct {
	Time      string `json:"time"`
	TokenName string `json:"token_name"`
	Quota     int    `json:"quota"`
	Count     int    `json:"count"`
}

// UserStatistics 按用户统计的数据结构
type UserStatistics struct {
	Time     string `json:"time"`
	Username string `json:"username"`
	Quota    int    `json:"quota"`
	Count    int    `json:"count"`
}

// GetChannelStatistics 获取按渠道统计的数据
func GetChannelStatistics(startTimestamp, endTimestamp int, username, tokenName, modelName string, channel int, group, defaultTime string) ([]ChannelStatistics, error) {
	var statistics []ChannelStatistics

	// 构建查询条件
	tx := LOG_DB.Table("logs").Select(`
		created_at,
		channel_id,
		COUNT(*) as count,
		SUM(quota) as quota
	`).Where("type = ?", 2) // LogTypeConsume = 2

	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if channel != 0 {
		tx = tx.Where("channel_id = ?", channel)
	}
	if group != "" {
		tx = tx.Where("`group` = ?", group)
	}
	if startTimestamp != 0 {
		tx = tx.Where("created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("created_at <= ?", endTimestamp)
	}

	// 按时间粒度分组 - 使用数据库无关的函数
	var groupBy string
	switch defaultTime {
	case "hour":
		if common.UsingMySQL {
			groupBy = "channel_id, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d %H:00:00')"
		} else if common.UsingPostgreSQL {
			groupBy = "channel_id, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD HH24:00:00')"
		} else {
			// SQLite
			groupBy = "channel_id, strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'unixepoch'))"
		}
	case "day":
		if common.UsingMySQL {
			groupBy = "channel_id, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "channel_id, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "channel_id, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	case "week":
		if common.UsingMySQL {
			groupBy = "channel_id, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%U')"
		} else if common.UsingPostgreSQL {
			groupBy = "channel_id, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-IW')"
		} else {
			// SQLite
			groupBy = "channel_id, strftime('%Y-%W', datetime(created_at, 'unixepoch'))"
		}
	default:
		if common.UsingMySQL {
			groupBy = "channel_id, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "channel_id, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "channel_id, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	}

	// 按渠道和时间分组
	tx = tx.Group(groupBy).Order("MIN(created_at) ASC")

	// 执行查询
	var results []struct {
		CreatedAt  int64 `json:"created_at"`
		ChannelId  int   `json:"channel_id"`
		Count      int   `json:"count"`
		Quota      int   `json:"quota"`
	}

	err := tx.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 获取渠道名称映射
	channelIds := make([]int, 0)
	channelMap := make(map[int]string)
	for _, result := range results {
		if result.ChannelId != 0 {
			channelIds = append(channelIds, result.ChannelId)
		}
	}

	if len(channelIds) > 0 {
		var channels []struct {
			Id   int    `gorm:"column:id"`
			Name string `gorm:"column:name"`
		}
		if err = DB.Table("channels").Select("id, name").Where("id IN ?", channelIds).Find(&channels).Error; err != nil {
			return nil, err
		}
		for _, channel := range channels {
			channelMap[channel.Id] = channel.Name
		}
	}

	// 格式化时间并构建返回结果
	for _, result := range results {
		timeStr := formatTime(result.CreatedAt, defaultTime)
		statistics = append(statistics, ChannelStatistics{
			Time:        timeStr,
			ChannelId:   result.ChannelId,
			ChannelName: channelMap[result.ChannelId],
			Quota:       result.Quota,
			Count:       result.Count,
		})
	}

	return statistics, nil
}

// GetTokenStatistics 获取按令牌统计的数据
func GetTokenStatistics(startTimestamp, endTimestamp int, username, tokenName, modelName string, channel int, group, defaultTime string) ([]TokenStatistics, error) {
	var statistics []TokenStatistics

	// 构建查询条件
	tx := LOG_DB.Table("logs").Select(`
		created_at,
		token_name,
		COUNT(*) as count,
		SUM(quota) as quota
	`).Where("type = ?", 2) // LogTypeConsume = 2

	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if channel != 0 {
		tx = tx.Where("channel_id = ?", channel)
	}
	if group != "" {
		tx = tx.Where("`group` = ?", group)
	}
	if startTimestamp != 0 {
		tx = tx.Where("created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("created_at <= ?", endTimestamp)
	}

	// 按时间粒度分组 - 使用数据库无关的函数
	var groupBy string
	switch defaultTime {
	case "hour":
		if common.UsingMySQL {
			groupBy = "token_name, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d %H:00:00')"
		} else if common.UsingPostgreSQL {
			groupBy = "token_name, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD HH24:00:00')"
		} else {
			// SQLite
			groupBy = "token_name, strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'unixepoch'))"
		}
	case "day":
		if common.UsingMySQL {
			groupBy = "token_name, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "token_name, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "token_name, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	case "week":
		if common.UsingMySQL {
			groupBy = "token_name, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%U')"
		} else if common.UsingPostgreSQL {
			groupBy = "token_name, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-IW')"
		} else {
			// SQLite
			groupBy = "token_name, strftime('%Y-%W', datetime(created_at, 'unixepoch'))"
		}
	default:
		if common.UsingMySQL {
			groupBy = "token_name, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "token_name, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "token_name, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	}

	// 按令牌和时间分组
	tx = tx.Group(groupBy).Order("MIN(created_at) ASC")

	// 执行查询
	var results []struct {
		CreatedAt int64  `json:"created_at"`
		TokenName string `json:"token_name"`
		Count     int    `json:"count"`
		Quota     int    `json:"quota"`
	}

	err := tx.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 格式化时间并构建返回结果
	for _, result := range results {
		timeStr := formatTime(result.CreatedAt, defaultTime)
		statistics = append(statistics, TokenStatistics{
			Time:      timeStr,
			TokenName: result.TokenName,
			Quota:     result.Quota,
			Count:     result.Count,
		})
	}

	return statistics, nil
}

// GetUserStatistics 获取按用户统计的数据
func GetUserStatistics(startTimestamp, endTimestamp int, username, tokenName, modelName string, channel int, group, defaultTime string) ([]UserStatistics, error) {
	var statistics []UserStatistics

	// 构建查询条件
	tx := LOG_DB.Table("logs").Select(`
		created_at,
		username,
		COUNT(*) as count,
		SUM(quota) as quota
	`).Where("type = ?", 2) // LogTypeConsume = 2

	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
	}
	if modelName != "" {
		tx = tx.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if channel != 0 {
		tx = tx.Where("channel_id = ?", channel)
	}
	if group != "" {
		tx = tx.Where("`group` = ?", group)
	}
	if startTimestamp != 0 {
		tx = tx.Where("created_at >= ?", startTimestamp)
	}
	if endTimestamp != 0 {
		tx = tx.Where("created_at <= ?", endTimestamp)
	}

	// 按时间粒度分组 - 使用数据库无关的函数
	var groupBy string
	switch defaultTime {
	case "hour":
		if common.UsingMySQL {
			groupBy = "username, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d %H:00:00')"
		} else if common.UsingPostgreSQL {
			groupBy = "username, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD HH24:00:00')"
		} else {
			// SQLite
			groupBy = "username, strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'unixepoch'))"
		}
	case "day":
		if common.UsingMySQL {
			groupBy = "username, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "username, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "username, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	case "week":
		if common.UsingMySQL {
			groupBy = "username, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%U')"
		} else if common.UsingPostgreSQL {
			groupBy = "username, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-IW')"
		} else {
			// SQLite
			groupBy = "username, strftime('%Y-%W', datetime(created_at, 'unixepoch'))"
		}
	default:
		if common.UsingMySQL {
			groupBy = "username, DATE_FORMAT(FROM_UNIXTIME(created_at), '%Y-%m-%d')"
		} else if common.UsingPostgreSQL {
			groupBy = "username, TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY-MM-DD')"
		} else {
			// SQLite
			groupBy = "username, strftime('%Y-%m-%d', datetime(created_at, 'unixepoch'))"
		}
	}

	// 按用户和时间分组
	tx = tx.Group(groupBy).Order("MIN(created_at) ASC")

	// 执行查询
	var results []struct {
		CreatedAt int64  `json:"created_at"`
		Username  string `json:"username"`
		Count     int    `json:"count"`
		Quota     int    `json:"quota"`
	}

	err := tx.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 格式化时间并构建返回结果
	for _, result := range results {
		timeStr := formatTime(result.CreatedAt, defaultTime)
		statistics = append(statistics, UserStatistics{
			Time:     timeStr,
			Username: result.Username,
			Quota:    result.Quota,
			Count:    result.Count,
		})
	}

	return statistics, nil
}

// formatTime 根据时间粒度格式化时间
func formatTime(timestamp int64, defaultTime string) string {
	t := time.Unix(timestamp, 0)
	switch defaultTime {
	case "hour":
		return t.Format("2006-01-02 15:00")
	case "day":
		return t.Format("2006-01-02")
	case "week":
		_, week := t.ISOWeek()
		return t.Format("2006-01") + "-W" + strconv.Itoa(week)
	default:
		return t.Format("2006-01-02")
	}
}