.PHONY: help build up down logs logs-tenant logs-user shell-tenant shell-user test clean

help:
	@echo "Docker Compose 微服务开发命令"
	@echo "================================"
	@echo "make build          - 构建所有镜像"
	@echo "make up             - 启动所有服务"
	@echo "make down           - 停止所有服务"
	@echo "make down-clean     - 停止所有服务并删除数据卷"
	@echo "make logs           - 查看所有服务日志"
	@echo "make logs-tenant    - 查看 tenant-service 日志"
	@echo "make logs-user      - 查看 user-service 日志"
	@echo "make shell-tenant   - 进入 tenant-service 容器"
	@echo "make shell-user     - 进入 user-service 容器"
	@echo "make test-dns       - 测试 DNS 解析"
	@echo "make test-api       - 测试服务间通信"
	@echo "make clean          - 清理容器和镜像"

build:
	docker-compose build

up:
	docker-compose up -d
	@echo "✅ 所有服务已启动"
	@echo "Tenant Service: http://localhost:8001"
	@echo "User Service:   http://localhost:8002"

down:
	docker-compose down

down-clean:
	docker-compose down -v

logs:
	docker-compose logs -f

logs-tenant:
	docker-compose logs -f tenant-service

logs-user:
	docker-compose logs -f user-service

shell-tenant:
	docker-compose exec tenant-service sh

shell-user:
	docker-compose exec user-service sh

test-dns:
	@echo "测试 DNS 解析..."
	docker-compose exec tenant-service nslookup user-service
	docker-compose exec user-service nslookup tenant-service

test-api:
	@echo "测试服务间通信..."
	@echo "\n1️⃣ 从 tenant-service 访问 user-service"
	docker-compose exec tenant-service curl -v http://user-service:8080/health || true
	@echo "\n2️⃣ 从 user-service 访问 tenant-service"
	docker-compose exec user-service curl -v http://tenant-service:8080/health || true

clean:
	docker-compose down -v
	docker system prune -f

ps:
	docker-compose ps

