package server

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes() {
	r := gin.Default()
	r.GET("/health-check", healthCheck)
	r.POST("/user", addNewUser)
	r.Run()
}