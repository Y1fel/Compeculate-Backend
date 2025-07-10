-- Compeculate API 数据库初始化脚本
-- 数据库: weChat

-- 创建数据库
CREATE DATABASE IF NOT EXISTS weChat DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE weChat;

-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS quiz_result;
DROP TABLE IF EXISTS user_score;

-- 创建用户分数表 (user_score)
-- 对应 models.UserScore 结构
CREATE TABLE user_score (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '唯一ID',
    username VARCHAR(100) NOT NULL COMMENT '用户名称',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    total_questions INT DEFAULT 0 COMMENT '答题总数',
    correct_answers INT DEFAULT 0 COMMENT '正确题目总数',
    score INT DEFAULT 0 COMMENT '分数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username),
    INDEX idx_score (score DESC),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户分数表';

-- 创建答题结果表 (quiz_result)
-- 对应 models.QuizResult 结构
CREATE TABLE quiz_result (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '唯一ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    total_questions INT NOT NULL COMMENT '答题总数',
    correct_answers INT NOT NULL COMMENT '正确题目数',
    score INT NOT NULL COMMENT '本次得分',
    total_time INT NOT NULL COMMENT '答题用时(秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_score (score DESC),
    FOREIGN KEY (user_id) REFERENCES user_score(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='答题结果表';

-- 插入初始用户数据
INSERT INTO user_score (username, avatar_url, total_questions, correct_answers, score) VALUES
('张三', 'https://example.com/avatar1.jpg', 20, 18, 180),
('李四', 'https://example.com/avatar2.jpg', 15, 14, 140),
('王五', 'https://example.com/avatar3.jpg', 25, 22, 220),
('赵六', 'https://example.com/avatar4.jpg', 10, 9, 90),
('钱七', 'https://example.com/avatar5.jpg', 30, 27, 270),
('孙八', 'https://example.com/avatar6.jpg', 0, 0, 0),
('周九', 'https://example.com/avatar7.jpg', 5, 5, 50),
('吴十', 'https://example.com/avatar8.jpg', 12, 11, 110),
('郑十一', 'https://example.com/avatar9.jpg', 18, 16, 160),
('王十二', 'https://example.com/avatar10.jpg', 8, 7, 70),
('冯十三', 'https://example.com/avatar11.jpg', 40, 40, 400),
('陈十四', 'https://example.com/avatar12.jpg', 10, 2, 20),
('褚十五', 'https://example.com/avatar13.jpg', 50, 25, 250),
('卫十六', 'https://example.com/avatar14.jpg', 3, 1, 10),
('蒋十七', 'https://example.com/avatar15.jpg', 1, 0, 0);

-- 插入初始答题结果数据
INSERT INTO quiz_result (user_id, total_questions, correct_answers, score, total_time) VALUES
(1, 10, 9, 90, 300), (1, 10, 9, 90, 280),
(2, 10, 8, 80, 320), (2, 5, 6, 60, 150),
(3, 10, 10, 100, 250), (3, 10, 9, 90, 290), (3, 5, 3, 30, 180),
(4, 10, 9, 90, 310),
(5, 10, 10, 100, 240), (5, 10, 9, 90, 260), (5, 10, 8, 80, 300),
(6, 5, 0, 0, 200),
(7, 5, 5, 50, 120),
(8, 6, 5, 50, 180), (8, 6, 6, 60, 160),
(9, 9, 8, 80, 210), (9, 9, 8, 80, 220),
(10, 8, 7, 70, 200),
(11, 20, 20, 200, 400), (11, 20, 20, 200, 390),
(12, 10, 2, 20, 300),
(13, 25, 12, 120, 500), (13, 25, 13, 130, 480),
(14, 3, 1, 10, 60),
(15, 1, 0, 0, 30);

-- 显示表结构
DESCRIBE user_score;
DESCRIBE quiz_result;

-- 显示用户数据（按分数排序）
SELECT 
    id,
    username,
    avatar_url,
    total_questions,
    correct_answers,
    score,
    ROUND(CASE 
        WHEN total_questions = 0 THEN 0 
        ELSE (correct_answers * 100.0 / total_questions) 
    END, 1) AS accuracy,
    created_at,
    updated_at
FROM user_score 
ORDER BY score DESC;

-- 显示答题结果数据
SELECT 
    qr.id,
    qr.user_id,
    us.username,
    qr.total_questions,
    qr.correct_answers,
    qr.score,
    qr.total_time,
    qr.created_at
FROM quiz_result qr
JOIN user_score us ON qr.user_id = us.id
ORDER BY qr.created_at DESC;

-- 测试排行榜查询（模拟代码中的查询逻辑）
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
LIMIT 10;
