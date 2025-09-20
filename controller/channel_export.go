package controller

import (
	"bytes"
	"net/http"
	"one-api/common"
	"one-api/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ExportChannels(c *gin.Context) {
	// 获取所有渠道数据
	channels, err := model.GetAllChannels(0, 0, true, false)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			common.SysError("Error closing Excel file: " + err.Error())
		}
	}()

	// 创建工作表
	sheetName := "Channels"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{
		"ID", "名称", "类型", "状态", "密钥", "组织", "测试模型",
		"权重", "创建时间", "测试时间", "响应时间", "基础URL", "其他",
		"余额", "余额更新时间", "模型", "分组", "已用配额",
		"模型映射", "状态码映射", "优先级", "自动禁用", "标签", "额外设置", "参数覆盖",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for i, channel := range channels {
		row := i + 2 // 从第二行开始填充数据
		data := []interface{}{
			channel.Id,
			channel.Name,
			channel.Type,
			channel.Status,
			channel.Key,
			channel.OpenAIOrganization,
			channel.TestModel,
			channel.Weight,
			channel.CreatedTime,
			channel.TestTime,
			channel.ResponseTime,
			channel.BaseURL,
			channel.Other,
			channel.Balance,
			channel.BalanceUpdatedTime,
			channel.Models,
			channel.Group,
			channel.UsedQuota,
			channel.ModelMapping,
			channel.StatusCodeMapping,
			channel.Priority,
			channel.AutoBan,
			channel.Tag,
			channel.Setting,
			channel.ParamOverride,
		}

		for j, value := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=channels.xlsx")
	c.Header("Cache-Control", "no-cache")

	// 将Excel文件写入缓冲区
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		common.ApiError(c, err)
		return
	}

	// 返回Excel文件
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
