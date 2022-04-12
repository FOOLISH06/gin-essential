package main

import (
	"github.com/foolish06/gin-essential/common"
	"github.com/gin-gonic/gin"
	"log"
)



func main() {
	_ = common.GetDB()

	router := gin.Default()
	_ = router.SetTrustedProxies([]string{"localhost"})
	router = collectRouter(router)

	if err := router.Run(); err != nil {
		log.Fatalln(err.Error())
	} // listen and serve on 0.0.0.0:8080
}





