package api

import (
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
