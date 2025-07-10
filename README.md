# Compeculate API 后端
## 目录

- [功能特性](#功能特性)
- [项目结构](#项目结构)
- [使用指南](#使用指南)
- [API接口文档](#api接口文档)
- [配置说明](#配置说明)
- [开发指南](#开发指南)
- [故障排除](#故障排除)

## 功能特性

- **用户认证系统**: 基于JWT的用户登录/注册
- **用户统计管理**: 答题总数、正确率、分数统计
- **答题结果处理**: 提交答题结果并更新用户统计
- **实时排行榜**: 支持分页的排行榜功能
- **数据持久化**: MySQL数据库存储
- **模块化架构**: 清晰的分层设计
- **完整错误处理**: 统一的错误响应格式
- **RESTful API**: 标准的REST API设计

## 项目结构

```
WechatGo/
├── config/              # 配置管理
│   └── config.go        # 应用配置和环境变量
├── handlers/            # API处理器
│   ├── auth.go          # 认证相关处理器
│   ├── user.go          # 用户相关处理器
│   ├── quiz.go          # 答题相关处理器
│   └── ranking.go       # 排行榜相关处理器
├── middleware/          # 中间件
│   └── auth.go          # JWT认证中间件
├── models/              # 数据模型
│   └── user_score.go    # 用户分数数据模型
├── services/            # 业务服务
│   └── database.go      # 数据库操作服务
├── utils/               # 工具包
│   └── jwt.go           # JWT工具函数
├── main.go              # 主程序入口
├── go.mod               # Go模块文件
├── go.sum               # 依赖校验文件
├── API_DOCUMENTATION.md # 详细API接口文档
├── test_api.py          # API测试脚本
└── README.md            # 项目说明文档
```
## 使用指南
### 1. 配置数据库连接
#### 在config.go中进行数据库连接相关配置
```
func LoadConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "release"),
			ReadTimeout:  10,
			WriteTimeout: 10,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     3306,
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "yyf@0221"),
			Database: getEnv("DB_DATABASE", "weChat"),
			Charset:  "utf8mb4",
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			ExpireTime: 24, // 24小时
		},
	}
```
### 2.使用init.sql文件创建数据库
```bash
# 直接执行init.sql
mysql -u root -p < init.sql
```
### 3. 运行服务

```bash
# 开发模式运行
go run main.go

# 或者编译后运行
go build -o wechatgo-api main.go
./wechatgo-api
```

服务启动后，访问 `http://localhost:8080` 即可。

## API接口文档

### 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: Bearer Token (JWT)
- **请求头**: `Authorization: Bearer <token>`
- **响应格式**: JSON

### 接口列表

#### 1. 用户认证

**用户登录/注册**
```http
POST /api/user/login
Content-Type: application/json

{
  "nickName": "用户昵称",
  "avatarUrl": "头像URL"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "用户昵称",
      "avatarUrl": "头像URL",
      "totalQuestions": 0,
      "correctAnswers": 0,
      "score": 0,
      "accuracy": 0.0
    }
  }
}
```

#### 2. 用户相关

**获取用户统计信息**
```http
GET /api/user/stats
Authorization: Bearer <token>
```

**获取用户详细信息**
```http
GET /api/user/profile
Authorization: Bearer <token>
```

#### 3. 答题相关

**提交答题结果**
```http
POST /api/quiz/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "totalQuestions": 10,
  "correctAnswers": 8,
  "totalTime": 300,
  "score": 80
}
```

#### 4. 排行榜

**获取排行榜列表**
```http
GET /api/ranking/list?page=1&limit=20
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "张三",
        "avatarUrl": "https://example.com/avatar1.jpg",
        "totalQuestions": 20,
        "correctAnswers": 18,
        "score": 180,
        "accuracy": 90.0,
        "rank": 1
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

### 响应格式

所有API接口都使用统一的响应格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {
  }
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 配置说明

### 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| SERVER_PORT | 8080 | 服务器端口 |
| GIN_MODE | debug | Gin运行模式 (debug/release) |
| DB_HOST | localhost | 数据库主机地址 |
| DB_USERNAME | root | 数据库用户名 |
| DB_PASSWORD | - | 数据库密码 |
| DB_DATABASE | weChat | 数据库名称 |
| JWT_SECRET | - | JWT签名密钥 |

### 性能优化

1. **数据库优化**:
   - 为常用查询字段添加索引
   - 配置数据库连接池
   - 定期优化数据库

2. **应用优化**:
   - 使用Gin的release模式
   - 配置适当的并发数
   - 启用压缩中间件

3. **服务器优化**:
   - 使用反向代理（Nginx）
   - 配置SSL证书
   - 启用HTTP/2

## 开发指南

### 项目架构

项目采用分层架构设计：

```
┌─────────────────┐
│   Handlers      │  ← HTTP请求处理
├─────────────────┤
│   Middleware    │  ← 中间件（认证、日志等）
├─────────────────┤
│   Services      │  ← 业务逻辑
├─────────────────┤
│   Models        │  ← 数据模型
├─────────────────┤
│   Database      │  ← 数据库操作
└─────────────────┘
```

### 开发流程

#### 添加新的API接口

1. 在 `handlers/` 目录下创建新的处理器文件
2. 在 `models/` 目录下定义相关数据模型
3. 在 `services/` 目录下添加业务逻辑
4. 在 `main.go` 中注册新的路由
5. 添加相应的测试用例

#### 添加新的数据库表

1. 在 `models/` 目录下创建新的模型文件
2. 在 `services/database.go` 中添加相关操作方法
3. 更新数据库迁移逻辑
4. 添加相应的索引优化

#### 添加新的中间件

1. 在 `middleware/` 目录下创建新的中间件文件
2. 在 `main.go` 中注册中间件
3. 测试中间件功能

### 代码规范

1. **命名规范**:
   - 包名使用小写字母
   - 函数名使用驼峰命名
   - 常量使用大写字母

2. **错误处理**:
   - 所有错误都要处理
   - 使用统一的错误响应格式
   - 记录详细的错误日志

3. **注释规范**:
   - 公共函数必须有注释
   - 复杂逻辑需要详细注释
   - 使用中文注释

## 故障排除

### 常见问题

#### 1. 数据库连接失败

**错误信息**: `dial tcp: connect: connection refused`

**解决方案**:
```bash
# 检查MySQL服务状态
sudo systemctl status mysql

# 启动MySQL服务
sudo systemctl start mysql

# 检查端口是否开放
netstat -tlnp | grep 3306
```

#### 2. 端口被占用

**错误信息**: `listen tcp :8080: bind: address already in use`

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <进程ID>

# 或者修改端口
export SERVER_PORT=8081
```

#### 3. JWT认证失败

**错误信息**: `401 Unauthorized`

**解决方案**:
- 检查JWT_SECRET环境变量是否正确设置
- 确认token格式是否正确
- 检查token是否过期

#### 4. 数据库表不存在

**错误信息**: `Table 'weChat.user_score' doesn't exist`

**解决方案**:
```sql
-- 检查数据库是否存在
SHOW DATABASES;

-- 使用数据库
USE weChat;

-- 检查表是否存在
SHOW TABLES;

-- 如果表不存在，执行创建语句
```

### 日志分析

#### 查看应用日志

```bash
# 查看实时日志
tail -f app.log

# 查看错误日志
grep "ERROR" app.log

# 查看特定时间段的日志
sed -n '/2024-01-01 10:00:00/,/2024-01-01 11:00:00/p' app.log
```

#### 查看系统日志

```bash
# 查看系统日志
sudo journalctl -u wechatgo -f

# 查看最近的日志
sudo journalctl -u wechatgo --since "1 hour ago"
```
