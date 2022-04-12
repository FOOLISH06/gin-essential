package main

import (
	"github.com/foolish06/gin-essential/controller"
	"github.com/foolish06/gin-essential/middlerware"
	"github.com/gin-gonic/gin"
)

func collectRouter(router *gin.Engine) *gin.Engine {
	router.POST("/api/auth/register", controller.Register)
	router.POST("/api/auth/login", controller.Login)
	router.GET("/api/auth/info", middlerware.AuthMiddleware(), controller.Info)

	return router
}
