# 刷题小程序 Go 后端

基于 Go + Gin + GORM + MySQL 的刷题小程序后端服务。

## 技术栈

| 技术 | 说明 |
|------|------|
| Go 1.21+ | 语言 |
| Gin | HTTP 框架 |
| GORM | ORM |
| MySQL 8.0 | 数据库 |
| JWT | 认证 |
| Docker | 容器化部署 |

## 项目结构

```
baokaobao/
├── cmd/server/          # 程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理层
│   ├── middleware/      # 中间件
│   ├── pkg/             # 公共包
│   │   ├── jwt/         # JWT 工具
│   │   ├── response/    # 统一响应
│   │   └── wechat/      # 微信 SDK
│   └── router/          # 路由注册
├── migrations/          # 数据库迁移
├── config/             # 配置文件
└── docs/               # 文档
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Docker (可选)

### 配置

```yaml
# config/config.yaml
app:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  name: baokaobao

jwt:
  secret: your-jwt-secret
  expire_hours: 720

wechat:
  appid: your_wechat_appid
  secret: your_wechat_secret
```

### 运行

```bash
# 安装依赖
go mod tidy

# 编译
go build -o bin/server ./cmd/server

# 运行
./bin/server -config config/config.yaml
```

### Docker 部署

```bash
docker-compose up -d
```

## API 文档

### 小程序 API (`/api/v1/`)

#### 认证
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /auth/login_by_wechat | 微信一键登录 |
| POST | /auth/decrypt_phone | 解密手机号 |
| POST | /auth/logout | 登出 |

#### 用户
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /user/profile | 获取个人信息 |
| PUT | /user/profile | 更新个人信息 |

#### 题目
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /questions | 题目列表 |
| GET | /questions/:id | 题目详情 |
| GET | /questions/random | 随机出题 |

#### 答题
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /quiz/submit | 提交答题 |
| GET | /quiz/history | 答题历史 |
| GET | /quiz/wrong_questions | 错题本 |

#### 得分
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /score/my | 我的得分 |
| GET | /score/ranking | 排行榜 |
| GET | /score/stats | 个人统计 |

### 后台 API (`/admin/api/v1/`)

#### 认证
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /login | 后台登录 |
| POST | /logout | 登出 |

#### 题库管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /question_banks | 题库列表 |
| POST | /question_banks | 创建题库 |
| PUT | /question_banks/:id | 更新题库 |
| DELETE | /question_banks/:id | 删除题库 |

#### 题目管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /questions | 题目列表 |
| POST | /questions | 创建题目 |
| PUT | /questions/:id | 更新题目 |
| DELETE | /questions/:id | 删除题目 |
| POST | /questions/import | 批量导入 |

#### 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /users | 用户列表 |
| GET | /users/:id | 用户详情 |
| PUT | /users/:id/status | 更新状态 |

#### 数据统计
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /dashboard | 数据概览 |
| GET | /stats/overview | 统计概览 |

## 数据库

自动迁移，程序启动时自动创建表。

主要表：
- `users` - 小程序用户
- `admin_users` - 后台用户
- `question_banks` - 题库
- `questions` - 题目
- `question_options` - 题目选项
- `user_answers` - 答题记录
- `wrong_questions` - 错题本
- `scores` - 得分
- `exam_records` - 考试记录
- `user_bank_access` - 用户题库权限

## License

MIT
