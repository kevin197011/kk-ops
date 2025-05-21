-- 添加SSH认证相关字段
ALTER TABLE servers ADD COLUMN IF NOT EXISTS ssh_key_path VARCHAR(255);
ALTER TABLE servers ADD COLUMN IF NOT EXISTS auth_type VARCHAR(50) NOT NULL DEFAULT 'password';

-- 将现有记录的认证类型设置为password
UPDATE servers SET auth_type = 'password' WHERE auth_type IS NULL OR auth_type = '';

-- 添加注释
COMMENT ON COLUMN servers.ssh_key_path IS 'SSH密钥文件路径';
COMMENT ON COLUMN servers.auth_type IS '认证类型: password 或 key';