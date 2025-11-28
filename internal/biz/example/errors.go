package example

import "github.com/jaggerzhuang1994/kratos-foundation-template/api/example_service/example_pb"

// 业务异常定义在 biz，然后在 repo 包或者其他依赖 biz 的包中调用

var ErrUserNotFound = example_pb.ErrorUserNotFound("user not found")
