package models

import (
	"time"
)

// UserScore 用户分数模型
type UserScore struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username       string    `json:"username" gorm:"type:varchar(100);not null"`
	AvatarURL      string    `json:"avatar_url" gorm:"type:varchar(255)"`
	TotalQuestions int       `json:"total_questions" gorm:"default:0"`
	CorrectAnswers int       `json:"correct_answers" gorm:"default:0"`
	Score          int       `json:"score" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UserScore) TableName() string {
	return "user_score"
}

// GetAccuracy 计算准确率
func (u *UserScore) GetAccuracy() float64 {
	if u.TotalQuestions == 0 {
		return 0.0
	}
	accuracy := float64(u.CorrectAnswers) / float64(u.TotalQuestions) * 100
	// 保留一位小数
	return float64(int(accuracy*10+0.5)) / 10
}

// UpdateScore 更新用户分数
func (u *UserScore) UpdateScore(totalQuestions, correctAnswers, score int) {
	u.TotalQuestions += totalQuestions
	u.CorrectAnswers += correctAnswers
	u.Score += score
}

// RankingUser 排行榜用户信息（匹配SQL查询结果）
type RankingUser struct {
	Rank           int     `json:"rank" gorm:"column:ranking"`
	ID             int     `json:"id"`
	NickName       string  `json:"nickName" gorm:"column:username"`
	AvatarURL      string  `json:"avatarUrl" gorm:"column:avatar_url"`
	TotalQuizzes   int     `json:"totalQuizzes" gorm:"column:total_questions"`
	CorrectAnswers int     `json:"correct_answers" gorm:"column:correct_answers"`
	TotalScore     int     `json:"totalScore" gorm:"column:score"`
	Accuracy       float64 `json:"accuracy"`
}

// QuizResult 答题结果
type QuizResult struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int       `json:"user_id" gorm:"not null"`
	TotalQuestions int       `json:"total_questions" gorm:"not null"`
	CorrectAnswers int       `json:"correct_answers" gorm:"not null"`
	Score          int       `json:"score" gorm:"not null"`
	TotalTime      int       `json:"total_time" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName 指定表名
func (QuizResult) TableName() string {
	return "quiz_result"
}
