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

## 开发环境准备

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

- 创建PostgreSQL数据库
- 创建配置文件

```bash
cp .env.example .env
# 编辑.env文件，配置数据库连接信息
```

数据库默认配置:
- 数据库名称: cmdb
- 用户名: admin
- 密码: 123456
- 主机: localhost
- 端口: 5432

### 3. 执行数据库迁移

```bash
# 创建数据库
createdb cmdb

# 执行SQL迁移
psql -d cmdb -f migrations/001_init_schema.sql
```

## 启动应用

```bash
go run cmd/server/main.go
```

应用将在 http://localhost:8080 运行，默认管理员账号/密码：admin/admin123


## 项目结构

```
kk-ops/
├── cmd/                   # 应用入口
│   └── server/            # 服务器入口
├── internal/              # 内部包
│   ├── api/               # API控制器
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── service/           # 业务逻辑
│   ├── config/            # 配置
│   └── utils/             # 工具函数
├── ui/                    # 前端资源
│   ├── css/               # CSS样式
│   ├── js/                # JavaScript脚本
│   └── templates/         # HTML模板
├── migrations/            # 数据库迁移脚本
├── scripts/               # 工具脚本
├── .env.example           # 环境变量示例
├── go.mod                 # Go模块定义
└── README.md              # 项目说明
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

MIT