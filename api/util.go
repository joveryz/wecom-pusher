package api

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/grafana"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/TongboZhang/wecom-pusher/wecom"
	"github.com/gin-gonic/gin"
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

func generateAndPushWeComMessageFromContext(context *gin.Context, aliases []string) (isSucc bool) {
	msgType := parseFromContext(context, "type")
	if msgType == "text" {
		return generateAndPushWeComTextMessageFromContext(context, aliases)
	} else if msgType == "textcard" {
		return generateAndPushWeComTextCardMessageFromContext(context, aliases)
	} else if msgType == "grafana" {
		return generateAndPushGrafanaTextCardMessageFromContext(context, aliases)
	} else {
		return false
	}
}

func generateAndPushWeComTextMessageFromContext(context *gin.Context, aliases []string) (isSucc bool) {
	content := parseFromContext(context, "content")
	isSucc = true
	for _, alias := range aliases {
		msg := wecom.TextMessage{
			Touser:  config.Config.WeComConfigs[alias].Receiver,
			Msgtype: "text",
			Agentid: config.Config.WeComConfigs[alias].AgentId,
		}
		msg.Text.Content = content
		data, err := json.Marshal(msg)
		if err != nil {
			isSucc = false
			logger.Errorf("generate text message failed, skip pushing, error: %v", err)
			continue
		}
		err = wecom.SendTextMessage(data, alias)
		if err != nil {
			isSucc = false
			logger.Errorf("push text message failed, error: %v", err)
			continue
		}
	}
	return isSucc
}

func generateAndPushWeComTextCardMessageFromContext(context *gin.Context, aliases []string) (isSucc bool) {
	title := parseFromContext(context, "title")
	description := parseFromContext(context, "description")
	cardUrl := parseFromContext(context, "cardUrl")
	isSucc = true
	for _, alias := range aliases {
		msg := wecom.TextCardMessage{
			Touser:  config.Config.WeComConfigs[alias].Receiver,
			Msgtype: "textcard",
			Agentid: config.Config.WeComConfigs[alias].AgentId,
		}

		msg.TextCard.Title = title
		msg.TextCard.Description = description
		msg.TextCard.URL = cardUrl

		data, err := json.Marshal(msg)
		if err != nil {
			isSucc = false
			logger.Errorf("generate text message failed, skip pushing, error: %v", err)
			continue
		}
		err = wecom.SendTextMessage(data, alias)
		if err != nil {
			isSucc = false
			logger.Errorf("push text message failed, error: %v", err)
			continue
		}
	}
	return isSucc
}

func generateAndPushGrafanaTextCardMessageFromContext(context *gin.Context, aliases []string) (isSucc bool) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var grafanaMsg grafana.GrafanaMessage
	err := json.Unmarshal(body, &grafanaMsg)
	if err != nil {
		logger.Errorf("unmarshal grafana request body failed, error: %+v, body: %s", err, string(body))
		return false
	}

	messages := strings.Split(grafanaMsg.Message, "\n")
	description := ""
	cardUrl := "http://grafana.sys.ink:8080"
	isSucc = true
	for _, m := range messages {
		if strings.HasPrefix(m, " - summary =") {
			description = m
		}
		if strings.HasPrefix(m, "Source: ") {
			cardUrl = m
		}
	}

	for _, alias := range aliases {
		msg := wecom.TextCardMessage{
			Touser:  config.Config.WeComConfigs[alias].Receiver,
			Msgtype: "textcard",
			Agentid: config.Config.WeComConfigs[alias].AgentId,
		}

		msg.TextCard.Title = grafanaMsg.Title
		msg.TextCard.Description = description
		msg.TextCard.URL = cardUrl

		data, err := json.Marshal(msg)
		if err != nil {
			isSucc = false
			logger.Errorf("generate text message failed, skip pushing, error: %v", err)
			continue
		}
		err = wecom.SendTextMessage(data, alias)
		if err != nil {
			isSucc = false
			logger.Errorf("push text message failed, error: %v", err)
			continue
		}
	}
	return isSucc
}
