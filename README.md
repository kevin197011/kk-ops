# KK-OPS 运维管理平台

一个功能完整的运维管理平台，包含服务器管理、CICD流水线、用户管理等功能。

## 功能模块

- 仪表盘：系统状态概览
- 服务器管理：管理服务器资源
- CICD管理：持续集成/持续部署流水线
- 用户管理：用户权限控制

## 技术栈

- 前端：Bootstrap
- 后端：Golang + Gin框架
- 数据库：PostgreSQL
- 容器化：Docker + Docker Compose

## 开发环境准备

### 方式一：本地开发

#### 1. 安装依赖

```bash
go mod tidy
```

#### 2. 配置数据库

- 创建PostgreSQL数据库
- 创建配置文件

```bash
cp env.example .env
# 编辑.env文件，配置数据库连接信息
```

数据库默认配置:
- 数据库名称: cmdb
- 用户名: admin
- 密码: 123456
- 主机: localhost
- 端口: 5432

#### 3. 执行数据库迁移

```bash
# 创建数据库
createdb cmdb

# 执行SQL迁移
psql -d cmdb -f migrations/001_init_schema.sql
psql -d cmdb -f migrations/002_add_ssh_auth.sql
psql -d cmdb -f migrations/003_add_user_last_login.sql
```

#### 4. 启动应用

```bash
go run cmd/server/main.go
```

应用将在 http://localhost:8080 运行

### 方式二：Docker 部署

#### 1. 使用 Docker Compose 启动

```bash
# 构建并启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

应用将在 http://localhost:8080 运行，PostgreSQL数据库会自动初始化，并且migrations目录下的所有SQL迁移文件会按照文件名顺序自动执行。

## 数据库迁移

### 自动迁移机制

系统支持自动执行数据库迁移脚本：

1. Docker环境下，应用启动时会自动执行`migrations`目录中的所有`.sql`文件
2. 迁移文件按照名称顺序执行（例如：001_xxx.sql, 002_xxx.sql）
3. 可以通过添加新的编号SQL文件来扩展数据库结构

### 添加新的迁移

添加新的迁移步骤：

1. 在`migrations`目录创建新的SQL文件，文件名格式：`XXX_描述.sql`（XXX为数字序号）
2. 编写SQL迁移脚本
3. 重启应用或重新部署Docker容器时，新的迁移将自动执行

## 数据初始化

### 生成测试数据

系统提供了测试数据生成工具，可以快速创建测试用户、服务器和流水线：

```bash
# 运行测试数据生成脚本
go run cmd/mock/main.go
```

生成的测试数据包括：
- 管理员用户（admin/admin123）和普通用户
- 多种环境的测试服务器
- 示例流水线及运行记录

## 项目结构

```
kk-ops/
├── cmd/                   # 应用入口
│   ├── server/            # 服务器入口
│   └── mock/              # 测试数据生成工具
├── internal/              # 内部包
│   ├── api/               # API控制器
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── service/           # 业务逻辑
│   ├── config/            # 配置
│   └── utils/             # 工具函数
├── ui/                    # 前端资源
│   ├── static/            # 静态资源
│   │   ├── css/           # CSS样式
│   │   └── js/            # JavaScript脚本
│   └── templates/         # HTML模板
├── migrations/            # 数据库迁移脚本
├── scripts/               # 工具脚本
│   ├── db_init.sh         # 数据库初始化脚本
│   └── docker-entrypoint.sh # Docker入口点脚本
├── Dockerfile             # Docker镜像构建文件
├── docker-compose.yml     # Docker Compose配置
├── .dockerignore          # Docker忽略文件
├── env.example            # 环境变量示例
├── go.mod                 # Go模块定义
└── README.md              # 项目说明
```

## Docker 环境变量

Docker 环境中可配置的环境变量:

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 应用服务端口 | 8080 |
| DB_HOST | 数据库主机名 | db |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户名 | admin |
| DB_PASSWORD | 数据库密码 | 123456 |
| DB_NAME | 数据库名称 | cmdb |
| JWT_SECRET | JWT密钥 | kk-ops-secure-jwt-token-2025 |

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

MIT