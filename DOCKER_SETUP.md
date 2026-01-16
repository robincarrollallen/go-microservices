# Docker Compose 本地开发指南

## 📋 前置要求

- Docker Desktop (包含 Docker Engine 和 Docker Compose)
- Go 1.25+
- Make (可选，用于便捷命令)

## 🚀 快速开始

### 1. 启动所有服务

```bash
# 方式一：使用 Makefile
make up

# 方式二：直接使用 docker-compose
docker-compose up -d
```

### 2. 验证服务

```bash
# 查看运行中的容器
docker-compose ps

# 应该看到两个服务都在运行：
# NAME                COMMAND             SERVICE        STATUS
# tenant-service      ./app               tenant-service   Up (healthy)
# user-service        ./app               user-service     Up (healthy)
```

### 3. 访问服务

```bash
# 本地访问 tenant-service
curl http://localhost:8001/health

# 本地访问 user-service
curl http://localhost:8002/health
```

## 🔍 测试 DNS 解析

Docker 内置 DNS 允许容器通过服务名相互通信。

### 验证 DNS 工作正常

```bash
# 在 tenant-service 容器内测试 DNS 解析
make test-dns

# 或手动运行：
docker-compose exec tenant-service nslookup user-service

# 应该返回 user-service 的容器 IP，例如：
# Name:   user-service
# Address: 172.18.0.3
```

## 📡 测试服务间通信

```bash
make test-api

# 这会测试：
# 1. tenant-service 访问 user-service
# 2. user-service 访问 tenant-service
```

## 📝 查看日志

```bash
# 查看所有服务日志
make logs

# 只查看 tenant-service 日志
make logs-tenant

# 只查看 user-service 日志
make logs-user

# 实时查看特定服务的日志
docker-compose logs -f user-service
```

## 🐚 进入容器调试

```bash
# 进入 tenant-service 容器
make shell-tenant

# 进入 user-service 容器
make shell-user

# 在容器内可以运行：
curl http://user-service:8080/health
nslookup user-service
ps aux
```

## 🛑 停止和清理

```bash
# 停止所有服务（保留数据卷）
make down

# 停止所有服务并删除数据卷
make down-clean

# 清理所有容器和未使用的镜像
make clean
```

## 🔧 常见命令速查

| 命令 | 作用 |
|------|------|
| `make build` | 构建镜像 |
| `make up` | 启动服务 |
| `make down` | 停止服务 |
| `make ps` | 查看容器状态 |
| `make logs` | 查看日志 |
| `make shell-tenant` | 进入 tenant 容器 |
| `docker-compose restart tenant-service` | 重启某个服务 |
| `docker-compose build --no-cache` | 强制重建镜像 |

## 🌐 网络配置说明

所有服务都在同一个自定义网络 `microservices` 中：

```
┌─────────────────────────────────┐
│  Docker 网络: microservices      │
├─────────────────────────────────┤
│  ┌──────────────┐  ┌───────────┐│
│  │tenant-service│  │user-service││
│  │  :8001       │  │  :8002     ││
│  │ 容器内:8080  │  │容器内:8080 ││
│  └──────────────┘  └───────────┘│
│        ↕ DNS 解析         ↕      │
│  user-service:8080    tenant     │
│                    -service:8080 │
└─────────────────────────────────┘
```

**在容器内访问其他服务时，使用服务名而不是 localhost：**
- ❌ http://localhost:8080
- ✅ http://user-service:8080

## 💡 环境变量说明

### tenant-service
- `USER_SERVICE_URL`: user-service 的地址（Docker 内自动设置为 `http://user-service:8080`）
- `LOG_LEVEL`: 日志级别（debug/info/warn/error）
- `SERVICE_NAME`: 服务名称

### user-service
- `LOG_LEVEL`: 日志级别
- `SERVICE_NAME`: 服务名称

## 🧪 性能调试

### 查看网络连接
```bash
docker-compose exec tenant-service netstat -tuln
```

### 进行网络延迟测试
```bash
docker-compose exec tenant-service ping -c 4 user-service
```

### 查看容器资源使用情况
```bash
docker stats
```

## 📊 健康检查

每个服务都配置了健康检查：
- 检查间隔：10 秒
- 超时时间：5 秒
- 重试次数：3 次
- 启动等待：10 秒

查看健康状态：
```bash
docker-compose ps
# STATUS 列会显示 Up (healthy) 或 Up (unhealthy)
```

## 🆘 故障排除

### 问题 1: DNS 无法解析
```bash
# 检查网络配置
docker network ls
docker network inspect go-microservices_microservices

# 检查容器是否在正确的网络中
docker inspect tenant-service | grep Networks -A 10
```

### 问题 2: 连接被拒绝 (Connection refused)
```bash
# 确保服务已完全启动
docker-compose ps

# 查看日志找出错误
make logs

# 检查端口是否被占用
lsof -i :8080
```

### 问题 3: 构建失败
```bash
# 清理之前的构建并重试
docker-compose build --no-cache

# 查看构建日志
docker-compose build --verbose
```

### 问题 4: 容器无法访问互联网
```bash
# 在容器内测试网络
docker-compose exec tenant-service wget -O - https://google.com
```

## 🔗 参考资源

- [Docker Compose 官方文档](https://docs.docker.com/compose/)
- [Docker 网络官方文档](https://docs.docker.com/network/)
- [Docker DNS 服务](https://docs.docker.com/network/network-tutorial-standalone/)

## 📌 最佳实践

1. **使用命名网络**：自动启用 DNS 解析
2. **设置健康检查**：确保服务正常运行
3. **使用 depends_on**：控制启动顺序
4. **环境变量隔离**：便于配置切换
5. **多阶段构建**：减小镜像大小
6. **.dockerignore**：加快构建速度

## 📞 其他有用的命令

```bash
# 查看服务的完整输出
docker-compose logs --timestamps --tail=100 tenant-service

# 只查看最后 50 行日志
docker-compose logs --tail=50

# 导出日志到文件
docker-compose logs > docker-compose.log

# 进入容器执行单条命令
docker-compose exec tenant-service ls -la

# 重建特定服务
docker-compose up --build tenant-service
```

