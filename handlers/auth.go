package handlers

import (
	"net/http"
	"time"

	"WechatGo/models"
	"WechatGo/services"
	"WechatGo/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	NickName  string `json:"nickName" binding:"required"`
	AvatarURL string `json:"avatarUrl"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

// UserInfo 用户信息
type UserInfo struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	dbService *services.DatabaseService
	jwtUtil   *utils.JWTUtil
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(dbService *services.DatabaseService, jwtUtil *utils.JWTUtil) *AuthHandler {
	return &AuthHandler{
		dbService: dbService,
		jwtUtil:   jwtUtil,
	}
}

// Login 用户登录/注册
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 检查用户是否已存在
	user, err := h.dbService.GetUserByUsername(req.NickName)
	if err != nil {
		// 用户不存在，创建新用户
		user = &models.UserScore{
			Username:       req.NickName,
			AvatarURL:      req.AvatarURL,
			TotalQuestions: 0,
			CorrectAnswers: 0,
			Score:          0,
		}

		err = h.dbService.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建用户失败",
				"data":    nil,
			})
			return
		}
	}

	// 生成JWT token
	token, err := h.jwtUtil.GenerateToken(user.ID, user.Username, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成token失败",
			"data":    nil,
		})
		return
	}

	response := LoginResponse{
		Token: token,
		UserInfo: UserInfo{
			NickName:  user.Username,
			AvatarURL: user.AvatarURL,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}
