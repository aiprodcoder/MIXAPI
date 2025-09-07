package controller

import (
	"net/http"
	"one-api/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUsageStatistics 获取用量统计数据
func GetUsageStatistics(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tokenId, _ := strconv.Atoi(c.Query("token_id"))
	modelName := c.Query("model_name")

	// 设置默认查询最近7天
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 查询数据
	statistics, total, err := model.GetUsageStatistics(startDate, endDate, tokenId, modelName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取摘要信息
	summary, err := model.GetUsageStatisticsSummary(startDate, endDate, tokenId, modelName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items":     statistics,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"summary":   summary,
		},
	})
}

// GetUsageStatisticsSummary 获取用量统计摘要
func GetUsageStatisticsSummary(c *gin.Context) {
	// 获取查询参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tokenId, _ := strconv.Atoi(c.Query("token_id"))
	modelName := c.Query("model_name")

	// 设置默认查询最近7天
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// 获取摘要信息
	summary, err := model.GetUsageStatisticsSummary(startDate, endDate, tokenId, modelName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    summary,
	})
}

// GetUserUsageStatistics 获取用户自己的用量统计数据
func GetUserUsageStatistics(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户ID不能为空",
		})
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tokenId, _ := strconv.Atoi(c.Query("token_id"))
	modelName := c.Query("model_name")

	// 设置默认查询最近7天
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 如果指定了tokenId，需要验证token是否属于当前用户
	if tokenId > 0 {
		token, err := model.GetTokenByIds(tokenId, userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "Token不存在或无权访问",
			})
			return
		}
		if token.UserId != userId {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无权访问该Token的统计数据",
			})
			return
		}
	}

	// 查询数据 - 这里需要修改model函数来支持用户过滤
	statistics, total, err := model.GetUserUsageStatistics(userId, startDate, endDate, tokenId, modelName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取摘要信息
	summary, err := model.GetUserUsageStatisticsSummary(userId, startDate, endDate, tokenId, modelName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"items":     statistics,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"summary":   summary,
		},
	})
}
