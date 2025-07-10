-- 创建数据库
CREATE DATABASE IF NOT EXISTS weChat DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE weChat;

-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS quiz_result;
DROP TABLE IF EXISTS user_score;

-- 创建用户分数表
CREATE TABLE user_score (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '唯一ID',
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

-- 创建答题结果表
CREATE TABLE quiz_result (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '唯一ID',
    user_id INT NOT NULL COMMENT '用户ID',
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

-- 插入测试数据（可选）
INSERT INTO user_score (username, avatar_url, total_questions, correct_answers, score) VALUES
('张三', 'https://example.com/avatar1.jpg', 20, 18, 180),
('李四', 'https://example.com/avatar2.jpg', 15, 14, 140),
('王五', 'https://example.com/avatar3.jpg', 25, 22, 220),
('赵六', 'https://example.com/avatar4.jpg', 10, 9, 90),
('钱七', 'https://example.com/avatar5.jpg', 30, 27, 270),
('孙八', 'https://example.com/avatar6.jpg', 0, 0, 0),
('周九', 'https://example.com/avatar7.jpg', 5, 5, 50);
-- 显示表结构
DESCRIBE user_score;
DESCRIBE quiz_result;

-- 显示测试数据
SELECT * FROM user_score ORDER BY score DESC;
