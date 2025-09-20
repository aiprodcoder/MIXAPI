package controller

import (
	"fmt"
	"net/http"
	"one-api/common"
	"one-api/model"
	"strconv"

	// "strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ImportChannels(c *gin.Context) {
	// 从请求中获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		common.ApiError(c, fmt.Errorf("failed to get file: %w", err))
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		common.ApiError(c, fmt.Errorf("failed to open file: %w", err))
		return
	}
	defer src.Close()

	// 读取Excel文件
	f, err := excelize.OpenReader(src)
	if err != nil {
		common.ApiError(c, fmt.Errorf("failed to open Excel file: %w", err))
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			common.SysError("Error closing Excel file: " + err.Error())
		}
	}()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		common.ApiError(c, fmt.Errorf("no sheets found in Excel file"))
		return
	}

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		common.ApiError(c, fmt.Errorf("failed to read rows: %w", err))
		return
	}

	// 检查是否有数据
	if len(rows) <= 1 {
		common.ApiError(c, fmt.Errorf("no data found in Excel file"))
		return
	}

	// 解析表头
	headers := rows[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// 解析数据行
	channels := make([]model.Channel, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// 创建渠道对象
		channel := model.Channel{}

		// 根据表头映射填充数据
		for header, colIndex := range headerMap {
			if colIndex >= len(row) {
				continue
			}
			value := row[colIndex]

			switch header {
			case "ID":
				// ID由数据库自动生成，不需要设置
			case "名称":
				channel.Name = value
			case "类型":
				if v, err := strconv.Atoi(value); err == nil {
					channel.Type = v
				}
			case "状态":
				if v, err := strconv.Atoi(value); err == nil {
					channel.Status = v
				}
			case "密钥":
				channel.Key = value
			case "组织":
				if value != "" {
					channel.OpenAIOrganization = &value
				}
			case "测试模型":
				if value != "" {
					channel.TestModel = &value
				}
			case "权重":
				if value != "" {
					if v, err := strconv.ParseUint(value, 10, 32); err == nil {
						weight := uint(v)
						channel.Weight = &weight
					}
				}
			case "创建时间":
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					channel.CreatedTime = v
				}
			case "测试时间":
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					channel.TestTime = v
				}
			case "响应时间":
				if v, err := strconv.Atoi(value); err == nil {
					channel.ResponseTime = v
				}
			case "基础URL":
				if value != "" {
					channel.BaseURL = &value
				}
			case "其他":
				channel.Other = value
			case "余额":
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					channel.Balance = v
				}
			case "余额更新时间":
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					channel.BalanceUpdatedTime = v
				}
			case "模型":
				channel.Models = value
			case "分组":
				channel.Group = value
			case "已用配额":
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					channel.UsedQuota = v
				}
			case "模型映射":
				if value != "" {
					channel.ModelMapping = &value
				}
			case "状态码映射":
				if value != "" {
					channel.StatusCodeMapping = &value
				}
			case "优先级":
				if value != "" {
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						channel.Priority = &v
					}
				}
			case "自动禁用":
				if value != "" {
					if v, err := strconv.Atoi(value); err == nil {
						channel.AutoBan = &v
					}
				}
			case "标签":
				if value != "" {
					channel.Tag = &value
				}
			case "额外设置":
				if value != "" {
					channel.Setting = &value
				}
			case "参数覆盖":
				if value != "" {
					channel.ParamOverride = &value
				}
			}
		}

		// 设置默认值
		if channel.CreatedTime == 0 {
			channel.CreatedTime = common.GetTimestamp()
		}

		// 如果没有设置状态，默认为启用
		if channel.Status == 0 {
			channel.Status = 1
		}

		channels = append(channels, channel)
	}

	// 批量插入渠道
	err = model.BatchInsertChannels(channels)
	if err != nil {
		common.ApiError(c, fmt.Errorf("failed to insert channels: %w", err))
		return
	}

	// 初始化渠道缓存
	model.InitChannelCache()

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("成功导入 %d 个渠道", len(channels)),
	})
}
