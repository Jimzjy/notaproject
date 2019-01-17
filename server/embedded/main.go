package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	var err error

	err = getConfig(&config)
	if err != nil {
		log.Println(err)
		return
	}

	go uploadStats()

	router := setupRouter()
	err = router.Run(config.LocalPort)
	if err != nil {
		log.Println(err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// 人脸识别
	router.GET("/face_search", searchFace)

	// 设置
	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, config)
	})
	router.POST("/config", func(c *gin.Context) {
		if err := setConfig(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "set config failed"})
		}
	})

	return router
}
