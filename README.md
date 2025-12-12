# Kratos Foundation Template

基于 [Go-Kratos](https://go-kratos.dev/) 微服务框架的项目模板，提供了开箱即用的项目结构和最佳实践。

## 特性

- 基于 Kratos v2 微服务框架
- 清晰的分层架构（Service/Biz/Data/Client）
- Wire 依赖注入
- 完善的 Makefile 工具链
- Protocol Buffers 驱动开发
- 自动生成客户端代码
- WebSocket 支持
- Docker 容器化部署
- Job 定时任务支持
- OpenTelemetry 可观测性（Tracing/Metrics）

## 前置要求

- Go >= 1.24
- Protocol Buffers 编译器
- Kratos CLI 工具
- Docker（可选，用于容器化部署）

## 快速开始

### 1. 安装 Kratos CLI

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### 2. 使用模板创建项目

```bash
# 创建项目
kratos new 你的项目名 -r https://github.com/jaggerzhuang1994/kratos-foundation-template

# 进入项目目录
cd 你的项目名

# 安装依赖
make init

# 生成代码
make all

# 运行服务
make run
```

### 3. 重命名模块名（可选）

如果项目需要部署为 Go 私有包，需要将 go module 命名为仓库名：

> ⚠️ 此工具是全文件文本替换，可能会错误替换非模块名的文本。执行前请备份代码或确保 git 工作空间干净。

```bash
make rename-module NEW=github.com/your-org/your-project
```

## 项目结构

```
.
├── api/                    # API 定义（protobuf）
│   └── example_service/    # 示例服务 API
├── cmd/                    # 应用入口
│   └── server/             # 服务入口
│       ├── main.go         # 主程序
│       ├── bootstrap.go    # 服务注册与初始化
│       ├── wire.go         # Wire 依赖定义
│       └── wire_gen.go     # Wire 生成代码
├── configs/                # 配置文件
├── docker/                 # Docker 相关配置
│   ├── .env.example        # 环境变量模板
│   └── configs/            # Docker 环境配置
├── internal/               # 内部代码
│   ├── biz/                # 业务逻辑层
│   ├── client/             # 外部服务客户端
│   ├── conf/               # 配置定义与加载
│   ├── data/               # 数据访问层
│   └── service/            # 服务实现层
├── third_party/            # 第三方 proto 文件
├── docker-compose.yaml     # Docker Compose 编排
├── Dockerfile              # Docker 构建文件
└── Makefile                # 构建工具
```

## 开发指南

### 创建新服务

```bash
# 创建 proto 文件
kratos proto add api/your_service/your_service.proto

# 生成代码
make all
```

### 分层架构

项目采用分层架构，各层职责如下：

| 层级          | 目录                 | 职责                             |
|-------------|--------------------|--------------------------------|
| **Service** | `internal/service` | 实现 protobuf 服务接口，调用 Biz 层      |
| **Biz**     | `internal/biz`     | 核心业务逻辑，定义实体和接口                 |
| **Data**    | `internal/data`    | 数据访问，实现 Biz 层定义的 Repository 接口 |
| **Client**  | `internal/client`  | 外部服务调用，实现 Biz 层定义的 Client 接口   |
| **Conf**    | `internal/conf`    | 配置定义与加载                        |

### 服务注册

在 `cmd/server/bootstrap.go` 中注册 HTTP/gRPC/WebSocket 服务：

```go
package main

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/service"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/server/websocket"
)

func NewBootstrap(
	httpServer *http.Server,
	grpcServer *grpc.Server,
	wss *websocket.Server,
	exampleService *service.ExampleService,
	exampleWsHandler *service.ExampleWsHandler,
) bootstrap.Bootstrap {
	// WebSocket 路由
	wss.Handle("/echo", exampleWsHandler, websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	// 绑定 HTTP 服务
	example_pb.RegisterExampleServiceHTTPServer(httpServer, exampleService)

	// 绑定 gRPC 服务
	example_pb.RegisterExampleServiceServer(grpcServer, exampleService)

	return nil
}

```

### WebSocket 处理器

实现 WebSocket 处理器接口：

```go
package service

import (
	"net/http"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/server/websocket"
)

type ExampleWsHandler struct{}

func NewExampleWsHandler() *ExampleWsHandler {
	return &ExampleWsHandler{}
}

func (h *ExampleWsHandler) OnHandshake(request *http.Request) error {
	// 连接前校验
	return nil
}

func (h *ExampleWsHandler) OnConnect(client *websocket.Client) {
	// 连接建立后回调
}

func (h *ExampleWsHandler) OnMessage(client *websocket.Client, message []byte, messageType websocket.MessageType) {
	// 收到消息后处理
	client.SendText("echo: " + string(message))
}

func (h *ExampleWsHandler) OnError(client *websocket.Client, err error) {
	// 错误处理
}

func (h *ExampleWsHandler) OnClose(client *websocket.Client) {
	// 连接关闭后回调
}

```

### Job 定时任务

框架内置了基于 Cron 的定时任务支持，集成 OpenTelemetry 可观测性。

#### 配置定义

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

**调度表达式**：

| 表达式           | 说明                 |
|---------------|--------------------|
| `@every 5s`   | 每 5 秒执行一次          |
| `@every 1m`   | 每 1 分钟执行一次         |
| `@every 1h`   | 每 1 小时执行一次         |
| `@hourly`     | 每小时整点执行            |
| `@daily`      | 每天午夜执行             |
| `@weekly`     | 每周日午夜执行            |
| `@monthly`    | 每月 1 日午夜执行         |
| `0 2 * * *`   | 标准 Cron：每天 2:00 执行 |
| `*/5 * * * *` | 标准 Cron：每 5 分钟执行   |

**并发策略**：

| 策略        | 说明                    |
|-----------|-----------------------|
| `DELAY`   | 上一次任务未完成时，等待完成后再执行下一次 |
| `SKIP`    | 上一次任务未完成时，跳过本次执行      |
| `OVERLAP` | 允许多个实例并发执行            |

#### 实现 Job 接口

创建 Job 实现文件 `internal/job/sync_data.go`：

```go
package job

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
)

type SyncDataJob struct {
	log *log.Helper
	db  *database.Manager
}

func NewSyncDataJob(logger log.Logger, db *database.Manager) *SyncDataJob {
	return &SyncDataJob{
		log: log.NewHelper(logger),
		db:  db,
	}
}

// Run 实现 job.Job 接口
func (j *SyncDataJob) Run(ctx context.Context) error {
	j.log.WithContext(ctx).Info("开始执行数据同步任务")
	// 业务逻辑...
	j.log.WithContext(ctx).Info("数据同步任务完成")
	return nil
}
```

#### 注册 Job

在 `cmd/server/bootstrap.go` 中注册：

```go
package main

import (
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/job"
)

func NewBootstrap(
// ...其他参数
	register *job.Register,
	syncDataJob *myjob.SyncDataJob,
) bootstrap.Bootstrap {
	// 注册 Job（名称需与配置文件中的 key 对应）
	register.Register("sync_data", syncDataJob)
	// ...
	return nil
}

```

#### Wire 依赖注入

`internal/job/wire.go`：

```go
package job

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSyncDataJob,
)

```

#### 函数式 Job（快捷方式）

```go
package main

import (
	"fmt"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/bootstrap"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/job"
	job2 "github.com/jaggerzhuang1994/kratos-foundation/pkg/component/job/job"
)

func NewBootstrap(
// ...其他参数
	register *job.Register,
) bootstrap.Bootstrap {
	register.Register("simple_task", job2.FuncJob(func(ctx context.Context) error {
		fmt.Println("执行简单任务")
		return nil
	}))
	// ...
	return nil
}

```

#### Job 可观测性

框架自动集成 OpenTelemetry：

- **Tracing**：每次执行创建独立 Span
- **Metrics**：记录执行次数、耗时、成功/失败率
- **Logging**：日志自动携带 TraceID

### 客户端注入

protobuf 生成的 `_client.pb.go` 提供自动生成的 Wire Provider：

```go
package client

import (
	"github.com/google/wire"
	"your-project/api/example_service/example_pb"
)

type BizImpl struct {
	api example_pb.ExampleServiceApi
}

func NewBizImpl(api example_pb.ExampleServiceApi) *BizImpl {
	return &BizImpl{api}
}

var ProviderSet = wire.NewSet(
	example_pb.ExampleServiceApiProvider,
	NewBizImpl,
)
```

### 依赖倒置原则

**Biz 层定义接口，Data/Client 层实现接口**

```
Service → Biz ← Data
            ↑
            └── Client
```

#### 示例：用户查询

**Biz 层定义** (`internal/biz/user.go`)：

```go
package biz

import "context"

type User struct {
	ID   int64
	Name string
}

// 数据访问接口（由 Data 层实现）
type GetUserRepo interface {
	GetUser(ctx context.Context, id int64) (User, error)
}

type GetUserBiz struct {
	repo GetUserRepo
}

func NewGetUserBiz(repo GetUserRepo) *GetUserBiz {
	return &GetUserBiz{repo}
}

func (b *GetUserBiz) GetUser(ctx context.Context, id int64) (User, error) {
	return b.repo.GetUser(ctx, id)
}
```

**Data 层实现** (`internal/data/user_repo.go`)：

```go
package data

import (
	"context"
	"your-project/internal/biz"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
)

type UserTableRepo struct {
	db *database.Manager
}

func NewUserTableRepo(db *database.Manager) *UserTableRepo {
	return &UserTableRepo{db}
}

func (r *UserTableRepo) GetUser(ctx context.Context, id int64) (biz.User, error) {
	// 查询数据库
	return biz.User{}, nil
}
```

**Wire 绑定** (`internal/data/wire.go`)：

```go
package data

import (
	"github.com/google/wire"
	"your-project/internal/biz"
)

var _ biz.GetUserRepo = (*UserTableRepo)(nil)

var ProviderSet = wire.NewSet(
	NewUserTableRepo,
	wire.Bind(new(biz.GetUserRepo), new(*UserTableRepo)),
)
```

## Docker 部署

### 本地开发

```bash
# 复制环境变量模板
cp docker/.env.example docker/.env

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 构建镜像

```bash
# 构建指定目标
docker build --build-arg TARGET=server -t your-app:latest .

# 或使用 Makefile
make build-server
```

### 环境配置

**配置优先级**：

- 本地环境：远程配置 > 本地文件
- 其他环境：本地文件 > 远程配置

**配置文件加载顺序**：

```
configs/config.yaml
configs/{app_name}.yaml
configs/{env}.config.yaml
configs/{env}.{app_name}.yaml
```

## 常用命令

| 命令                   | 说明         |
|----------------------|------------|
| `make init`          | 安装依赖工具     |
| `make all`           | 生成所有代码     |
| `make api`           | 生成 API 代码  |
| `make config`        | 生成配置代码     |
| `make generate`      | 生成 Wire 代码 |
| `make run`           | 运行服务       |
| `make build`         | 构建所有二进制    |
| `make build-server`  | 构建指定目标     |
| `make lint`          | 代码检查       |
| `make test`          | 运行测试       |
| `make test-coverage` | 生成覆盖率报告    |
| `make help`          | 查看帮助       |

## 相关资源

- [Go-Kratos 官方文档](https://go-kratos.dev/)
- [Kratos Foundation](https://github.com/jaggerzhuang1994/kratos-foundation)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Wire 依赖注入](https://github.com/google/wire)

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。
