package main

import (
	"github.com/TongboZhang/wecom-pusher/api"
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
)

func main() {
	// wecom.SendTextMessage("test", "downloadpush")
	// wecom.SendTextCardMessage("test cccc", "test", "https://www.baidu.com", "downloadpush")
	err := api.Start()
	logger.Error(err)
}

func init() {
	jsonConfig := "D:\\config.json"
	logger.Infof("Loading config from %s", jsonConfig)
	config.LoadConfig(jsonConfig)
}
