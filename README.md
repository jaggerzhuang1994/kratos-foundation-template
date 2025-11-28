# kratos 项目模板

## 依赖

* kratos

## 安装 kratos

```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## 使用模板创建项目

```
> kratos new example -r https://github.com/jaggerzhuang1994/kratos-foundation-template
> cd example
> make all
```

### 重命名go模块名

```
> make rename-module NEW=新的模块名
```

### 创建服务

```
# 创建pb
> kratos proto add api/server/server.proto
> make all
```

### 编写业务代码

* 在 internal/service 编写服务入口，需实现 server.UnimplementedServerServer 接口
* 在 internal/bootstrap/bootstrap.go 中向 httpServer/grpcServer 注册服务 server.RegisterServerServer
* 在 internal/biz 中编写业务代码，并在 internal/service 中注入并调用 biz 业务逻辑
* 在 internal/data 中实现 internal/biz 需要的数据接口
* 在 internal/client 中实现 internal/biz 需要的客户端接口
* 在 internal/conf 中添加业务需要的配置，并执行 make config 生成pb，在业务代码中注入 *conf.Bootstrap 即可访问

### 客户端注入

在生成的pb中会存在一个 _client.pb.go 的文件，会自动生成 wire 注入器 `ServerApiProvider`，客户端实例注入使用 `ServerApi`

```go
package client

import (
	"github.com/google/wire"
	"github.com/jaggerzhuang1994/kratos-foundation-template/api/server"
)

type BizImpl struct {
	api server.ServerApi
}

func NewBizImpl(api server.ServerApi) *BizImpl {
	return &BizImpl{api}
}

// ProviderSet 在 client 的wire中添加以下 provider 即可
var ProviderSet = wire.NewSet(
	server.ServerApiProvider,
	NewBizImpl,
)

```

### biz业务接口定义以及外部实现的规范

在 biz 中如果遇到外部逻辑（外部api/数据库/缓存/队列等）都应该在 biz 包内定义逻辑的
interface，然后在对应包中实现并用wire绑定，biz注入该逻辑接口即可

#### 例如

biz/user.go 实现获取用户信息的业务逻辑`GetUserBiz`，定义获取用户数据的逻辑接口 `GetUserRepo`

```go
package biz

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation/proto/kratos_foundation_pb"
)

// 在 biz 包内定义业务错误
// 这里不只可以用我们自己生成的 errors，也可以用foundation包内的 errors 的错误（这些错误会自动渲染为对应的http响应）
var ErrUserNotFound = kratos_foundation_pb.ErrorNotFound("user not found")

// 在 biz 包内定义业务实体
type User struct {
	ID   int64
	Name string
}

// 在biz中定义外部获取的逻辑接口
type GetUserRepo interface {
	GetUser(ctx context.Context, id int64) (User, error)
}

// 在业务实体上注入
type GetUserBiz struct {
	getUserRepo GetUserRepo
}

// 构造函数需要注册在biz.ProviderSet内
func NewGetUserBiz(getUserRepo GetUserRepo) *GetUserBiz {
	return &GetUserBiz{getUserRepo}
}

// 定义方法给 service 调用
func (biz *GetUserBiz) GetUser(ctx context.Context, id int64) (User, error) {
	return biz.getUserRepo.GetUser(ctx, id)
}

```

在 data 中需要实现这个获取数据的逻辑接口

data/user_table_repo.go

```go
package data

import (
	"context"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
)

type UserTableRepo struct {
	db *database.Manager
}

// 需要在 data.ProviderSet 中注册
func NewUserTableRepo(db *database.Manager) *UserTableRepo {
	return &UserTableRepo{db}
}

// 需要实现 biz 中的业务逻辑
func (repo *UserTableRepo) GetUser(ctx context.Context, id int64) (biz.User, error) {
	// 这里需要使用 model 获取用户，并转成 biz.User 返回
	// 这里不能直接用 biz.User 作为数据库读写的model
	// 应该在 data/po 或者 data/model 下定义对应的模型结构用作 gorm 的模型
	// 不能在 biz 包内断言 err 是不是 gorm 的 RecordNotFound 错误，需要在这里判断 db 的数据是否存在，不存在 需要返回 biz.ErrUserNotFound 错误
}

```

实现完读数据接口后，向 data.ProviderSet 注册实现和绑定接口实现，并且约束接口实现。

data/wire.go

```go
package data

import "github.com/google/wire"

// repo 实现哪些接口在这里约束 wire.Bind 只会在 wire 时提示找不到接口或者接口不满足约束，这里直接写可以实时提示是否满足接口约束
var _ biz.GetUserRepo = (*UserTableRepo)(nil) // 约束接口

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewUserTableRepo,                                     // 注册实现
	wire.Bind(new(biz.GetUserRepo), new(*UserTableRepo)), // 绑定接口
)

```
