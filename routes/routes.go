package routes

import (
	"bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	return r
}
