package main

import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gocv.io/x/gocv"
	"image"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"
	"unsafe"
)
// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -fopenmp -L${SRCDIR}/lib -lncnnwrapper -lncnn -lm -lstdc++ -lopencv_core -lopencv_videoio
// #include "ncnnwrapper.h"
import "C"

const (
	FaceDetect = 0
	BodyDetect = 1

	ConfigFileName = "config.json"
)

var config Config

func main() {
	var err error

	getConfig(&config)
	if err != nil {
		log.Fatal(err)
		return
	}

	go uploadStats()

	router := setupRouter()
	router.Run(config.LocalPort)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// 人脸识别
	router.POST("/face_search", func(c *gin.Context) {
		searchFace(c)
	})

	// 设置
	router.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, config)
	})
	router.POST("/config", func(c *gin.Context) {
		if err := setConfig(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonError{Error: "set config failed"})
		}
	})

	return router
}

// 获取设置文件信息
func getConfig(config *Config) error {
	var err error

	data, err := ioutil.ReadFile(ConfigFileName)
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

// POST /faceSearch
func searchFace(c *gin.Context) {
	var err error

	camPath := c.PostForm("cam_path")
	faceSetToken := c.PostForm("faceset_token")

	personData, err := getSearchData(camPath, faceSetToken)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, JsonError{Error: "can not get person data"})
		return
	}

	c.JSON(http.StatusOK, PersonDataFaces{Faces: personData})
}

// 获得 人脸 Rect 和 Token
func getSearchData(camPath, faceSetToken string) ([]PersonData, error) {
	var err error

	detectedImage, err := getDetectedImage(camPath, FaceDetect)
	if err != nil {
		return nil, err
	}

	personData := make([]PersonData, len(detectedImage))
	params := map[string]string{
		"api_key": config.ApiKey,
		"api_secret": config.ApiSecret,
		"faceset_token": faceSetToken,
	}
	for i := 0; i < len(detectedImage); i++ {
		_image := detectedImage[i]
		response, err := fileUploadRequest(config.SearchFaceUrl, params,
			"image_file", _image.Data, "image.jpg")
		if err != nil {
			return nil, err
		}
		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, fmt.Errorf(response.Status)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			response.Body.Close()
			return nil, err
		}

		response.Body.Close()

		var searchFaceResponse SearchFaceResults
		err = json.Unmarshal(body, &searchFaceResponse)
		if err != nil {
			return nil, err
		}
		results := searchFaceResponse.Results
		if len(results) > 0 {
			confidence := results[0].Confidence
			if confidence > 75 {
				personData[i] = PersonData{_image.DetectedData, results[0].FaceToken}
			} else {
				personData[i] = PersonData{_image.DetectedData, ""}
			}
		} else {
			personData[i] = PersonData{_image.DetectedData, ""}
		}
	}

	return personData, nil
}

// 获得当前 Mat
func getCameraImage(camPath string, img *gocv.Mat) error {
	webCam, err := gocv.OpenVideoCapture(camPath)
	if err != nil {
		return err
	}
	defer webCam.Close()
	if ok := webCam.Read(img); !ok {
		return fmt.Errorf("can not read from webCam")
	}

	return nil
}

// 获得输入 camPath 中的人脸检测结果
// return 人脸 Rect 和 jpg
func getDetectedImage(camPath string, mode int) ([]DetectedImage, error) {
	var err error

	img := gocv.NewMat()
	defer img.Close()
	err = getCameraImage(camPath, &img)
	if err != nil {
		return nil, err
	}

	detectedData, err := getDetectedData(img, mode)
	if err != nil {
		return nil, err
	}

	detectedImage := make([]DetectedImage, len(detectedData))
	for i := 0; i < len(detectedData); i++  {
		rect := detectedData[i]
		_img, err := img.FromRect(image.Rect(rect.X0, rect.Y0, rect.X1, rect.Y1))

		detectedImage[i].Data, err = gocv.IMEncode(".jpg", _img)
		detectedImage[i].DetectedData = rect
		if err != nil {
			return nil, err
		}
	}

	return detectedImage, nil
}


// 获得输入 Mat 中的人脸检测结果
// return 人脸 Rect
func getDetectedData(img gocv.Mat, mode int) ([]DetectedData, error) {
	ncnnnet := C.newNcnnnet()
	defer C.free(unsafe.Pointer(ncnnnet))
	var _param, _model string

	switch mode {
	case FaceDetect:
		_param = "ssd_face.param"
		_model = "ssd_face.bin"
		break
	case BodyDetect:
		_param = "ssd_body.param"
		_model = "ssd_body.bin"
		break
	default:
	}

	param := C.CString(_param)
	model := C.CString(_model)
	defer C.free(unsafe.Pointer(param))
	defer C.free(unsafe.Pointer(model))
	C.ncnnnetLoad(param, model, ncnnnet)

	data := (*C.uchar)(unsafe.Pointer(&(img.DataPtrUint8()[0])))
	rects := C.detectFromByte(data, C.int(img.Cols()), C.int(img.Rows()), ncnnnet, C.int(mode))
	rectsPointer := rects.rects
	size := int(rects.size)
	header := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(rectsPointer)),
		Len:  size,
		Cap:  size,
	}
	rectsSlice := *(*[]C.Rect)(unsafe.Pointer(header))

	detectedData := make([]DetectedData, len(rectsSlice))

	bias := 20
	for i := 0; i < len(rectsSlice); i++ {
		rect := rectsSlice[i]
		detectedData[i] = DetectedData{int(rect.x0) - bias, int(rect.y0) - bias, int(rect.x1) + bias, int(rect.y1) + bias}
	}

	return detectedData, nil
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

func uploadStats() {
	ticker := time.NewTicker(time.Duration(config.DetectInterval * 1000000000))
	var err error

	for range ticker.C {
		err = uploadStatsRequest()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func uploadStatsRequest() error {
	var err error

	img := gocv.NewMat()
	defer img.Close()
	classroomStats := make([]ClassroomStat, len(config.Classrooms))
	for i := 0; i < len(classroomStats); i++ {
		err = getCameraImage(config.Classrooms[i].CamPath, &img)
		if err != nil {
			return err
		}

		var data []DetectedData
		data, err = getDetectedData(img, BodyDetect)
		if err != nil {
			return err
		}
		classroomStats[i].Persons = data
		classroomStats[i].PersonCount = len(data)
		classroomStats[i].Classroom = config.Classrooms[i].Classroom
	}

	cpuPer, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	memStats, err := mem.VirtualMemory()
	stats := Stats{
		UpdateTime: time.Now(),
		Stats: SystemStats{
			CpuUsed: arrayToString(cpuPer),
			MemUsed: fmt.Sprintf("%v %v", memStats.Total, memStats.Active),
		},
		Classrooms: classroomStats,
	}
	jsonStats, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%v/classrooms", config.ServerAddr)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStats))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(response.Status)
	}

	//if response.StatusCode != http.StatusOK {
	//	var body []byte
	//	body, err = ioutil.ReadAll(response.Body)
	//	if err != nil {
	//		return err
	//	}
	//	return fmt.Errorf(string(body))
	//}

	return nil
}

func arrayToString(v interface{}) string {
	tmp := fmt.Sprint(v)
	last := len(tmp) - 1

	if tmp[0] == '[' && tmp[last] == ']' {
		return tmp[1:last]
	}

	return tmp
}