package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const ConfigFileName = "config.json"

var config Config

func main() {
	var err error

	getConfig(&config)
	if err != nil {
		log.Println(err)
		return
	}

	router := setupRouter()
	router.Run(config.LocalPort)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.MaxMultipartMemory = 2 << 20

	// 班级
	router.POST("/classes", func(c *gin.Context) {
		if err := createClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create class failed"})
		}
	})
	router.GET("/classes", func(c *gin.Context) {
		sendClasses(c)
	})
	router.PATCH("/classes", func(c *gin.Context) {
		if err := updateClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "update class failed"})
		}
	})

	// 人脸
	router.POST("/faces", func(c *gin.Context) {
		if err := addFace(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "add face failed"})
		}
	})
	router.POST("/detect_face", func(c *gin.Context) {
		if err := detectFace(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "detect face failed"})
		}
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

func sendClasses(c *gin.Context) {
	if classes, err := getAllClasses(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, JsonError{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, classes)
	}
}

func updateClass(c *gin.Context) error {
	return nil
}

func addFace(c *gin.Context) error {
	var err error

	faceToken := c.PostForm("face_token")
	faceSetToken := c.PostForm("faceset_token")

	response, err := http.PostForm(config.DetectFaceUrl, url.Values{
		"api_key": {config.ApiKey},
		"api_secret": {config.ApiSecret},
		"faceset_token": {faceSetToken},
		"face_token": {faceToken},
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

func detectFace(c *gin.Context) error {
	var err error

	fileHeader, err := c.FormFile("image_file")
	if err != nil {
		return err
	}
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	params := map[string]string{
		"api_key": config.ApiKey,
		"api_secret": config.ApiSecret,
	}
	response, err := fileUploadRequest(config.DetectFaceUrl, params,
		"image_file", data, fileHeader.Filename)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var faceRectTokens FaceRectTokens
	err = json.Unmarshal(body, &faceRectTokens)
	if err != nil {
		return err
	}
	if len(faceRectTokens.Faces) < 1 {
		c.JSON(http.StatusBadRequest, "no person in image")
		return nil
	}

	stuFace := FaceNoToken{
		FaceToken: faceRectTokens.Faces[0].FaceToken,
		StudentNo: strings.TrimSuffix(fileHeader.Filename, path.Ext(fileHeader.Filename)),
	}
	c.JSON(http.StatusOK, stuFace)

	return nil
}

// 文件上传 Request
// return Response
func fileUploadRequest(url string, params map[string]string, fileParamName string, fileContent []byte, fileName string) (*http.Response, error) {
	var err error

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileParamName, fileName)
	if err != nil {
		return nil, err
	}
	part.Write(fileContent)

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, err
}