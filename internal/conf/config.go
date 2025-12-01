package conf

import (
	"fmt"

	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/config"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/consul"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/env"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/utils"
)

type LocalFilePath []string

func NewConfig(
	appInfo *app_info.AppInfo,
	consul *consul.Client,
	localFilePath LocalFilePath,
) (config.Config, func(), error) {
	// 本地配置
	localConfigSource, err := config.NewFilePatternSource(localFilePath)
	if err != nil {
		return nil, nil, err
	}

	appEnv := env.AppEnv()
	appName := appInfo.GetName()
	var consulConfigList = []string{
		// 不指定环境配置
		"app/configs/common.",
		"app/secrets/common.",
		fmt.Sprintf("app/configs/%s.common.", appEnv),
		fmt.Sprintf("app/secrets/%s.common.", appEnv),
		fmt.Sprintf("app/configs/%s.", appName),
		fmt.Sprintf("app/secrets/%s.", appName),
		fmt.Sprintf("app/configs/%s.%s.", appEnv, appName),
		fmt.Sprintf("app/secrets/%s.%s.", appEnv, appName),
	}
	// 远程配置
	remoteConfigSource := config.NewConsulSource(consul, consulConfigList)

	var sources []config.Source
	// 如果是本地环境，本地source优先级更高
	if env.IsLocal() {
		sources = append(sources, remoteConfigSource)
		sources = append(sources, localConfigSource)
	} else {
		// 其他环境远程配置优先级更高
		sources = append(sources, localConfigSource)
		sources = append(sources, remoteConfigSource)
	}

	sources = utils.FilterZero(sources)
	c := config.New(config.WithSource(config.NewPriorityConfigSource(sources)))

	if err := c.Load(); err != nil {
		_ = c.Close() // release config watcher
		return nil, nil, err
	}

	return c, func() {
		_ = c.Close()
	}, nil
}
