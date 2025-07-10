package handlers

import (
	"net/http"
	"strconv"

	"WechatGo/services"

	"github.com/gin-gonic/gin"
)

// RankingHandler 排行榜处理器
type RankingHandler struct {
	dbService *services.DatabaseService
}

// NewRankingHandler 创建排行榜处理器
func NewRankingHandler(dbService *services.DatabaseService) *RankingHandler {
	return &RankingHandler{
		dbService: dbService,
	}
}

// GetRankingList 获取排行榜列表
func (h *RankingHandler) GetRankingList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 验证参数
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 从数据库获取排行榜数据
	rankingUsers, total, err := h.dbService.GetRankingList(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取排行榜失败",
			"data":    nil,
		})
		return
	}

	// 根据接口文档，直接返回数组，分页信息放在外层
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rankingUsers,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}
