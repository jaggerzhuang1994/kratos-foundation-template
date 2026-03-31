package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/config"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/env"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/utils"
)

// FileConfigSource main需要提供这个类型
type FileConfigSource string

// NewFileSource 文件配置源
func NewFileSource(
	appInfo app_info.AppInfo,
	fileConfigSource FileConfigSource,
) (config.FileSourcePathList, error) {
	fileConfigSourcePath := string(fileConfigSource)
	if fileConfigSourcePath == "" {
		return nil, nil
	}

	dir, err := isDir(fileConfigSourcePath)
	// 如果传入不存在的文件路径，则告警，不要中断运行
	if err != nil {
		log.Warnf("failed to stat config path: path=%s, error=%v", fileConfigSourcePath, err)
		return nil, nil
	}

	// 如果 fileConfigSource 是文件路径，则直接返回
	if !dir {
		return []string{
			fileConfigSourcePath,
		}, nil
	}

	// 如果 fileConfigSource 是目录，则读取目录下存在的文件
	//   config.yaml
	//   {appInfo.Name}.yaml
	//   {env}.config.yaml
	//   {env}.{appInfo.Name}.yaml
	fileList := []string{
		filepath.Join(fileConfigSourcePath, "config.yaml"),
		filepath.Join(fileConfigSourcePath, fmt.Sprintf("%s.yaml", appInfo.GetName())),
		filepath.Join(fileConfigSourcePath, fmt.Sprintf("%s.config.yaml", env.AppEnv())),
		filepath.Join(fileConfigSourcePath, fmt.Sprintf("%s.%s.yaml", env.AppEnv(), appInfo.GetName())),
	}
	fileList = utils.Filter(fileList, fileExists)
	if len(fileList) == 0 {
		return nil, nil
	}
	return fileList, nil
}

// NewConsulSource consul 的配置源
func NewConsulSource(appInfo app_info.AppInfo) config.ConsulSourcePathList {
	appEnv := env.AppEnv()
	appName := appInfo.GetName()
	return []string{
		// 不指定环境配置
		"app/configs/common.",
		"app/secrets/common.",
		fmt.Sprintf("app/configs/%s.", appName),
		fmt.Sprintf("app/secrets/%s.", appName),
		// 指定环境
		fmt.Sprintf("app/configs/%s.common.", appEnv),
		fmt.Sprintf("app/secrets/%s.common.", appEnv),
		fmt.Sprintf("app/configs/%s.%s.", appEnv, appName),
		fmt.Sprintf("app/secrets/%s.%s.", appEnv, appName),
	}
}

func isDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
