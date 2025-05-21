-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- admin, user
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 创建服务器表
CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    ip VARCHAR(50) NOT NULL,
    port INTEGER NOT NULL DEFAULT 22,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'offline', -- online, offline, maintenance
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 创建CICD流水线表
CREATE TABLE IF NOT EXISTS pipelines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    git_repo VARCHAR(255) NOT NULL,
    branch VARCHAR(100) NOT NULL DEFAULT 'main',
    build_script TEXT,
    deploy_script TEXT,
    server_id INTEGER REFERENCES servers(id),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive', -- inactive, active, running, success, failed
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_run_at TIMESTAMP WITH TIME ZONE
);

-- 创建流水线运行记录表
CREATE TABLE IF NOT EXISTS pipeline_runs (
    id SERIAL PRIMARY KEY,
    pipeline_id INTEGER REFERENCES pipelines(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL, -- pending, running, success, failed
    logs TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    trigger_user INTEGER REFERENCES users(id)
);

-- 创建仪表盘统计表
CREATE TABLE IF NOT EXISTS dashboard_stats (
    id SERIAL PRIMARY KEY,
    stat_date DATE NOT NULL,
    server_count INTEGER NOT NULL DEFAULT 0,
    online_server_count INTEGER NOT NULL DEFAULT 0,
    pipeline_count INTEGER NOT NULL DEFAULT 0,
    pipeline_success_count INTEGER NOT NULL DEFAULT 0,
    pipeline_failed_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- 插入默认管理员用户 (密码: 123456)
INSERT INTO users (username, password, email, role, created_at, updated_at)
VALUES (
    'admin',
    '$2a$10$GkqLnj8StDgj6ho7BHWS.u0P/H0QuZwWE81.85qawnVt0/Zv.Cebi',
    'admin@example.com',
    'admin',
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;