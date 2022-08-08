package api

import (
	"encoding/json"
	"errors"
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
	"github.com/gin-gonic/gin"
	"strings"
)

func parseFromContext(context *gin.Context, key string) (value string) {
	value = context.GetHeader(key)
	if value == "" {
		value = context.Query(key)
	}
	if value == "" {
		value = context.PostForm(key)
	}

	return value
}

func parseDestinationFromContext(context *gin.Context) (value []string) {
	destinationString := parseFromContext(context, "destination")
	var destination []string
	if destinationString != "" {
		destination = strings.Split(destinationString, "|")
	} else {
		destination = config.Config.Aliases
	}

	return destination
}

func generateWeComMessageFromContext(context *gin.Context, alias string) (data []byte, err error) {
	msgType := parseFromContext(context, "type")
	if msgType == "text" {
		return generateWeComTextMessageFromContext(context, alias)
	} else if msgType == "textcard" {
		return generateWeComTextCardMessageFromContext(context, alias)
	} else if msgType == "grafana" {
		return nil, err
	} else {
		return nil, errors.New("type not supported, type: " + msgType)
	}
}

func generateWeComTextMessageFromContext(context *gin.Context, alias string) (data []byte, err error) {
	content := parseFromContext(context, "content")
	msg := wecom.TextMessage{
		Touser:  config.Config.WeComConfigs[alias].Receiver,
		Msgtype: "text",
		Agentid: config.Config.WeComConfigs[alias].AgentId,
	}
	msg.Text.Content = content

	data, err = json.Marshal(msg)
	if err != nil {
		logger.Errorf("generate text message failed, error: %v", err)
		return nil, err
	}
	return data, err
}

func generateWeComTextCardMessageFromContext(context *gin.Context, alias string) (data []byte, err error) {
	title := parseFromContext(context, "title")
	description := parseFromContext(context, "description")
	cardUrl := parseFromContext(context, "cardUrl")
	msg := wecom.TextCardMessage{
		Touser:  config.Config.WeComConfigs[alias].Receiver,
		Msgtype: "textcard",
		Agentid: config.Config.WeComConfigs[alias].AgentId,
	}

	msg.TextCard.Title = title
	msg.TextCard.Description = description
	msg.TextCard.URL = cardUrl

	data, err = json.Marshal(msg)
	if err != nil {
		logger.Errorf("generate text message failed, error: %v", err)
		return nil, err
	}
	return data, err
}
