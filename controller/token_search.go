package controller

import (
	"net/http"
	"one-api/common"
	"one-api/model"
	"one-api/setting/ratio_setting"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenSearchRequest struct {
	Token string `json:"token" binding:"required"`
}

type TokenInfoResponse struct {
	TokenName     string  `json:"token_name"`
	RemainQuota   int     `json:"remain_quota"`
	UsedQuota     int     `json:"used_quota"`
	UnlimitedQuota bool    `json:"unlimited_quota"`
	ExpiredTime   int64   `json:"expired_time"`
	Status        int     `json:"status"`
	ModelRatio    float64 `json:"model_ratio"`
}

func SearchTokenByToken(c *gin.Context) {
	var req TokenSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiError(c, err)
		return
	}

	// 去掉前缀"sk-"后再查询数据库
	tokenKey := strings.TrimPrefix(req.Token, "sk-")

	// 获取token信息
	token, err := model.GetTokenByKey(tokenKey, false)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// 获取配比关系，默认使用gpt-3.5-turbo的配比
	modelRatio, _, _ := ratio_setting.GetModelRatio("gpt-3.5-turbo")

	// 构造返回信息
	response := TokenInfoResponse{
		TokenName:      token.Name,
		RemainQuota:    token.RemainQuota,
		UsedQuota:      token.UsedQuota,
		UnlimitedQuota: token.UnlimitedQuota,
		ExpiredTime:    token.ExpiredTime,
		Status:         token.Status,
		ModelRatio:     modelRatio,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    response,
	})
}