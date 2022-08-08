package api

import (
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/gin-gonic/gin"
)

var handler func(context *gin.Context)

func TokenAuth(context *gin.Context) {
	token := parseFromContext(context, "token")
	if token == config.Config.Token {
		logger.Info("token validation passed")
		return
	}

	if handler == nil {
		handler = gin.BasicAuth(gin.Accounts{config.Config.User: config.Config.Password})
	}

	handler(context)
}
