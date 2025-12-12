package data

import (
	"context"
	"fmt"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/example"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
)

// 实现业务需要的接口

type UserDbRepo struct {
	db *database.Manager
}

func NewUserDbRepo(db *database.Manager) *UserDbRepo {
	return &UserDbRepo{db: db}
}

func (r *UserDbRepo) GetUser(ctx context.Context, id int64) (*example.User, error) {
	if id == 0 {
		return nil, example.ErrUserNotFound
	}
	// 模拟调用db
	_ = r.db.GetConnection(ctx).Exec("show tables")
	return &example.User{
		ID:   id,
		Name: fmt.Sprintf("user%d", id),
	}, nil
}
