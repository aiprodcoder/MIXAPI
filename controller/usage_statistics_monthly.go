package controller

import (
	"net/http"
	"one-api/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetMonthlyUsageStatistics 获取月度用量统计数据
func GetMonthlyUsageStatistics(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tokenId, _ := strconv.Atoi(c.Query("token_id"))
	modelName := c.Query("model_name")

	// 设置默认查询最近6个月
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, -6, 0).Format("2006-01")
		endDate = now.Format("2006-01")
	}

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 查询数据
	statistics, total, err := model.GetMonthlyUsageStatistics(startDate, endDate, tokenId, modelName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取摘要信息
	summary, err := model.GetMonthlyUsageStatisticsSummary(startDate, endDate, tokenId, modelName)
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

// GetMonthlyUsageStatisticsSummary 获取月度用量统计摘要
func GetMonthlyUsageStatisticsSummary(c *gin.Context) {
	// 获取查询参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tokenId, _ := strconv.Atoi(c.Query("token_id"))
	modelName := c.Query("model_name")

	// 设置默认查询最近6个月
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, -6, 0).Format("2006-01")
		endDate = now.Format("2006-01")
	}

	// 获取摘要信息
	summary, err := model.GetMonthlyUsageStatisticsSummary(startDate, endDate, tokenId, modelName)
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

// GetUserMonthlyUsageStatistics 获取用户自己的月度用量统计数据
func GetUserMonthlyUsageStatistics(c *gin.Context) {
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

	// 设置默认查询最近6个月
	if startDate == "" && endDate == "" {
		now := time.Now()
		startDate = now.AddDate(0, -6, 0).Format("2006-01")
		endDate = now.Format("2006-01")
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

	// 查询数据
	statistics, total, err := model.GetUserMonthlyUsageStatistics(userId, startDate, endDate, tokenId, modelName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 获取摘要信息
	summary, err := model.GetUserMonthlyUsageStatisticsSummary(userId, startDate, endDate, tokenId, modelName)
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