FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装基本依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kk-ops ./cmd/server

# 创建精简的最终镜像
FROM alpine:3.17

# 安装基本工具和证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

WORKDIR /app

# 创建必要的目录
RUN mkdir -p /app/ui/templates /app/ui/static/css /app/ui/static/js /app/migrations

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/kk-ops /app/

# 复制必要的资源文件
COPY --from=builder /app/ui/templates /app/ui/templates/
COPY --from=builder /app/ui/static /app/ui/static/
COPY --from=builder /app/migrations /app/migrations/

# 暴露应用端口
EXPOSE 8080

# 启动应用
CMD ["/app/kk-ops"]