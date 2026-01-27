-- 创建额外的数据库
CREATE DATABASE user_db;
CREATE DATABASE tenant_db;

-- 在 tenant_db 中创建 tenants 表
\c tenant_db

CREATE TABLE IF NOT EXISTS tenants (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  status SMALLINT DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT idx_name UNIQUE (name)
);

CREATE INDEX idx_tenants_deleted_at ON tenants(deleted_at);

-- 创建 domains 表
CREATE TABLE IF NOT EXISTS domains (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  domain VARCHAR(255) NOT NULL,
  status SMALLINT DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT fk_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenants(id),
  CONSTRAINT idx_domain UNIQUE (domain)
);

CREATE INDEX idx_domains_tenant_id ON domains(tenant_id);
CREATE INDEX idx_domains_deleted_at ON domains(deleted_at);