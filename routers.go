package main

import (
	"github.com/foolish06/gin-essential/controller"
	"github.com/gin-gonic/gin"
)

func collectRouter(router *gin.Engine)  {
	router.POST("/api/auth/register", controller.Register)
}
