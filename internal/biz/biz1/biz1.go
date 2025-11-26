package biz1

type Obj string

type Repo interface {
	Get() Obj
}

type Biz struct {
	repo Repo
}

func NewBiz1(repo Repo) *Biz {
	return &Biz{repo}
}

func (b *Biz) Get() Obj {
	return b.repo.Get()
}
