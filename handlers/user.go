package handlers

import (
	"log"
	"net/http"

	"WechatGo/services"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	dbService *services.DatabaseService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(dbService *services.DatabaseService) *UserHandler {
	return &UserHandler{
		dbService: dbService,
	}
}

// GetUserStats 获取用户统计信息
func (h *UserHandler) GetUserStats(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
			"data":    nil,
		})
		return
	}

	userID := userIDInterface.(int64)

	// 从数据库获取用户统计信息
	user, err := h.dbService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户统计信息失败",
			"data":    nil,
		})
		return
	}

	// 获取用户排名
	rank, err := h.dbService.GetUserRank(userID)
	if err != nil {
		rank = 0 // 如果获取排名失败，设为0
	}

	response := gin.H{
		"totalQuizzes": user.TotalQuestions,
		"totalScore":   user.Score,
		"accuracy":     user.GetAccuracy(),
		"ranking":      rank,
	}
	log.Print(response)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// GetUserProfile 获取用户详细信息
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
			"data":    nil,
		})
		return
	}

	userID := userIDInterface.(int64)

	// 从数据库获取用户详细信息
	user, err := h.dbService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
			"data":    nil,
		})
		return
	}

	response := gin.H{
		"nickName":  user.Username,
		"avatarUrl": user.AvatarURL,
	}
	log.Print(response)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}
