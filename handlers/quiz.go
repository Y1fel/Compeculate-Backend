package handlers

import (
	"fmt"
	"net/http"
	"time"

	"WechatGo/models"
	"WechatGo/services"

	"github.com/gin-gonic/gin"
)

// QuizSubmitRequest 答题提交请求
type QuizSubmitRequest struct {
	TotalQuestions int `json:"totalQuestions" binding:"required,min=1"`
	CorrectAnswers int `json:"correctAnswers" binding:"required,min=0"`
	TotalTime      int `json:"totalTime" binding:"required,min=0"`
	Score          int `json:"score" binding:"required,min=0"`
}

// QuizSubmitResponse 答题提交响应
type QuizSubmitResponse struct {
	QuizID      string `json:"quizId"`
	SubmittedAt string `json:"submittedAt"`
}

// QuizHandler 答题处理器
type QuizHandler struct {
	dbService *services.DatabaseService
}

// NewQuizHandler 创建答题处理器
func NewQuizHandler(dbService *services.DatabaseService) *QuizHandler {
	return &QuizHandler{
		dbService: dbService,
	}
}

// SubmitQuizResult 提交答题结果
func (h *QuizHandler) SubmitQuizResult(c *gin.Context) {
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

	userID := userIDInterface.(int)

	var req QuizSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证参数逻辑
	if req.CorrectAnswers > req.TotalQuestions {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "正确答案数不能超过总题目数",
			"data":    nil,
		})
		return
	}

	// 保存答题结果到数据库
	quizResult := &models.QuizResult{
		UserID:         userID,
		TotalQuestions: req.TotalQuestions,
		CorrectAnswers: req.CorrectAnswers,
		Score:          req.Score,
		TotalTime:      req.TotalTime,
		CreatedAt:      time.Now(),
	}

	err := h.dbService.SaveQuizResult(quizResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存答题结果失败",
			"data":    nil,
		})
		return
	}

	// 更新用户总分
	err = h.dbService.UpdateUserScore(userID, req.TotalQuestions, req.CorrectAnswers, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户分数失败",
			"data":    nil,
		})
		return
	}

	// 生成响应
	response := QuizSubmitResponse{
		QuizID:      fmt.Sprintf("quiz_%d", quizResult.ID),
		SubmittedAt: quizResult.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "提交成功",
		"data":    response,
	})
}
