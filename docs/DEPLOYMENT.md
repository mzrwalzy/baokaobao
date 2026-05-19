# 刷题小程序 Go 语言重写 - 部署文档 & 功能规划

---

## 一、系统架构图

```
                        ┌─────────────────────────────────────────────────────────┐
                        │                      微信小程序                            │
                        │         (一键登录 / 刷题 / 答题 / 排行榜)                    │
                        └──────────────────────────┬────────────────────────────────┘
                                                   │ HTTPS (443)
                                                   ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│                               Nginx (反向代理 + SSL)                                │
│                     监听 :80/:443  分发到后端服务                                    │
└────────────────────────┬───────────────────────────────────┬─────────────────────┘
                         │                                   │
               ┌─────────▼─────────┐               ┌─────────▼─────────┐
               │   Go API 服务      │               │   Go Admin 服务    │
               │   (小程序用户)     │               │   (后台管理)        │
               │   :8080           │               │   :8081            │
               └─────────┬─────────┘               └─────────┬─────────┘
                         │                                   │
                         └──────────────┬───────────────────┘
                                        │
                              ┌─────────▼─────────┐
                              │      MySQL 8.0     │
                              │     (3306)         │
                              └───────────────────┘
```

---

## 二、部署拓扑

### 2.1 服务器规划

| 主机 | 角色 | 配置 | 数量 |
|------|------|------|------|
| API Server | Go API + Admin | 2核4G+ | 1-2 台 |
| Database | MySQL 8.0 | 2核4G+ 50G+ SSD | 1 台 |
| Nginx | 反向代理 | 1核1G+ | 1 台 |

> **说明**：初期可单台部署，流量增长后拆分

### 2.2 端口规划

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 80/443 | 对外入口 |
| Go API | 8080 | 小程序 API |
| Go Admin | 8081 | 后台管理 API |
| MySQL | 3306 | 数据库 |

---

## 三、环境要求

### 3.1 软件依赖

| 软件 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 推荐 1.21 LTS |
| MySQL | 8.0+ | 推荐 8.0.35+ |
| Nginx | 1.18+ | 反向代理 |
| Redis | 6.0+ | 可选 (缓存/会话) |

### 3.2 开发工具

| 工具 | 用途 |
|------|------|
| GoLand / VSCode | IDE |
| Navicat / DBeaver | 数据库客户端 |
| Postman / Apifox | API 测试 |
| Docker | 容器化部署 |

---

## 四、配置文件说明

### 4.1 config.yaml 结构

```yaml
# 小程序 API 服务配置
app:
  name: "baokaobao-api"
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug / release / test

# 数据库配置
database:
  host: "127.0.0.1"
  port: 3306
  username: "baokaobao"
  password: "your_password"
  name: "baokaobao"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

# JWT 配置
jwt:
  secret: "your-jwt-secret-key-change-in-production"
  expire_hours: 720  # 30天

# 微信小程序配置
wechat:
  appid: "wx_your_appid"
  secret: "your_app_secret"
  mchid: ""  # 微信商户号 (可选)
  apikey: "" # 微信支付API密钥 (可选)

# 管理后台配置
admin:
  host: "0.0.0.0"
  port: 8081

# 日志配置
log:
  level: "info"
  path: "./logs"
  max_size: 100
  max_backups: 30
```

### 4.2 环境变量覆盖 (生产环境)

```bash
# 推荐通过环境变量覆盖敏感配置
export BAOKAOBAO_JWT_SECRET="production-secret"
export BAOKAOBAO_DATABASE_PASSWORD="prod-password"
export BAOKAOBAO_WECHAT_SECRET="prod-wechat-secret"
```

---

## 五、数据库设计

### 5.1 ER 图

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│ admin_users  │       │    users     │       │  categories  │
├──────────────┤       ├──────────────┤       ├──────────────┤
│ id           │       │ id           │       │ id           │
│ username     │       │ openid       │       │ name         │
│ password_hash│       │ unionid      │       │ parent_id    │
│ nickname     │       │ nickname     │       │ sort         │
│ role         │       │ avatar_url   │       └──────┬───────┘
│ status       │       │ phone        │              │
└──────────────┘       │ is_admin     │              │
                       │ status       │              │
                       └──────┬───────┘              │
                              │                      │
                              │ 1:N                  │ 1:N
                              ▼                      ▼
                       ┌──────────────┐       ┌──────────────┐
                       │user_answers  │       │  questions   │
                       ├──────────────┤       ├──────────────┤
                       │ id           │       │ id           │
                       │ user_id      │       │ category_id  │
                       │ question_id  │       │ title        │
                       │ my_answer    │       │ content      │
                       │ is_correct   │       │ answer       │
                       │ score        │       │ analysis     │
                       │ answered_at  │       │ difficulty   │
                       └──────┬───────┘       │ type         │
                              │               │ score        │
                              │               └──────┬───────┘
                              │                      │
                              │ N:1                   │
                              ▼                      │
                       ┌──────────────┐               │
                       │   scores    │◄──────────────┘
                       ├──────────────┤
                       │ id           │
                       │ user_id      │
                       │ total_score  │
                       │ quiz_date    │
                       └──────────────┘
```

### 5.2 表结构汇总

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `users` | 小程序用户 | openid, unionid, phone |
| `admin_users` | 后台用户 | username, password_hash, role |
| `categories` | 题目分类 | name, parent_id (树形) |
| `questions` | 题目库 | title, content, answer, difficulty |
| `user_answers` | 答题记录 | user_id, question_id, my_answer |
| `scores` | 得分记录 | user_id, total_score, quiz_date |
| `question_tags` | 题目标签 | tag_name |
| `question_tag_relations` | 题目标签关联 | question_id, tag_id |

---

## 六、API 路由总览

### 6.1 小程序 API (`/api/v1/`)

#### 认证模块
| 方法 | 路径 | 说明 | 需认证 |
|------|------|------|--------|
| POST | `/auth/login_by_wechat` | 微信一键登录 | ❌ |
| POST | `/auth/decrypt_phone` | 解密手机号 | ✅ |
| POST | `/auth/logout` | 登出 | ✅ |

#### 用户模块
| 方法 | 路径 | 说明 | 需认证 |
|------|------|------|--------|
| GET | `/user/profile` | 获取个人信息 | ✅ |
| PUT | `/user/profile` | 更新个人信息 | ✅ |
| POST | `/user/avatar` | 上传头像 | ✅ |

#### 题目模块
| 方法 | 路径 | 说明 | 需认证 |
|------|------|------|--------|
| GET | `/categories` | 获取分类树 | ✅ |
| GET | `/questions` | 获取题目列表 | ✅ |
| GET | `/questions/:id` | 题目详情 | ✅ |
| GET | `/questions/random` | 随机出题 | ✅ |

#### 答题模块
| 方法 | 路径 | 说明 | 需认证 |
|------|------|------|--------|
| POST | `/quiz/submit` | 提交答题 | ✅ |
| GET | `/quiz/history` | 答题历史 | ✅ |
| GET | `/quiz/wrong_questions` | 错题本 | ✅ |
| POST | `/quiz/add_wrong/:qid` | 加入错题本 | ✅ |

#### 得分模块
| 方法 | 路径 | 说明 | 需认证 |
|------|------|------|--------|
| GET | `/score/my` | 我的得分 | ✅ |
| GET | `/score/ranking` | 排行榜 | ✅ |
| GET | `/score/stats` | 个人统计 | ✅ |

### 6.2 后台管理 API (`/admin/api/v1/`)

#### 认证
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 后台登录 |
| POST | `/logout` | 登出 |

#### 题目管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/questions` | 题目列表 |
| POST | `/questions` | 创建题目 |
| PUT | `/questions/:id` | 更新题目 |
| DELETE | `/questions/:id` | 删除题目 |
| POST | `/questions/import` | 批量导入 |

#### 分类管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/categories` | 分类列表 |
| POST | `/categories` | 创建分类 |
| PUT | `/categories/:id` | 更新分类 |
| DELETE | `/categories/:id` | 删除分类 |

#### 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/users` | 用户列表 |
| GET | `/users/:id` | 用户详情 |
| PUT | `/users/:id/status` | 启用/禁用 |

#### 数据统计
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/stats/overview` | 数据概览 |
| GET | `/stats/users` | 用户统计 |
| GET | `/stats/questions` | 题目统计 |
| GET | `/stats/activity` | 活跃度统计 |

---

## 七、微信登录详细流程

### 7.1 一键登录时序

```
┌────────┐    ┌─────────────┐    ┌──────────┐    ┌──────────────┐
│ 小程序  │    │   Go API    │    │ 微信服务器 │    │   MySQL      │
└──┬─────┘    └──────┬──────┘    └────┬─────┘    └──────┬───────┘
   │                 │                │                │
   │ 1.wx.login()   │                │                │
   │────────────────►                │                │
   │                 │                │                │
   │ 2.code          │                │                │
   │◄────────────────│                │                │
   │                 │                │                │
   │ 3.code          │                │                │
   │────────────────►│                │                │
   │                 │ 4.code2session  │                │
   │                 │────────────────►│                │
   │                 │                │                │
   │                 │ 5.openid        │                │
   │                 │   session_key   │                │
   │                 │◄────────────────│                │
   │                 │                │                │
   │                 │ 6.查找/创建用户  │                │
   │                 │────────────────►│ 7.insert      │
   │                 │                │◄──────────────│
   │                 │                │                │
   │                 │ 8.生成JWT       │                │
   │                 │────────────────│                │
   │ 9.token         │                │                │
   │◄────────────────│                │                │
```

### 7.2 关键安全说明

```
⚠️  session_key 处理策略：
   - 仅用于解密手机号等敏感数据
   - 解密后立即丢弃，不存储
   - 每次登录重新获取

✅  openid 存储：
   - 用 openid 作为用户唯一标识
   - 可重复使用，不失效

✅  JWT Token：
   - 包含 user_id, openid, exp
   - 有效期 30 天
   - 自动续期机制
```

---

## 八、Docker 部署

### 8.1 docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: baokaobao-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: baokaobao
      MYSQL_USER: baokaobao
      MYSQL_PASSWORD: baokaobao_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    command: --default-authentication-plugin=mysql_native_password
    networks:
      - baokaobao-net

  api:
    build:
      context: .
      dockerfile: Dockerfile.api
    container_name: baokaobao-api
    restart: always
    environment:
      BAOKAOBAO_DATABASE_HOST: mysql
      BAOKAOBAO_DATABASE_PASSWORD: baokaobao_password
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    networks:
      - baokaobao-net

  admin:
    build:
      context: .
      dockerfile: Dockerfile.admin
    container_name: baokaobao-admin
    restart: always
    environment:
      BAOKAOBAO_DATABASE_HOST: mysql
      BAOKAOBAO_DATABASE_PASSWORD: baokaobao_password
    ports:
      - "8081:8081"
    depends_on:
      - mysql
    networks:
      - baokaobao-net

  nginx:
    image: nginx:1.18
    container_name: baokaobao-nginx
    restart: always
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - api
      - admin
    networks:
      - baokaobao-net

volumes:
  mysql_data:

networks:
  baokaobao-net:
    driver: bridge
```

### 8.2 Nginx 配置

```nginx
upstream api_backend {
    server api:8080;
}

upstream admin_backend {
    server admin:8081;
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    # 小程序 API
    location /api/ {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 后台管理
    location /admin/ {
        proxy_pass http://admin_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 九、部署检查清单

### 9.1 部署前检查

- [ ] Go 1.21+ 已安装
- [ ] MySQL 8.0 已安装并运行
- [ ] 微信公众平台 AppID 和 AppSecret 已获取
- [ ] 域名已备案并解析
- [ ] SSL 证书已申请

### 9.2 数据库初始化

```bash
# 1. 创建数据库
mysql -u root -p
CREATE DATABASE baokaobao CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'baokaobao'@'%' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON baokaobao.* TO 'baokaobao'@'%';
FLUSH PRIVILEGES;

# 2. 自动迁移 (程序启动时自动创建表)
```

### 9.3 启动服务

```bash
# 编译
go build -o bin/api ./cmd/api
go build -o bin/admin ./cmd/admin

# 启动 API 服务
./bin/api --config config.yaml

# 启动 Admin 服务
./bin/admin --config config.yaml
```

---

## 十、监控与日志

### 10.1 日志配置

```yaml
log:
  level: "info"           # debug / info / warn / error
  path: "./logs"          # 日志目录
  max_size: 100           # 单文件最大 MB
  max_backups: 30        # 保留文件数
  compress: true         # gzip 压缩
```

### 10.2 健康检查

```
GET /health
Response: {"status": "ok", "time": "2024-01-01T00:00:00Z"}
```

---

## 十一、备份策略

| 备份类型 | 频率 | 保留时间 | 说明 |
|----------|------|----------|------|
| 全量备份 | 每天 | 30 天 | mysqldump |
| 增量备份 | 每小时 | 7 天 | binlog |
| 配置备份 | 每次变更 | 90 天 | config.yaml |

```bash
# 备份脚本示例
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -u root -p baokaobao > /backup/baokaobao_$DATE.sql
find /backup -name "*.sql" -mtime +30 -delete
```
