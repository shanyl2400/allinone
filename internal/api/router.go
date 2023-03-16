package api

import "github.com/gin-gonic/gin"

func route(e *gin.Engine) {
	_api := newAPI()
	e.GET("/gomss/branches", _api.getGomssBranches)
	e.GET("/zrtc/path", _api.getZRTCPath)
	e.GET("/publish/records", _api.getRecentPublish)
	e.GET("/publish/logs", _api.getPublishLogs)
	e.POST("/publish", _api.publish)
}
