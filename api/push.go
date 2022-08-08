package api

import (
	"fmt"
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
	"github.com/gin-gonic/gin"
	"strings"
)

func BasicPush(context *gin.Context) {
	content := parseFromContext(context, "content")
	destinationString := parseFromContext(context, "destination")
	var destination []string
	if destinationString != "" {
		destination = strings.Split(destinationString, "|")
	} else {
		destination = config.Config.Aliases
	}

	var errString []string
	for _, alias := range destination {
		err := wecom.SendTextMessage(content, alias)
		if err != nil {
			errString = append(errString, fmt.Sprintf("basic push to wecom failed, content: %s, alias: %s, error: %+v", content, alias, err))
		}
	}

	if len(errString) != 0 {
		context.String(500, strings.Join(errString, "|"))
		logger.Error(strings.Join(errString, "|"))
		return
	}

	context.String(200, "ok")
}
