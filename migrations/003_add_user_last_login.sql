-- 添加用户最后登录时间字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login TIMESTAMP WITH TIME ZONE;

-- 添加注释
COMMENT ON COLUMN users.last_login IS '用户最后登录时间';