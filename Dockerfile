# 构建容器
FROM golang:1.24 AS builder
# 声明构建目标
ARG TARGET

WORKDIR /src

# 1) 安装 wire（固定版本，避免 latest 触发缓存失效）
RUN --mount=type=cache,id=gomod,sharing=locked,target=/go/pkg/mod \
    go install github.com/google/wire/cmd/wire@latest

# 2) 只拷依赖声明，保证高命中
COPY go.mod go.sum ./

# 3) 共享 module 缓存（跨服务/多构建也可复用）
RUN --mount=type=cache,id=gomod,sharing=locked,target=/go/pkg/mod \
    go mod download

COPY . .

# 4) 复用缓存构建
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    make build-${TARGET}

# 运行时容器
FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
