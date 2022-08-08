package api

import (
	"github.com/gin-gonic/gin"
)

func Push(context *gin.Context) {
	destination := parseDestinationFromContext(context)
	isSucc := generateAndPushWeComMessageFromContext(context, destination)
	if !isSucc {
		context.JSON(500, "push failed")
		return
	}

	context.String(200, "ok")
}
