package main

import (
	"flag"

	"github.com/jaggerzhuang1994/kratos-foundation-template/internal/conf"
	"github.com/jaggerzhuang1994/kratos-foundation/pkg/app_info"
	_ "github.com/jaggerzhuang1994/kratos-foundation/pkg/setup"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// appInfo
	appInfo := app_info.NewAppInfo(Version)
	app_info.PrintAppInfo(appInfo)

	// wireApp
	app, cleanup, err := wireApp(appInfo, conf.FileConfigSource(flagconf))
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
