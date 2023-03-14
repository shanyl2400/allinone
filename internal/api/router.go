package api

import "github.com/gin-gonic/gin"

func route(e *gin.Engine) {
	_api := newAPI()
	e.GET("/gomss/branches", _api.getGomssBranches)
	e.GET("/zrtc/path", _api.getZRTCPath)
	e.POST("/gomss", _api.publish)
}
