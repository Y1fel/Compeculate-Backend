package services

import (
	"fmt"
	"log"

	"WechatGo/config"
	"WechatGo/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseService 数据库服务
type DatabaseService struct {
	DB *gorm.DB
}

// NewDatabaseService 创建数据库服务
func NewDatabaseService(cfg *config.Config) (*DatabaseService, error) {
	dsn := cfg.Database.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.UserScore{}, &models.QuizResult{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully")
	return &DatabaseService{DB: db}, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *DatabaseService) GetUserByUsername(username string) (*models.UserScore, error) {
	var user models.UserScore
	err := s.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户
func (s *DatabaseService) CreateUser(user *models.UserScore) error {
	return s.DB.Create(user).Error
}

// UpdateUserScore 更新用户分数
func (s *DatabaseService) UpdateUserScore(userID int, totalQuestions, correctAnswers, score int) error {
	return s.DB.Model(&models.UserScore{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"total_questions": gorm.Expr("total_questions + ?", totalQuestions),
			"correct_answers": gorm.Expr("correct_answers + ?", correctAnswers),
			"score":           gorm.Expr("score + ?", score),
		}).Error
}

// GetUserStats 获取用户统计信息
func (s *DatabaseService) GetUserStats(userID int) (*models.UserScore, error) {
	var user models.UserScore
	err := s.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetRankingList 获取排行榜列表
func (s *DatabaseService) GetRankingList(page, limit int) ([]models.RankingUser, int64, error) {
	var total int64

	// 获取总数
	err := s.DB.Model(&models.UserScore{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 使用原生SQL查询获取排行榜数据
	offset := (page - 1) * limit
	query := `
		SELECT 
			RANK() OVER (ORDER BY score DESC) AS ranking,
			id,
			username,
			avatar_url,
			total_questions,
			correct_answers,
			score,
			ROUND(CASE 
				WHEN total_questions = 0 THEN 0 
				ELSE (correct_answers * 100.0 / total_questions) 
			END, 1) AS accuracy
		FROM user_score
		ORDER BY score DESC
		LIMIT ? OFFSET ?
	`

	var rankingUsers []models.RankingUser
	err = s.DB.Raw(query, limit, offset).Scan(&rankingUsers).Error
	if err != nil {
		return nil, 0, err
	}

	return rankingUsers, total, nil
}

// GetUserRank 获取用户排名
func (s *DatabaseService) GetUserRank(userID int) (int, error) {
	// 使用原生SQL查询获取用户排名
	query := `
		SELECT ranking FROM (
			SELECT 
				id,
				RANK() OVER (ORDER BY score DESC) AS ranking
			FROM user_score
		) ranked_users
		WHERE id = ?
	`

	var rank int
	err := s.DB.Raw(query, userID).Scan(&rank).Error
	if err != nil {
		return 0, err
	}

	return rank, nil
}

// SaveQuizResult 保存答题结果
func (s *DatabaseService) SaveQuizResult(result *models.QuizResult) error {
	return s.DB.Create(result).Error
}

// GetUserByID 根据ID获取用户
func (s *DatabaseService) GetUserByID(userID int) (*models.UserScore, error) {
	var user models.UserScore
	err := s.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
