package main

import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gocv.io/x/gocv"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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
	RootAddress = "/home/pi/embeddedSystem/"
	ConfigFileName = RootAddress + "config.json"
)

var (
	config Config
	globalWidth float64
	globalHeight float64
	)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var ncnnnetFace = C.newNcnnnet()
var ncnnnetBody = C.newNcnnnet()
var client = &http.Client{
	Timeout: time.Second * 5,
}

func main() {
	var err error

	gin.SetMode(gin.ReleaseMode)

	err = getConfig(&config)
	if err != nil {
		log.Println(err)
		return
	}

	paramFace := C.CString(config.FaceDetectParam)
	modelFace := C.CString(config.FaceDetectBin)
	paramBody := C.CString(config.BodyDetectParam)
	modelBody := C.CString(config.BodyDetectBin)
	C.ncnnnetLoad(paramFace, modelFace, ncnnnetFace)
	C.ncnnnetLoad(paramBody, modelBody, ncnnnetBody)
	defer C.free(unsafe.Pointer(ncnnnetFace))
	defer C.free(unsafe.Pointer(ncnnnetBody))
	C.free(unsafe.Pointer(paramFace))
	C.free(unsafe.Pointer(modelFace))
	C.free(unsafe.Pointer(paramBody))
	C.free(unsafe.Pointer(modelBody))

	//go uploadStats()

	router := setupRouter()
	err = router.Run(config.LocalPort)
	if err != nil {
		log.Println(err)
	}
}

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

func sendPersonStatus(c *gin.Context) (err error) {
	camStreamPath := c.PostForm("cam_stream_path")
	basicSize := 5

	_, faceDetectResults, err := getFacePPDetect(camStreamPath)
	if err != nil {
		return
	}

	if len(faceDetectResults.Faces) <= basicSize {
		c.JSON(http.StatusOK, faceDetectResults)
	} else {
		var body []byte
		var err2 error
		facesTokens := []string{""}
		count := 0
		index := 0

		for i := 5; i < len(faceDetectResults.Faces); i++ {
			facesTokens[index] += faceDetectResults.Faces[i].FaceToken + ","
			count++
			if count >= 5 {
				facesTokens[index] = strings.TrimSuffix(facesTokens[index], ",")
				facesTokens = append(facesTokens, "")
				index++
				count = 0
			}
		}

		facesTokens[len(facesTokens) - 1] = strings.TrimSuffix(facesTokens[len(facesTokens) - 1], ",")

		for k, v := range facesTokens {
			body, err2 = sendPostForm(url.Values{
				"api_key": {config.ApiKey},
				"api_secret": {config.ApiSecret},
				"face_tokens": {v},
				"return_attributes": {"headpose,eyestatus,emotion"},
			}, config.AnalyzeFaceUrl)
			if err2 != nil {
				log.Println(err2)
				continue
			}

			var _faceAnalyzeResults FaceDetectResults
			err2 = json.Unmarshal(body, &_faceAnalyzeResults)
			if err2 != nil {
				log.Println(err2)
				continue
			}

			for k2, v2 := range _faceAnalyzeResults.Faces {
				faceDetectResults.Faces[basicSize + basicSize * k + k2].Attributes = v2.Attributes
			}
		}

		c.JSON(http.StatusOK, faceDetectResults)
	}
	return
}

func searchFace(c *gin.Context) {
	var err error

	camStreamPath := c.Query("camStreamPath")
	faceSetToken := c.Query("faceSetToken")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	done := make(chan string)
	go func() {
		defer close(done)

		for {
			_, _, err = conn.ReadMessage()
			if err != nil {
				log.Println("error from read from deviceManager:", err)
				break
			}
		}
	}()

	chPersonData := make(chan PersonData)
	go getSearchData(camStreamPath, faceSetToken, chPersonData)

	for {
		select {
		case <-done:
			return
		case personData, ok := <-chPersonData:
			if ok {
				var data []byte
				data, err = json.Marshal(personData)
				if err != nil {
					continue
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("error from write to deviceManager", err)
					return
				}
			} else {
				err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("error from write to deviceManager", err)
					return
				}
				return
			}
		}
	}
}

func getSearchData(camStreamPath, faceSetToken string, chPersonData chan PersonData) {
	var err error
	defer close(chPersonData)

	gImage, faceDetectResults, err := getFacePPDetect(camStreamPath)
	if err != nil {
		log.Println(err)
		return
	}

	var body []byte
	body, err = fileUploadRequest(fmt.Sprintf("http://%v/images", config.ServerAddr),
		map[string]string{}, "file", gImage, "image.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	var jsonMessage JsonMessage
	err = json.Unmarshal(body, &jsonMessage)
	if err != nil {
		log.Println(err)
		return
	}

	personCount := len(faceDetectResults.Faces)
	personData := PersonData{
		PersonCount: personCount,
		ImageUrl: jsonMessage.Message,
		GlobalWidth: globalWidth,
		GlobalHeight: globalHeight,
	}

	chPostCtrl := make(chan string, config.Qps)
	chActionFinished := make(chan byte)
	count := 0
	for _, v := range faceDetectResults.Faces {
		_personData := personData

		go sendFacePostForm(faceSetToken, v, _personData, chPersonData, chPostCtrl, &count, chActionFinished, personCount)
	}

	<- chActionFinished
}

func sendFacePostForm(faceSetToken string, face FaceAnalyzeResult, _personData PersonData, chPersonData chan PersonData, chPostCtrl chan string, count *int, chActionFinished chan byte, personCount int) {
	chPostCtrl <- ""
	defer func() {
		if *count < personCount - 1 {
			*count += 1
		} else {
			chActionFinished <- 0
		}
		<- chPostCtrl
	}()

	body, err := sendPostForm(url.Values{
		"api_key": {config.ApiKey},
		"api_secret": {config.ApiSecret},
		"faceset_token": {faceSetToken},
		"face_token": {face.FaceToken},
	}, config.SearchFaceUrl)
	if err != nil {
		log.Println(err)
		_personData.Face.FaceRectangle = face.FaceRectangle
		chPersonData <- _personData
		return
	}

	var searchFaceResponse SearchFaceResults
	err = json.Unmarshal(body, &searchFaceResponse)
	if err != nil {
		log.Println(err)
		_personData.Face.FaceRectangle = face.FaceRectangle
		chPersonData <- _personData
		return
	}
	results := searchFaceResponse.Results
	if len(results) > 0 {
		confidence := results[0].Confidence
		if confidence > 70 {
			_personData.Face.FaceToken = results[0].FaceToken
			_personData.Face.FaceRectangle = face.FaceRectangle
			chPersonData <- _personData
		} else {
			_personData.Face.FaceRectangle = face.FaceRectangle
			chPersonData <- _personData
		}
	} else {
		_personData.Face.FaceRectangle = face.FaceRectangle
		chPersonData <- _personData
	}
}

func getCameraImage(camStreamPath string, img *gocv.Mat) error {
	webCam, err := gocv.OpenVideoCapture(camStreamPath)
	if err != nil {
		return err
	}
	defer webCam.Close()

	globalWidth = webCam.Get(gocv.VideoCaptureFrameWidth)
	globalHeight = webCam.Get(gocv.VideoCaptureFrameHeight)

	if ok := webCam.Read(img); !ok {
		return fmt.Errorf("can not read from webCam")
	}

	return nil
}

func getFacePPDetect(camStreamPath string) (gImage []byte, faceDetectResults FaceDetectResults, err error) {
	img := gocv.NewMat()
	defer img.Close()
	err = getCameraImage(camStreamPath, &img)
	if err != nil {
		return
	}

	gImage, err = gocv.IMEncode(".jpg", img)
	if err != nil {
		return
	}

	faceDetectBody, err := fileUploadRequest(config.DetectFaceUrl, map[string]string{
		"api_key": config.ApiKey,
		"api_secret": config.ApiSecret,
		"return_attributes": "headpose,eyestatus,emotion",
	}, "image_file", gImage, "image.jpg")
	if err != nil {
		return
	}

	err = json.Unmarshal(faceDetectBody, &faceDetectResults)
	if err != nil {
		return
	}

	return
}

func getDetectedData(img gocv.Mat, mode int) ([]FaceRectangle, error) {
	bias := 0

	data := (*C.uchar)(unsafe.Pointer(&(img.DataPtrUint8()[0])))
	var rects C.Rects
	switch mode {
	case FaceDetect:
		bias = 10
		rects = C.detectFromByte(data, C.int(img.Cols()), C.int(img.Rows()), ncnnnetFace, C.int(mode))
	case BodyDetect:
		rects = C.detectFromByte(data, C.int(img.Cols()), C.int(img.Rows()), ncnnnetBody, C.int(mode))
	}

	rectsPointer := rects.rects
	size := int(rects.size)
	header := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(rectsPointer)),
		Len:  size,
		Cap:  size,
	}
	rectsSlice := *(*[]C.Rect)(unsafe.Pointer(header))

	detectedData := make([]FaceRectangle, len(rectsSlice) - 1)

	for i := 1; i < len(rectsSlice); i++ {
		rect := rectsSlice[i]
		detectedData[i-1] = FaceRectangle{int(rect.top) - bias, int(rect.left) - bias, int(rect.width) + bias, int(rect.height) + bias}
		//fmt.Println(detectedData[i-1])
	}

	return detectedData, nil
}

func fileUploadRequest(url string, params map[string]string, fileParamName string, fileContent []byte, fileName string) (bodyResp []byte, err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileParamName, fileName)
	if err != nil {
		return
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return
	}

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return
		}
	}
	err = writer.Close()
	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("response not ok for upload file")

		//var body2 []byte
		//body2, err = ioutil.ReadAll(response.Body)
		//if err != nil {
		//	return
		//}
		//err = fmt.Errorf(string(body2))

		return
	}

	bodyResp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}

func uploadStats() {
	ticker := time.NewTicker(time.Duration(config.DetectInterval * 1e9))
	var err error

	for range ticker.C {
		err = uploadStatsRequest()
		if err != nil {
			log.Println(err)
		}
	}
}

func uploadStatsRequest() error {
	var err error

	img := gocv.NewMat()
	defer img.Close()
	classroomStats := make([]ClassroomStat, len(config.Classrooms))
	for i := 0; i < len(classroomStats); i++ {
		err = getCameraImage(config.Classrooms[i].CamStreamPath, &img)
		if err != nil {
			return err
		}

		var data []FaceRectangle
		data, err = getDetectedData(img, BodyDetect)
		if err != nil {
			return err
		}
		classroomStats[i].Persons = data
		classroomStats[i].PersonCount = len(data)
		classroomStats[i].ClassroomNo = config.Classrooms[i].ClassroomNo
		classroomStats[i].GlobalWidth = globalWidth
		classroomStats[i].GlobalHeight = globalHeight
	}

	cpuUsed, err := cpu.Percent(time.Second, false)
	if err != nil {
		return err
	}
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	stats := Stats{
		UpdateTime: time.Now(),
		Stats: SystemStats{
			CpuUsed: getFloatPrecision(cpuUsed[0], "1"),
			MemUsed: getFloatPrecision(
				(float64(memStats.Active) / float64(memStats.Total)) * 100,
				"1"),
		},
		Classrooms: classroomStats,
	}
	jsonStats, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	_url := fmt.Sprintf("http://%v/classroom_stats", config.ServerAddr)
	request, err := http.NewRequest("POST", _url, bytes.NewBuffer(jsonStats))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response not ok for upload stats")
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

func getFloatPrecision(number float64, p string) float64 {
	format := "%." + p + "f"

	numberStr := fmt.Sprintf(format, number)

	numberP, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0
	}

	return numberP
}

// still have problem of server hangs, maybe it's the bug of nanoPi
func sendPostForm(params url.Values, url string) (body []byte, err error) {
	response, err := client.PostForm(url, params)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("response not ok for send post form")
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}