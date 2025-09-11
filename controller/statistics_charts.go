package controller

import (
	"net/http"
	"one-api/common"
	"one-api/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetChannelStatistics 获取按渠道统计的数据
func GetChannelStatistics(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	group := c.Query("group")
	defaultTime := c.Query("default_time")

	// 设置默认时间范围为最近7天
	if startTimestamp == 0 && endTimestamp == 0 {
		endTimestamp = time.Now().Unix()
		startTimestamp = endTimestamp - 7*24*3600
	}

	// 设置默认时间粒度
	if defaultTime == "" {
		defaultTime = "day"
	}

	statistics, err := model.GetChannelStatistics(int(startTimestamp), int(endTimestamp), username, tokenName, modelName, channel, group, defaultTime)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    statistics,
	})
}

// GetTokenStatistics 获取按令牌统计的数据
func GetTokenStatistics(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	group := c.Query("group")
	defaultTime := c.Query("default_time")

	// 设置默认时间范围为最近7天
	if startTimestamp == 0 && endTimestamp == 0 {
		endTimestamp = time.Now().Unix()
		startTimestamp = endTimestamp - 7*24*3600
	}

	// 设置默认时间粒度
	if defaultTime == "" {
		defaultTime = "day"
	}

	statistics, err := model.GetTokenStatistics(int(startTimestamp), int(endTimestamp), username, tokenName, modelName, channel, group, defaultTime)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    statistics,
	})
}

// GetUserStatistics 获取按用户统计的数据
func GetUserStatistics(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	group := c.Query("group")
	defaultTime := c.Query("default_time")

	// 设置默认时间范围为最近7天
	if startTimestamp == 0 && endTimestamp == 0 {
		endTimestamp = time.Now().Unix()
		startTimestamp = endTimestamp - 7*24*3600
	}

	// 设置默认时间粒度
	if defaultTime == "" {
		defaultTime = "day"
	}

	statistics, err := model.GetUserStatistics(int(startTimestamp), int(endTimestamp), username, tokenName, modelName, channel, group, defaultTime)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    statistics,
	})
}