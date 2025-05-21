#!/bin/bash
set -e

# 等待PostgreSQL启动
echo "Waiting for PostgreSQL to start..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "PostgreSQL started, running migrations..."

# 获取所有.sql文件并按照文件名排序
sql_files=$(find /app/migrations -name "*.sql" -type f | sort)

# 执行每个SQL文件
for sql_file in $sql_files; do
  filename=$(basename "$sql_file")
  echo "Applying migration: $filename"
  PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
  echo "Migration $filename applied successfully!"
done

echo "All migrations completed successfully!"