# Kratos Foundation Template

> 基于 [Go-Kratos](https://go-kratos.dev/) 微服务框架的企业级项目模板，提供开箱即用的工程化实践和最佳实践。

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Kratos](https://img.shields.io/badge/Kratos-v2.9+-blue?style=flat)](https://go-kratos.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

---

## 特性

- 🏗️ **清晰架构** - Service/Biz/Data/Client 分层，依赖倒置原则
- 🔧 **Wire 依赖注入** - 编译时依赖注入，类型安全
- 📡 **多协议支持** - HTTP/gRPC/WebSocket 统一框架
- ⏰ **定时任务** - Cron 表达式支持，多种并发策略
- 📊 **可观测性** - OpenTelemetry Tracing/Metrics/Logging
- 🐳 **容器化** - Docker 多阶段构建，支持 Docker Compose
- 🔄 **自动生成** - Protobuf 驱动开发，客户端代码自动生成
- 📝 **配置验证** - JSON Schema 自动生成，IDE 智能提示

---

## 快速开始

### 前置要求

- **Go** >= 1.24
- **Protocol Buffers Compiler** (protoc)
- **Kratos CLI** - `go install github.com/go-kratos/kratos/cmd/kratos/v2@latest`
- **Docker** (可选，用于容器化部署)

### 创建项目

```bash
# 使用模板创建项目
kratos new your-project -r https://github.com/jaggerzhuang1994/kratos-foundation-template

# 进入项目目录
cd your-project

# 安装依赖和工具
make init

# 生成所有代码
make all

# 运行服务
make run
```

### 重命名模块（可选）

如果需要将项目作为私有 Go 模块部署，执行以下命令重命名：

> ⚠️ **注意**：此命令会全局替换所有文件中的模块路径，执行前请确保 Git 工作区干净或已备份。

```bash
make rename-module NEW=github.com/your-org/your-project
```

---

## 项目结构

```
.
├── api/                        # API 定义（Protobuf）
│   └── example_service/        # 示例服务 API
├── cmd/                        # 应用入口
│   └── server/                 # 服务端入口
│       ├── main.go             # 主程序
│       ├── bootstrap.go        # 服务注册与初始化
│       ├── wire.go             # Wire 依赖定义
│       └── wire_gen.go         # Wire 生成代码（勿手动修改）
├── configs/                    # 配置文件
│   ├── config.yaml             # 基础配置
│   └── config.local.yaml       # 本地开发配置（不提交）
├── docker/                     # Docker 相关
│   ├── .env.example            # 环境变量模板
│   └── configs/                # Docker 环境配置
├── internal/                   # 内部代码（不可被外部引用）
│   ├── biz/                    # 业务逻辑层
│   │   └── example/            # 示例业务领域
│   ├── client/                 # 外部服务客户端
│   ├── conf/                   # 配置定义与加载
│   ├── data/                   # 数据访问层
│   ├── job/                    # 定时任务
│   └── service/                # 服务实现层
├── third_party/                # 第三方 proto 文件
├── docker-compose.yaml         # Docker Compose 编排
├── Dockerfile                  # Docker 构建文件
├── Makefile                    # 构建工具
└── openapi.swagger.yaml        # OpenAPI 文档
```

### 分层架构

| 层级          | 目录                 | 职责                             |
|-------------|--------------------|--------------------------------|
| **Service** | `internal/service` | 实现 Protobuf 服务接口，调用 Biz 层      |
| **Biz**     | `internal/biz`     | 核心业务逻辑，定义实体和接口                 |
| **Data**    | `internal/data`    | 数据访问，实现 Biz 层定义的 Repository 接口 |
| **Client**  | `internal/client`  | 外部服务调用，实现 Biz 层定义的 Client 接口   |
| **Conf**    | `internal/conf`    | 配置定义与加载                        |

---

## 开发指南

### 常用命令

| 命令                   | 说明                                |
|----------------------|-----------------------------------|
| `make init`          | 安装依赖工具                            |
| `make all`           | 生成所有代码（API + Config + Wire）+ Lint |
| `make api`           | 生成 API 代码（Protobuf → Go）          |
| `make config`        | 生成配置代码和 JSON Schema               |
| `make generate`      | 生成 Wire 依赖注入代码                    |
| `make run`           | 运行服务（开发模式）                        |
| `make build`         | 构建所有二进制文件                         |
| `make build-server`  | 构建 server 二进制                     |
| `make lint`          | 运行代码检查                            |
| `make test`          | 运行单元测试                            |
| `make test-coverage` | 生成测试覆盖率报告                         |
| `make help`          | 显示所有可用命令                          |

### 创建新服务

```bash
# 1. 创建 proto 文件
kratos proto add api/your_service/v1/your_service.proto

# 2. 生成代码
make all

# 3. 在 internal/service 中实现服务
# 4. 在 internal/biz 中定义业务逻辑
# 5. 在 internal/data 中实现数据访问
# 6. 在 cmd/server/bootstrap.go 中注册服务
```

### 服务注册示例

在 `cmd/server/bootstrap.go` 中注册 HTTP/gRPC/WebSocket 服务：

```go
func NewBootstrap(
httpServer server.HttpServer,
grpcServer server.GrpcServer,
wss server.WebsocketServer,
exampleService *service.ExampleService,
exampleWsHandler *service.ExampleWsHandler,
) bootstrap.Bootstrap {
// WebSocket 路由
wss.Handle("/echo", exampleWsHandler, websocket.Upgrader{
CheckOrigin: func (r *http.Request) bool {
return true // 生产环境应进行严格的来源校验
},
})

// HTTP 服务
example_pb.RegisterExampleServiceHTTPServer(httpServer, exampleService)

// gRPC 服务
example_pb.RegisterExampleServiceServer(grpcServer, exampleService)

return nil
}
```

### WebSocket 处理器

```go
type ExampleWsHandler struct{}

func NewExampleWsHandler() *ExampleWsHandler {
return &ExampleWsHandler{}
}

// 建立连接前校验
func (h *ExampleWsHandler) OnHandshake(r *http.Request) error {
// TODO: 实现 JWT 验证、IP 白名单等逻辑
return nil
}

// 建立连接后回调
func (h *ExampleWsHandler) OnConnect(client server.WebsocketClient) {
log.Info("WebSocket 连接建立", "remote_addr", client.Request().RemoteAddr)
}

// 收到消息后处理
func (h *ExampleWsHandler) OnMessage(client server.WebsocketClient, message []byte, messageType server.MessageType) {
// TODO: 解析消息并调用 Biz 层处理
client.SendText("echo: " + string(message))
}

// 错误处理
func (h *ExampleWsHandler) OnError(client server.WebsocketClient, err error) {
log.Error("WebSocket 错误", "error", err)
}

// 连接关闭后回调
func (h *ExampleWsHandler) OnClose(client server.WebsocketClient) {
log.Info("WebSocket 连接关闭", "remote_addr", client.Request().RemoteAddr)
}
```

---

## 定时任务（Job）

框架内置基于 Cron 的定时任务支持，集成 OpenTelemetry 可观测性。

### 配置定义

在 `configs/config.yaml` 中定义 Job：

```yaml
job:
  jobs:
    # Job 名称（与代码中注册的名称对应）
    sync_data:
      # 调度表达式（支持标准 Cron 和快捷语法）
      schedule: "@every 5s"
      # 并发策略：DELAY（延迟执行）、SKIP（跳过）、OVERLAP（允许并发）
      concurrent_policy: DELAY

    daily_report:
      schedule: "0 2 * * *"  # 每天凌晨 2 点
      concurrent_policy: SKIP
```

**调度表达式：**

| 表达式           | 说明                 |
|---------------|--------------------|
| `@every 5s`   | 每 5 秒执行一次          |
| `@every 1m`   | 每 1 分钟执行一次         |
| `@hourly`     | 每小时整点执行            |
| `@daily`      | 每天午夜执行             |
| `@weekly`     | 每周日午夜执行            |
| `@monthly`    | 每月 1 日午夜执行         |
| `0 2 * * *`   | 标准 Cron：每天 2:00 执行 |
| `*/5 * * * *` | 标准 Cron：每 5 分钟执行   |

**并发策略：**

| 策略        | 说明                    |
|-----------|-----------------------|
| `DELAY`   | 上一次任务未完成时，等待完成后再执行下一次 |
| `SKIP`    | 上一次任务未完成时，跳过本次执行      |
| `OVERLAP` | 允许多个实例并发执行            |

### 实现 Job

```go
package job

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/database"
)

type SyncDataJob struct {
	log log.Log
	db  *database.Manager
}

func NewSyncDataJob(log log.Log, db *database.Manager) *SyncDataJob {
	return &SyncDataJob{log: log, db: db}
}

// Run 实现 job.Job 接口
func (j *SyncDataJob) Run(ctx context.Context) error {
	j.log.Info("开始执行数据同步任务")
	// 业务逻辑...
	j.log.Info("数据同步任务完成")
	return nil
}
```

### 注册 Job

在 `cmd/server/bootstrap.go` 中注册：

```go
func NewBootstrap(
register job.Register,
syncDataJob *job.SyncDataJob,
) bootstrap.Bootstrap {
// 注册 Job（名称需与配置文件中的 key 对应）
register.Register("sync_data", syncDataJob)
return nil
}
```

---

## 依赖倒置实践

### Biz 层定义接口

```go
// internal/biz/user.go
package biz

type User struct {
	ID   int64
	Name string
}

// 数据访问接口（由 Data 层实现）
type UserRepo interface {
	GetUser(ctx context.Context, id int64) (*User, error)
}

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetUser(ctx, id)
}
```

### Data 层实现接口

```go
// internal/data/user_repo.go
package data

import (
	"context"
	"your-project/internal/biz"
)

type UserRepo struct {
	db *database.Manager
}

func NewUserRepo(db *database.Manager) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	// 查询数据库
	return &biz.User{ID: id, Name: "user1"}, nil
}
```

### Wire 绑定接口

```go
// internal/data/wire.go
package data

import (
	"github.com/google/wire"
	"your-project/internal/biz"
)

var _ biz.UserRepo = (*UserRepo)(nil)

var ProviderSet = wire.NewSet(
	NewUserRepo,
	wire.Bind(new(biz.UserRepo), new(*UserRepo)),
)
```

---

## 客户端注入

Protobuf 生成的 `_client.pb.go` 提供自动生成的 Wire Provider：

```go
// internal/client/biz_impl.go
package client

import (
	"github.com/google/wire"
	"your-project/api/example_service/example_pb"
)

type BizImpl struct {
	api example_pb.ExampleServiceApi
}

func NewBizImpl(api example_pb.ExampleServiceApi) *BizImpl {
	return &BizImpl{api: api}
}

var ProviderSet = wire.NewSet(
	example_pb.ExampleServiceApiProvider,
	NewBizImpl,
)
```

---

## Docker 部署

### 本地开发

```bash
# 复制环境变量模板
cp docker/.env.example docker/.env

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 构建镜像

```bash
# 构建指定目标
docker build --build-arg TARGET=server -t your-app:latest .

# 或使用 Makefile
make build-server
```

### 环境配置

**配置优先级：**

- 本地环境：远程配置 > 本地文件
- 其他环境：本地文件 > 远程配置

**配置文件加载顺序：**

```
configs/config.yaml
configs/{app_name}.yaml
configs/{env}.config.yaml
configs/{env}.{app_name}.yaml
```

---

## 配置说明

项目支持通过 `configs/config.yaml` 进行配置，自动生成 JSON Schema 用于 IDE 智能提示和验证。

### 主要配置项

```yaml
# 日志配置
log:
  level: debug              # 日志级别：debug/info/warn/error
  std:
    disable: false          # 是否禁用标准输出
  file:
    disable: false          # 是否禁用文件日志
    path: ./app.log         # 日志文件路径

# 服务配置
server:
  stop_delay: 1s            # 优雅停机延迟
  http:
    addr: 0.0.0.0:8000      # HTTP 服务地址
  grpc:
    addr: 0.0.0.0:9000      # gRPC 服务地址
  middleware:
    timeout:
      default: 1s           # 默认超时
      routes: # 路由级超时
        - path: /example.ExampleService/GetUser
          timeout: 0.5s

# 数据库配置
database:
  default: default          # 默认连接
  connections:
    default:
      dsn: user:pass@tcp(host:3306)/db?parseTime=true
      max_idle_conns: 10
      max_open_conns: 100

# Redis 配置
redis:
  default: default
  connections:
    default:
      addr: localhost:6379
      db: 0

# 客户端配置
client:
  clients:
    example_service:
      protocol: GRPC        # GRPC/HTTP
      target: "127.0.0.1:9000"  # 直连地址或服务发现
      middleware:
        timeout:
          default: 1s

# 定时任务配置
job:
  jobs:
    job_name:
      schedule: "@every 1m"
      concurrent_policy: DELAY
```

---

## 相关资源

- [Go-Kratos 官方文档](https://go-kratos.dev/)
- [Kratos Framework](https://github.com/jaggerzhuang1994/kratos-foundation)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Wire 依赖注入](https://github.com/google/wire)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)

---

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some amazing feature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

---

**Built with ❤️ using [Go-Kratos](https://go-kratos.dev/)
and [Kratos Foundation](https://github.com/jaggerzhuang1994/kratos-foundation)**
