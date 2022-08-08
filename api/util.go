package api

import (
	"encoding/json"
	"errors"
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/grafana"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		return generateGrafanaTextCardMessageFromContext(context, alias)
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
		logger.Errorf("generate text card message failed, error: %v", err)
		return nil, err
	}
	return data, err
}

func generateGrafanaTextCardMessageFromContext(context *gin.Context, alias string) (data []byte, err error) {
	bytes, _ := ioutil.ReadAll(context.Request.Body)
	logger.Errorf("request body, body: %s", string(bytes))
	var grafanaMsg grafana.GrafanaMessage
	err = json.Unmarshal(bytes, &grafanaMsg)
	if err != nil {
		logger.Errorf("unmarshal grafana message failed, error: %+v, body: %s", err, string(bytes))
		return nil, err
	}

	messages := strings.Split(grafanaMsg.Message, "\n")
	description := ""
	cardUrl := "http://grafana.sys.ink:8080"

	for _, m := range messages {
		if strings.HasPrefix(m, " - summary =") {
			description = m
		}
		if strings.HasPrefix(m, "Source =") {
			cardUrl = m
		}
	}

	msg := wecom.TextCardMessage{
		Touser:  config.Config.WeComConfigs[alias].Receiver,
		Msgtype: "textcard",
		Agentid: config.Config.WeComConfigs[alias].AgentId,
	}

	msg.TextCard.Title = grafanaMsg.Title
	msg.TextCard.Description = description
	msg.TextCard.URL = cardUrl

	data, err = json.Marshal(msg)
	if err != nil {
		logger.Errorf("generate grafana text card message failed, error: %v", err)
		return nil, err
	}

	return data, err
}
