package api

import (
	"fmt"
	"gomssbuilder/internal/config"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.Use(corsMiddleware())
	route(r)
	r.Run(fmt.Sprintf(":%d", config.GetConfig().HttpPort)) // listen and serve on 0.0.0.0:8080
}
