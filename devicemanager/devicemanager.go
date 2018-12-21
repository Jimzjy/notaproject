package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const ConfigFileName = "config.json"

var config Config

func main() {
	var err error

	getConfig(&config)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := setupRouter()
	router.Run(config.LocalPort)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// 班级
	router.POST("/classes", func(c *gin.Context) {
		if err := createClass(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create class failed"})
		}
	})
	router.GET("/classes", func(c *gin.Context) {

	})
	router.PATCH("/classes", func(c *gin.Context) {

	})

	// 人脸
	router.POST("/faces", func(c *gin.Context) {

	})
	router.GET("/faces", func(c *gin.Context) {

	})
	router.PATCH("/faces", func(c *gin.Context) {

	})

	// 图片
	router.GET("/images/:name", func(c *gin.Context) {

	})

	// 点名
	router.POST("/face_count", func(c *gin.Context) {

	})

	// 教室状态
	router.POST("/classrooms", func(c *gin.Context) {

	})
	router.GET("/classrooms", func(c *gin.Context) {

	})

	// 设备
	router.POST("/devices", func(c *gin.Context) {

	})
	router.GET("/devices", func(c *gin.Context) {

	})
	router.PATCH("/devices", func(c *gin.Context) {

	})

	// 设置
	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, config)
	})
	router.POST("/config", func(c *gin.Context) {
		if err := setConfig(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonError{Error: "can not set config"})
		}
	})

	return router
}

// 获取设置文件信息
func getConfig(config *Config) error {
	var err error

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

func setConfig(c *gin.Context) error {
	var err error

	var reqConfig Config
	err = c.ShouldBindJSON(&reqConfig)
	if err != nil {
		return err
	}

	tmp, err := json.Marshal(reqConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(ConfigFileName, tmp, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createClass(c *gin.Context) error {
	var err error

	className := c.PostForm("class_name")

	response, err := http.PostForm(config.DetectFaceUrl, url.Values{
		"api_key": {config.ApiKey},
		"api_secret": {config.ApiSecret},
		"display_name": {className},
	})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var faceCountToken FaceCountToken
	err = json.Unmarshal(body, &faceCountToken)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, faceCountToken)
	return nil
}