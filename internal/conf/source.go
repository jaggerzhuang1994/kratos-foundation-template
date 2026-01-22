package conf

import (
	"fmt"
	"os"
	"path/filepath"

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
	if err != nil {
		return nil, err
	}
	// 如果fileConfigSource是目录，则读取目录下
	//   config.yaml
	//   {appInfo.Name}.yaml
	//   {env}.config.yaml
	//   {env}.{appInfo.Name}.yaml
	if dir {
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
	// 如果fileConfigSource不是目录，则当作文件返回
	return []string{
		fileConfigSourcePath,
	}, nil
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
