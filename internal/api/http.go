package api

import "github.com/gin-gonic/gin"

func Start() {
	r := gin.Default()
	route(r)
	r.Run(":8088") // listen and serve on 0.0.0.0:8080
}
