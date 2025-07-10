package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
	jwt.RegisteredClaims
}

// JWTUtil JWT工具
type JWTUtil struct {
	secret string
}

// NewJWTUtil 创建JWT工具实例
func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{
		secret: secret,
	}
}

// GenerateToken 生成JWT token
func (j *JWTUtil) GenerateToken(userID int64, nickName string, expireTime time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		NickName: nickName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ParseToken 解析JWT token
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证token是否有效
func (j *JWTUtil) ValidateToken(tokenString string) bool {
	_, err := j.ParseToken(tokenString)
	return err == nil
}
