package main

import (
	"fmt"

	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
)

var (
	BuildTime string
	GoVersion string
	GitHead   string
)

func main() {
	fmt.Println("helloworld")
	logger.Error()
}

func init() {
	jsonConfig := "D:\\config.json"
	logger.Infof("Loading config from %s", jsonConfig)
	config.LoadConfig(jsonConfig)
}
