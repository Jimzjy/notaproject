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
	"image"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
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

func main() {
	var err error

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

	go uploadStats()

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
				log.Println("read from deviceManager:", err)
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
					log.Println("write to deviceManager", err)
					return
				}
			} else {
				err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write to deviceManager", err)
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

	gImage, detectedImage, err := getDetectedImage(camStreamPath, FaceDetect)
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

	personData := PersonData{
		PersonCount: len(detectedImage),
		ImageUrl: jsonMessage.Message,
		GlobalWidth: globalWidth,
		GlobalHeight: globalHeight,
	}

	params := map[string]string{
		"api_key": config.ApiKey,
		"api_secret": config.ApiSecret,
		"faceset_token": faceSetToken,
	}
	for i := 0; i < len(detectedImage); i++ {
		_image := detectedImage[i]
		_personData := personData

		body, err = fileUploadRequest(config.SearchFaceUrl, params,
			"image_file", _image.Data, "image.jpg")
		if err != nil {
			log.Println(err)
			if err.Error() != "response not ok" {
				return
			} else {
				_personData.DetectedData = _image.DetectedData
				chPersonData <- _personData
				continue
			}
		}

		var searchFaceResponse SearchFaceResults
		err = json.Unmarshal(body, &searchFaceResponse)
		if err != nil {
			log.Println(err)
			return
		}
		results := searchFaceResponse.Results
		if len(results) > 0 {
			confidence := results[0].Confidence
			if confidence > 40 {
				_personData.DetectedData = _image.DetectedData
				_personData.Token = results[0].FaceToken
				chPersonData <- _personData
			} else {
				_personData.DetectedData = _image.DetectedData
				chPersonData <- _personData
			}
		} else {
			_personData.DetectedData = _image.DetectedData
			chPersonData <- _personData
		}
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

func getDetectedImage(camStreamPath string, mode int) (gImage []byte, detectedImage []DetectedImage, err error) {
	img := gocv.NewMat()
	defer img.Close()
	err = getCameraImage(camStreamPath, &img)
	if err != nil {
		return
	}

	detectedData, err := getDetectedData(img, mode)
	if err != nil {
		return
	}

	detectedImage = make([]DetectedImage, len(detectedData))
	for i := 0; i < len(detectedData); i++  {
		rect := detectedData[i]
		var _img gocv.Mat
		_img, err = img.FromRect(image.Rect(rect.X0, rect.Y0, rect.X1, rect.Y1))
		if err != nil {
			return
		}
		if rect.X1 - rect.X0 < 80 {
			fx := math.Ceil(80 / float64(rect.X1 - rect.X0))
			gocv.Resize(_img, &_img, image.Pt(0, 0), fx, fx, gocv.InterpolationArea)
		}

		detectedImage[i].Data, err = gocv.IMEncode(".jpg", _img)
		if err != nil {
			return
		}
		err = _img.Close()
		if err != nil {
			log.Println(err)
		}
		detectedImage[i].DetectedData = rect
	}

	gImage, err = gocv.IMEncode(".jpg", img)
	if err != nil {
		return
	}
	return
}

func getDetectedData(img gocv.Mat, mode int) ([]DetectedData, error) {
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

	detectedData := make([]DetectedData, len(rectsSlice) - 1)

	for i := 1; i < len(rectsSlice); i++ {
		rect := rectsSlice[i]
		detectedData[i-1] = DetectedData{int(rect.x0) - bias, int(rect.y0) - bias, int(rect.x1) + bias, int(rect.y1) + bias}
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

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("response not ok")
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

		var data []DetectedData
		data, err = getDetectedData(img, BodyDetect)
		if err != nil {
			return err
		}
		classroomStats[i].Persons = data
		classroomStats[i].PersonCount = len(data)
		classroomStats[i].ClassroomNo = config.Classrooms[i].ClassroomNo
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

	url := fmt.Sprintf("http://%v/classroom_stats", config.ServerAddr)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStats))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
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

func getFloatPrecision(number float64, p string) float64 {
	format := "%." + p + "f"

	numberStr := fmt.Sprintf(format, number)

	numberP, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0
	}

	return numberP
}