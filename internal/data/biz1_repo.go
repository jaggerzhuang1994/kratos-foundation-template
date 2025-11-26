package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/biz/biz1"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/database"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/component/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/utils"
)

type biz1Repo struct {
	log *log.Helper
	db  *database.Manager
}

// NewBiz1Repo .
func NewBiz1Repo(log *log.Log, db *database.Manager) biz1.Repo {
	return &biz1Repo{
		log.WithModule("data/biz1_repo").NewHelper(),
		db,
	}
}

func (repo *biz1Repo) Get() biz1.Obj {
	var rows []struct {
		Database string
	}
	err := repo.db.GetConnection(context.Background()).Raw("show databases;").Scan(&rows).Error
	fmt.Println("err", err)
	return biz1.Obj(strings.Join(utils.Map(rows, func(t struct{ Database string }) string {
		return t.Database
	}), " "))
}
