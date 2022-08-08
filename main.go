package main

import (
	"flag"

	"github.com/TongboZhang/wecom-pusher/api"
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
)

func main() {
	err := api.Start()
	logger.Error(err)
}

func init() {
	jsonConfig := flag.String("c", "config.json", "config path")
	flag.Parse()
	logger.Infof("Loading config from %s", *jsonConfig)
	config.LoadConfig(*jsonConfig)
}
