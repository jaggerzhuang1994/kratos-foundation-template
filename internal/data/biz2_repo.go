package data

import (
	"context"
	"errors"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz2"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/redis"
)

type Biz2Repo struct {
	log   *log.Helper
	redis *redis.Manager
}

var _ biz2.Repo = (*Biz2Repo)(nil)
var _ biz2.Repo2 = (*Biz2Repo)(nil)

// NewBiz2Repo .
func NewBiz2Repo(log *log.Log, redis *redis.Manager) *Biz2Repo {
	return &Biz2Repo{
		log.WithModule("data/biz2_repo").NewHelper(),
		redis,
	}
}

func (r *Biz2Repo) Get(ctx context.Context, k string) (biz2.Obj, bool) {
	cmd := r.redis.Get(ctx, k)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return "", false
		} else {
			return "", false
		}
	}
	return biz2.Obj(cmd.Val()), true
}

func (r *Biz2Repo) Set(ctx context.Context, k string, v biz2.Obj) error {
	cmd := r.redis.Set(ctx, k, string(v), redis.KeepTTL)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
