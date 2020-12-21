package server

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes() {
	r := gin.Default()
	r.GET("/health-check", healthCheck)
	r.POST("/users", addNewUser)
	r.GET("/users", getAllUsers)
	r.GET("/users/:id", getItemById)
	r.Run()
}