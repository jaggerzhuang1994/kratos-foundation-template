# Kratos Foundation Template

基于 [Go-Kratos](https://go-kratos.dev/) 微服务框架的项目模板，提供了开箱即用的项目结构和最佳实践。

## 特性

- 🚀 基于 Kratos v2 微服务框架
- 📦 清晰的分层架构（Service/Biz/Data/Client）
- 🔌 Wire 依赖注入
- 🛠️ 完善的 Makefile 工具链
- 📝 Protocol Buffers 驱动开发
- 🔧 自动生成客户端代码

## 前置要求

- Go >= 1.22
- Protocol Buffers 编译器
- Kratos CLI 工具

## 快速开始

### 1. 安装 Kratos CLI

```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### 2. 使用模板创建项目

```bash
# 创建项目
kratos new example -r https://github.com/jaggerzhuang1994/kratos-foundation-template

# 进入项目目录
cd example

# 安装依赖
make init

# 生成代码
make all

# 运行服务
make run
```

### 3. 重命名模块名（可选）

如果项目有独立的 Git 仓库，建议重命名 Go 模块名：

```bash
make rename-module NEW=github.com/your-org/your-project
```

## 项目结构

```
.
├── api/                    # API 定义（protobuf）
│   └── server/
├── cmd/                    # 应用入口
│   └── app/
├── internal/              # 内部代码
│   ├── bootstrap/         # 应用启动与初始化
│   ├── biz/               # 业务逻辑层
│   ├── client/            # 外部服务客户端
│   ├── conf/              # 配置定义
│   ├── data/              # 数据访问层
│   └── service/           # 服务实现层
├── configs/               # 配置文件
└── Makefile              # 构建工具
```

## 开发指南

### 创建新服务

```bash
# 创建 proto 文件
kratos proto add api/server/server.proto

# 生成代码
make all
```

### 编写业务代码

项目采用分层架构，开发流程如下：

1. **Service 层** (`internal/service`)
   - 实现 protobuf 生成的服务接口
   - 调用 Biz 层业务逻辑
   - 在 `internal/bootstrap/bootstrap.go` 中注册服务

2. **Biz 层** (`internal/biz`)
   - 编写核心业务逻辑
   - 定义业务实体和错误
   - 定义外部依赖接口（Repository/Client）

3. **Data 层** (`internal/data`)
   - 实现 Biz 层定义的数据接口
   - 处理数据库访问、缓存等

4. **Client 层** (`internal/client`)
   - 实现 Biz 层定义的客户端接口
   - 处理外部 API 调用

5. **配置管理** (`internal/conf`)
   - 在 proto 中定义配置结构
   - 执行 `make config` 生成配置代码
   - 通过注入 `*conf.Bootstrap` 访问配置

### 客户端注入

生成的 protobuf 代码会包含 `_client.pb.go` 文件，提供自动生成的 Wire Provider：

```go
package client

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/server"
)

type BizImpl struct {
	api server.ServerApi  // 使用生成的 ServerApi 接口
}

func NewBizImpl(api server.ServerApi) *BizImpl {
	return &BizImpl{api}
}

var ProviderSet = wire.NewSet(
	server.ServerApiProvider,  // 使用自动生成的 Provider
	NewBizImpl,
)
```

### 依赖倒置原则

**Biz 层定义接口，外部层实现接口**

#### 示例：用户查询功能

**1. Biz 层定义** (`internal/biz/user.go`)

```go
package biz

import (
	"context"
	"github.com/jaggerzhuang1994/kratos-foundation/proto/kratos_foundation_pb"
)

// 业务错误定义
var ErrUserNotFound = kratos_foundation_pb.ErrorNotFound("user not found")

// 业务实体
type User struct {
	ID   int64
	Name string
}

// 数据访问接口（由 Data 层实现）
type GetUserRepo interface {
	GetUser(ctx context.Context, id int64) (User, error)
}

// 业务逻辑
type GetUserBiz struct {
	getUserRepo GetUserRepo
}

func NewGetUserBiz(getUserRepo GetUserRepo) *GetUserBiz {
	return &GetUserBiz{getUserRepo}
}

func (biz *GetUserBiz) GetUser(ctx context.Context, id int64) (User, error) {
	return biz.getUserRepo.GetUser(ctx, id)
}
```

**2. Data 层实现** (`internal/data/user_table_repo.go`)

```go
package data

import (
	"context"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
)

type UserTableRepo struct {
	db *database.Manager
}

func NewUserTableRepo(db *database.Manager) *UserTableRepo {
	return &UserTableRepo{db}
}

// 实现 biz.GetUserRepo 接口
func (repo *UserTableRepo) GetUser(ctx context.Context, id int64) (biz.User, error) {
	// 1. 使用 data/po 或 data/model 中定义的 GORM 模型
	// 2. 查询数据库
	// 3. 转换为 biz.User 实体
	// 4. 处理错误（如记录不存在，返回 biz.ErrUserNotFound）

	// 示例代码（伪代码）
	// var userPO UserPO
	// if err := repo.db.First(&userPO, id).Error; err != nil {
	//     if errors.Is(err, gorm.ErrRecordNotFound) {
	//         return biz.User{}, biz.ErrUserNotFound
	//     }
	//     return biz.User{}, err
	// }
	// return biz.User{ID: userPO.ID, Name: userPO.Name}, nil

	return biz.User{}, nil
}
```

**3. Wire 绑定** (`internal/data/wire.go`)

```go
package data

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz"
)

// 编译时接口约束检查
var _ biz.GetUserRepo = (*UserTableRepo)(nil)

var ProviderSet = wire.NewSet(
	NewUserTableRepo,
	wire.Bind(new(biz.GetUserRepo), new(*UserTableRepo)),  // 绑定接口实现
)
```

## 设计原则

### 分层职责

- **Service**: 协议适配，参数校验，调用 Biz
- **Biz**: 业务逻辑，定义接口，不依赖具体实现
- **Data**: 数据访问，实现 Biz 定义的接口
- **Client**: 外部调用，实现 Biz 定义的接口

### 依赖方向

```
Service → Biz ← Data
            ↑
            └── Client
```

### 错误处理

- Biz 层定义业务错误
- Data/Client 层将底层错误转换为业务错误
- 使用 kratos-foundation 提供的标准错误类型

### 实体映射

- **Biz Entity**: 业务实体，贫血模型
- **Data PO/Model**: 数据库模型，与 ORM 绑定
- **Service DTO**: API 请求/响应，由 protobuf 生成

不同层次使用不同的数据结构，避免跨层污染。

## 常用命令

```bash
# 初始化依赖
make init

# 生成所有代码
make all

# 仅生成 API 代码
make api

# 仅生成配置代码
make config

# 生成 Wire 依赖注入代码
make generate

# 运行服务
make run

# 构建二进制
make build

# 代码检查
make lint

# 查看帮助
make help
```

## 相关资源

- [Go-Kratos 官方文档](https://go-kratos.dev/)
- [Kratos Foundation](https://github.com/jaggerzhuang1994/kratos-foundation)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Wire 依赖注入](https://github.com/google/wire)

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。
