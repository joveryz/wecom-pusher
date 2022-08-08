package api

import (
	"github.com/TongboZhang/wecom-pusher/config"
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/gin-gonic/gin"
)

func BasicAuth() func(context *gin.Context) {
	return gin.BasicAuth(gin.Accounts{
		config.Config.User: config.Config.Password,
	})
}

func TokenAuth(context *gin.Context) {
	token := parseFromContext(context, "token")
	if token != config.Config.Token {
		logger.Warnf("token validation failed, token: %s", token)
		context.String(401, "token validation failed")
		context.Abort()
		return
	}
}
