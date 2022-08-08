package api

import (
	"github.com/TongboZhang/wecom-pusher/logger"
	"github.com/gin-gonic/gin"
)

func Start() (err error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(logger.ApiLogger())
	SetRoute(engine)
	err = engine.Run("0.0.0.0:18080")
	return err
}

func SetRoute(r *gin.Engine) {
	r.POST("/push", TokenAuth, Push)
}
