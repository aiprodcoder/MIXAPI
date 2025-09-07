package model

import (
	"errors"
	"fmt"
	"log"
	"one-api/common"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"gorm.io/gorm"
)

type Token struct {
	Id                 int            `json:"id"`
	UserId             int            `json:"user_id" gorm:"index"`
	Key                string         `json:"key" gorm:"type:char(48);uniqueIndex"`
	Status             int            `json:"status" gorm:"default:1"`
	Name               string         `json:"name" gorm:"index" `
	CreatedTime        int64          `json:"created_time" gorm:"bigint"`
	AccessedTime       int64          `json:"accessed_time" gorm:"bigint"`
	ExpiredTime        int64          `json:"expired_time" gorm:"bigint;default:-1"` // -1 means never expired
	RemainQuota        int            `json:"remain_quota" gorm:"default:0"`
	UnlimitedQuota     bool           `json:"unlimited_quota"`
	ModelLimitsEnabled bool           `json:"model_limits_enabled"`
	ModelLimits        string         `json:"model_limits" gorm:"type:varchar(1024);default:''"`
	AllowIps           *string        `json:"allow_ips" gorm:"default:''"`
	UsedQuota          int            `json:"used_quota" gorm:"default:0"` // used quota
	Group              string         `json:"group" gorm:"default:''"`
	DailyUsageCount    int            `json:"daily_usage_count" gorm:"default:0"`     // 今日使用次数
	TotalUsageCount    int            `json:"total_usage_count" gorm:"default:0"`     // 总使用次数
	LastUsageDate      string         `json:"last_usage_date" gorm:"default:''"`      // 最后使用日期(YYYY-MM-DD)
	RateLimitPerMinute int            `json:"rate_limit_per_minute" gorm:"default:0"` // 每分钟访问次数限制，0表示不限制
	RateLimitPerDay    int            `json:"rate_limit_per_day" gorm:"default:0"`    // 每日访问次数限制，0表示不限制
	LastRateLimitReset int64          `json:"last_rate_limit_reset" gorm:"default:0"` // 最后重置时间戳
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (token *Token) Clean() {
	token.Key = ""
}

func (token *Token) GetIpLimitsMap() map[string]any {
	// delete empty spaces
	//split with \n
	ipLimitsMap := make(map[string]any)
	if token.AllowIps == nil {
		return ipLimitsMap
	}
	cleanIps := strings.ReplaceAll(*token.AllowIps, " ", "")
	if cleanIps == "" {
		return ipLimitsMap
	}
	ips := strings.Split(cleanIps, "\n")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		ip = strings.ReplaceAll(ip, ",", "")
		if common.IsIP(ip) {
			ipLimitsMap[ip] = true
		}
	}
	return ipLimitsMap
}

func GetAllUserTokens(userId int, startIdx int, num int) ([]*Token, error) {
	var tokens []*Token
	var err error
	err = DB.Where("user_id = ?", userId).Order("id desc").Limit(num).Offset(startIdx).Find(&tokens).Error
	return tokens, err
}

func SearchUserTokens(userId int, keyword string, token string) (tokens []*Token, err error) {
	if token != "" {
		token = strings.Trim(token, "sk-")
	}
	err = DB.Where("user_id = ?", userId).Where("name LIKE ?", "%"+keyword+"%").Where(commonKeyCol+" LIKE ?", "%"+token+"%").Find(&tokens).Error
	return tokens, err
}

func ValidateUserToken(key string) (token *Token, err error) {
	log.Println("===========", key)
	if key == "" {
		return nil, errors.New("未提供令牌")
	}
	token, err = GetTokenByKey(key, false)
	if err == nil {
		if token.Status == common.TokenStatusExhausted {
			keyPrefix := key[:3]
			keySuffix := key[len(key)-3:]
			return token, errors.New("该令牌额度已用尽 TokenStatusExhausted[sk-" + keyPrefix + "***" + keySuffix + "]")
		} else if token.Status == common.TokenStatusExpired {
			return token, errors.New("该令牌已过期")
		}
		if token.Status != common.TokenStatusEnabled {
			return token, errors.New("该令牌状态不可用")
		}
		if token.ExpiredTime != -1 && token.ExpiredTime < common.GetTimestamp() {
			if !common.RedisEnabled {
				token.Status = common.TokenStatusExpired
				err := token.SelectUpdate()
				if err != nil {
					common.SysError("failed to update token status" + err.Error())
				}
			}
			return token, errors.New("该令牌已过期")
		}
		if !token.UnlimitedQuota && token.RemainQuota <= 0 {
			if !common.RedisEnabled {
				// in this case, we can make sure the token is exhausted
				token.Status = common.TokenStatusExhausted
				err := token.SelectUpdate()
				if err != nil {
					common.SysError("failed to update token status" + err.Error())
				}
			}
			keyPrefix := key[:3]
			keySuffix := key[len(key)-3:]
			return token, errors.New(fmt.Sprintf("[sk-%s***%s] 该令牌额度已用尽 !token.UnlimitedQuota && token.RemainQuota = %d", keyPrefix, keySuffix, token.RemainQuota))
		}

		// 检查访问频率限制
		if err := CheckRateLimit(token); err != nil {
			return token, err
		}

		return token, nil
	}
	return nil, errors.New("无效的令牌")
}

func GetTokenByIds(id int, userId int) (*Token, error) {
	if id == 0 || userId == 0 {
		return nil, errors.New("id 或 userId 为空！")
	}
	token := Token{Id: id, UserId: userId}
	var err error = nil
	err = DB.First(&token, "id = ? and user_id = ?", id, userId).Error
	return &token, err
}

func GetTokenById(id int) (*Token, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	token := Token{Id: id}
	var err error = nil
	err = DB.First(&token, "id = ?", id).Error
	if shouldUpdateRedis(true, err) {
		gopool.Go(func() {
			if err := cacheSetToken(token); err != nil {
				common.SysError("failed to update user status cache: " + err.Error())
			}
		})
	}
	return &token, err
}

func GetTokenByKey(key string, fromDB bool) (token *Token, err error) {
	defer func() {
		// Update Redis cache asynchronously on successful DB read
		if shouldUpdateRedis(fromDB, err) && token != nil {
			gopool.Go(func() {
				if err := cacheSetToken(*token); err != nil {
					common.SysError("failed to update user status cache: " + err.Error())
				}
			})
		}
	}()
	if !fromDB && common.RedisEnabled {
		// Try Redis first
		token, err := cacheGetTokenByKey(key)
		if err == nil {
			return token, nil
		}
		// Don't return error - fall through to DB
	}
	fromDB = true
	err = DB.Where(commonKeyCol+" = ?", key).First(&token).Error
	return token, err
}

func (token *Token) Insert() error {
	var err error
	err = DB.Create(token).Error
	return err
}

// Update Make sure your token's fields is completed, because this will update non-zero values
func (token *Token) Update() (err error) {
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheSetToken(*token)
				if err != nil {
					common.SysError("failed to update token cache: " + err.Error())
				}
			})
		}
	}()
	err = DB.Model(token).Select("name", "status", "expired_time", "remain_quota", "unlimited_quota",
		"model_limits_enabled", "model_limits", "allow_ips", "group", "daily_usage_count", "total_usage_count", "last_usage_date",
		"rate_limit_per_minute", "rate_limit_per_day", "last_rate_limit_reset").Updates(token).Error
	return err
}

func (token *Token) SelectUpdate() (err error) {
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheSetToken(*token)
				if err != nil {
					common.SysError("failed to update token cache: " + err.Error())
				}
			})
		}
	}()
	// This can update zero values
	return DB.Model(token).Select("accessed_time", "status").Updates(token).Error
}

func (token *Token) Delete() (err error) {
	defer func() {
		if shouldUpdateRedis(true, err) {
			gopool.Go(func() {
				err := cacheDeleteToken(token.Key)
				if err != nil {
					common.SysError("failed to delete token cache: " + err.Error())
				}
			})
		}
	}()
	err = DB.Delete(token).Error
	return err
}

func (token *Token) IsModelLimitsEnabled() bool {
	return token.ModelLimitsEnabled
}

func (token *Token) GetModelLimits() []string {
	if token.ModelLimits == "" {
		return []string{}
	}
	return strings.Split(token.ModelLimits, ",")
}

func (token *Token) GetModelLimitsMap() map[string]bool {
	limits := token.GetModelLimits()
	limitsMap := make(map[string]bool)
	for _, limit := range limits {
		limitsMap[limit] = true
	}
	return limitsMap
}

func DisableModelLimits(tokenId int) error {
	token, err := GetTokenById(tokenId)
	if err != nil {
		return err
	}
	token.ModelLimitsEnabled = false
	token.ModelLimits = ""
	return token.Update()
}

func DeleteTokenById(id int, userId int) (err error) {
	// Why we need userId here? In case user want to delete other's token.
	if id == 0 || userId == 0 {
		return errors.New("id 或 userId 为空！")
	}
	token := Token{Id: id, UserId: userId}
	err = DB.Where(token).First(&token).Error
	if err != nil {
		return err
	}
	return token.Delete()
}

func IncreaseTokenQuota(id int, key string, quota int) (err error) {
	if quota < 0 {
		return errors.New("quota 不能为负数！")
	}
	if common.RedisEnabled {
		gopool.Go(func() {
			err := cacheIncrTokenQuota(key, int64(quota))
			if err != nil {
				common.SysError("failed to increase token quota: " + err.Error())
			}
		})
	}
	if common.BatchUpdateEnabled {
		addNewRecord(BatchUpdateTypeTokenQuota, id, quota)
		return nil
	}
	return increaseTokenQuota(id, quota)
}

func increaseTokenQuota(id int, quota int) (err error) {
	err = DB.Model(&Token{}).Where("id = ?", id).Updates(
		map[string]interface{}{
			"remain_quota":  gorm.Expr("remain_quota + ?", quota),
			"used_quota":    gorm.Expr("used_quota - ?", quota),
			"accessed_time": common.GetTimestamp(),
		},
	).Error
	return err
}

func DecreaseTokenQuota(id int, key string, quota int) (err error) {
	if quota < 0 {
		return errors.New("quota 不能为负数！")
	}
	if common.RedisEnabled {
		gopool.Go(func() {
			err := cacheDecrTokenQuota(key, int64(quota))
			if err != nil {
				common.SysError("failed to decrease token quota: " + err.Error())
			}
		})
	}
	if common.BatchUpdateEnabled {
		addNewRecord(BatchUpdateTypeTokenQuota, id, -quota)
		return nil
	}
	return decreaseTokenQuota(id, quota)
}

func decreaseTokenQuota(id int, quota int) (err error) {
	err = DB.Model(&Token{}).Where("id = ?", id).Updates(
		map[string]interface{}{
			"remain_quota":  gorm.Expr("remain_quota - ?", quota),
			"used_quota":    gorm.Expr("used_quota + ?", quota),
			"accessed_time": common.GetTimestamp(),
		},
	).Error
	return err
}

// CountUserTokens returns total number of tokens for the given user, used for pagination
func CountUserTokens(userId int) (int64, error) {
	var total int64
	err := DB.Model(&Token{}).Where("user_id = ?", userId).Count(&total).Error
	return total, err
}

// BatchDeleteTokens 删除指定用户的一组令牌，返回成功删除数量
func BatchDeleteTokens(ids []int, userId int) (int, error) {
	if len(ids) == 0 {
		return 0, errors.New("ids 不能为空！")
	}

	tx := DB.Begin()

	var tokens []Token
	if err := tx.Where("user_id = ? AND id IN (?)", userId, ids).Find(&tokens).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Where("user_id = ? AND id IN (?)", userId, ids).Delete(&Token{}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	if common.RedisEnabled {
		gopool.Go(func() {
			for _, t := range tokens {
				_ = cacheDeleteToken(t.Key)
			}
		})
	}

	return len(tokens), nil
}

// IncreaseTokenUsageCount 增加Token的使用次数
func IncreaseTokenUsageCount(key string) error {
	if key == "" {
		return errors.New("key 不能为空")
	}

	// 获取当前日期
	currentDate := common.GetTimeString()[:10] // YYYY-MM-DD格式

	// 更新数据库
	err := DB.Model(&Token{}).Where(commonKeyCol+" = ?", key).Updates(map[string]interface{}{
		"total_usage_count": gorm.Expr("total_usage_count + 1"),
		"daily_usage_count": gorm.Expr("CASE WHEN last_usage_date = ? THEN daily_usage_count + 1 ELSE 1 END", currentDate),
		"last_usage_date":   currentDate,
		"accessed_time":     common.GetTimestamp(),
	}).Error

	// 更新缓存
	if common.RedisEnabled && err == nil {
		gopool.Go(func() {
			// 重新缓存Token信息
			token, getErr := GetTokenByKey(key, true) // 从DB获取最新数据
			if getErr == nil {
				_ = cacheSetToken(*token)
			}
		})
	}

	return err
}

// CheckRateLimit 检查令牌的访问频率限制
func CheckRateLimit(token *Token) error {
	if token == nil {
		return errors.New("token不能为空")
	}

	// 如果没有设置限制，直接返回
	if token.RateLimitPerMinute <= 0 && token.RateLimitPerDay <= 0 {
		return nil
	}

	currentTime := time.Now()
	currentTimestamp := currentTime.Unix()

	// 获取当前时间的分钟级时间戳（用于分钟级限制）
	currentMinute := currentTime.Truncate(time.Minute).Unix()

	// 获取当天开始的时间戳（用于日级限制）
	currentDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()

	// 检查是否需要重置计数器
	needUpdate := false
	originalToken := *token // 保存原始状态

	// 检查分钟级限制
	if token.RateLimitPerMinute > 0 {
		// 如果上次重置时间不在当前分钟内，重置分钟计数器
		if token.LastRateLimitReset < currentMinute {
			token.LastRateLimitReset = currentTimestamp
			needUpdate = true
		}

		// 计算当前分钟内的使用次数
		var minuteCount int64
		err := DB.Model(&TokenUsageLog{}).Where("token_id = ? AND created_at >= ?", token.Id, currentMinute).Count(&minuteCount).Error
		if err != nil {
			common.SysError("检查分钟级使用次数失败: " + err.Error())
			return errors.New("系统错误，请稍后再试")
		}

		if int(minuteCount) >= token.RateLimitPerMinute {
			return errors.New("超出分钟限制，请稍后再试")
		}
	}

	// 检查日级限制
	if token.RateLimitPerDay > 0 {
		// 计算当天内的使用次数
		var dayCount int64
		err := DB.Model(&TokenUsageLog{}).Where("token_id = ? AND created_at >= ?", token.Id, currentDay).Count(&dayCount).Error
		if err != nil {
			common.SysError("检查日级使用次数失败: " + err.Error())
			return errors.New("系统错误，请稍后再试")
		}

		if int(dayCount) >= token.RateLimitPerDay {
			return errors.New("超出日限制，请稍后再试")
		}
	}

	// 如果需要更新，更新数据库
	if needUpdate {
		err := DB.Model(token).Select("last_rate_limit_reset").Updates(token).Error
		if err != nil {
			common.SysError("更新令牌重置时间失败: " + err.Error())
			// 恢复原始状态
			*token = originalToken
		}
	}

	return nil
}

// TokenUsageLog Token使用日志表（用于频率限制）
type TokenUsageLog struct {
	Id        int   `json:"id" gorm:"primaryKey"`
	TokenId   int   `json:"token_id" gorm:"index:idx_token_created"`
	CreatedAt int64 `json:"created_at" gorm:"index:idx_token_created;index:idx_created"`
}

func (TokenUsageLog) TableName() string {
	return "token_usage_logs"
}

// RecordTokenUsage 记录令牌使用（用于频率限制）
func RecordTokenUsage(tokenId int) error {
	if tokenId <= 0 {
		return errors.New("tokenId不能为空")
	}

	usageLog := TokenUsageLog{
		TokenId:   tokenId,
		CreatedAt: time.Now().Unix(),
	}

	return DB.Create(&usageLog).Error
}
