package api

import (
	"fmt"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
	"github.com/gin-gonic/gin"
	"strings"
)

func Push(context *gin.Context) {
	destination := parseDestinationFromContext(context)
	var errStrings []string
	for _, alias := range destination {
		data, err := generateWeComMessageFromContext(context, alias)
		if err != nil {
			errStrings = append(errStrings, fmt.Sprintf("basic push to generate wecom msg failed, data: %s, alias: %s, error: %+v", string(data), alias, err))
			continue
		}
		err = wecom.SendTextMessage(data, alias)
		if err != nil {
			errStrings = append(errStrings, fmt.Sprintf("basic push to wecom failed, data: %s, alias: %s, error: %+v", string(data), alias, err))
		}
	}

	if len(errStrings) != 0 {
		context.JSON(500, errStrings)
		logger.Error(strings.Join(errStrings, "|"))
		return
	}

	context.String(200, "ok")
}
