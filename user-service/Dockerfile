# 多阶段构建 - 编译阶段
FROM golang:1.25-alpine AS builder

WORKDIR /build

# 复制 pkg 模块（本地共享代码）
COPY pkg ./pkg

# 复制当前服务的所有代码
WORKDIR /build/user-service
COPY user-service .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server

# 多阶段构建 - 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/user-service/app .
COPY --from=builder /build/user-service/configs ./configs

EXPOSE 8080
CMD ["./app"]

