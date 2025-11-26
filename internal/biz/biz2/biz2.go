package biz2

import "context"

type Obj string

type Repo interface {
	Set(ctx context.Context, k string, v Obj) error
}

type Repo2 interface {
	Get(ctx context.Context, k string) (Obj, bool)
}

type Biz struct {
	repo  Repo
	repo2 Repo2
}

func NewBiz1(repo Repo, repo2 Repo2) *Biz {
	return &Biz{repo, repo2}
}

func (biz *Biz) Set(ctx context.Context, obj Obj) error {
	return biz.repo.Set(ctx, "a", obj)
}

func (biz *Biz) Get(ctx context.Context) (Obj, bool) {
	return biz.repo2.Get(ctx, "a")
}
