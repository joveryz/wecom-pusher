package main

import (
	"fmt"

	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
)

func main() {
	fmt.Println("helloworld")
	wecom.SendTextMessage("test", "downloadpush")
	wecom.SendTextCardMessage("test cccc", "test", "https://www.baidu.com", "downloadpush")
}

func init() {
	jsonConfig := "D:\\config.json"
	logger.Infof("Loading config from %s", jsonConfig)
	config.LoadConfig(jsonConfig)
}
