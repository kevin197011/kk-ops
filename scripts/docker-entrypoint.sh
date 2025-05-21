#!/bin/bash
set -e

# 运行数据库迁移脚本
echo "Running database migrations..."
/app/scripts/db_init.sh

# 启动应用
echo "Starting kk-ops application..."
exec "$@"