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
	"strconv"
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
		if err := sendClasses(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "send classes failed"})
		}
	})
	router.PATCH("/classes", func(c *gin.Context) {
		if err := updateClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "update class failed"})
		}
	})

	// 人脸
	router.POST("/detect_face", func(c *gin.Context) {
		if err := detectFace(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "detect face failed"})
		}
	})

	// 点名
	router.POST("/face_count", func(c *gin.Context) {
		// TODO("face count")
	})

	// 教室状态
	router.POST("/classroom_stats", func(c *gin.Context) {
		if err := updateClassroomStats(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "update stats error"})
		}
	})
	router.GET("/classroom_stats", func(c *gin.Context) {
		if err := sendClassroomStats(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "get stats error"})
		}
	})

	// 设备
	router.POST("/devices", func(c *gin.Context) {
		if err := createDevice(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create device error"})
		}
	})
	router.GET("/devices", func(c *gin.Context) {
		if err := sendDevices(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "send device error"})
		}
	})
	router.PATCH("/devices", func(c *gin.Context) {
		// TODO("patch device")
	})

	// 学生
	router.POST("/students", func(c *gin.Context) {
		if err := createStudent(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create student error"})
		}
	})
	router.GET("/students", func(c *gin.Context) {
		if err := sendStudents(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "send students error"})
		}
	})
	router.PATCH("/students", func(c *gin.Context) {
		// TODO("patch students")
	})

	// 摄像头
	router.POST("/cameras", func(c *gin.Context) {
		if err := createCamera(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create camera error"})
		}
	})
	router.GET("/cameras", func(c *gin.Context) {
		if err := sendCameras(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "send cameras error"})
		}
	})
	router.PATCH("/cameras", func(c *gin.Context) {
		// TODO("patch cameras")
	})

	// 教室
	router.POST("/classrooms", func(c *gin.Context) {
		if err := createClassroom(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "create classroom error"})
		}
	})
	router.GET("/classrooms", func(c *gin.Context) {
		if err := sendClassrooms(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonError{Error: "send classrooms error"})
		}
	})
	router.PATCH("/classrooms", func(c *gin.Context) {
		// TODO("patch classrooms")
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

	var classResponse ClassResponse
	err = json.Unmarshal(body, &classResponse)
	if err != nil {
		return err
	}

	class := Class{
		FaceSetToken: classResponse.FaceSetToken,
		ClassName: &className,
	}
	err = createTableItem(&class)
	if err != nil {
		return err
	}

	classResponse.ClassName = className
	classResponse.ClassID = class.ID
	c.JSON(http.StatusOK, classResponse)
	return nil
}

func sendClasses(c *gin.Context) (err error) {
	classID, isNotMulti := c.GetQuery("class_id")

	if !isNotMulti {
		var classes []Class
		classes, err = getAllClasses()
		if err != nil {
			return
		}

		classesResp := make([]ClassResponse, len(classes))
		for i := 0; i < len(classes); i++ {
			classesResp[i].ClassID = classes[i].ID
			classesResp[i].ClassName = *classes[i].ClassName
			classesResp[i].FaceCount = len(classes[i].Students)
			classesResp[i].FaceSetToken = classes[i].FaceSetToken
		}

		c.JSON(http.StatusOK, ClassesResponse{
			Classes: classesResp,
		})
	} else {
		var id int
		id, err = strconv.Atoi(classID)
		if err != nil {
			return
		}

		var class *Class
		class, err = getClass(id)
		if err != nil {
			return
		}

		studentNos := make([]string, len(class.Students))
		for k, v := range class.Students {
			studentNos[k] = *v.StudentNo
		}

		c.JSON(http.StatusOK, ClassResponse{
			ClassID: class.ID,
			ClassName: *class.ClassName,
			FaceCount: len(class.Students),
			FaceSetToken: class.FaceSetToken,
			StudentNos: studentNos,
		})
	}

	return
}

func updateClass(c *gin.Context) error {
	// TODO("update class")
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
	//var body []byte
	//body, err = ioutil.ReadAll(response.Body)
	//if err != nil {
	//	return err
	//}

	//var faceCountToken FaceCountToken
	//err = json.Unmarshal(body, &faceCountToken)
	//if err != nil {
	//	return err
	//}
	//
	//c.JSON(http.StatusOK, faceCountToken)
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

func updateClassroomStats(c *gin.Context) (err error) {
	var stats Stats
	err = c.ShouldBindJSON(&stats)
	if err != nil {
		return
	}

	devicePath := strings.Split(c.Request.Host, ":")[0]
	device, err := getDeviceByPath(devicePath)
	if err != nil {
		return
	}

	err = createTableItem(&DeviceStatsTable{
		UpdateTime: stats.UpdateTime,
		CpuUsed: stats.SystemStats.CpuUsed,
		MemUsed: stats.SystemStats.MemUsed,
		DeviceID: device.ID,
	})
	if err != nil {
		return
	}

	var classroom *Classroom
	for _, classroomStats := range stats.Classrooms {
		classroom, err = getClassroom(int(classroomStats.ClassroomID))
		if err != nil {
			return
		}

		err = createTableItem(&ClassroomStatsTable{
			UpdateTime: stats.UpdateTime,
			PersonCount: classroomStats.PersonCount,
			Persons: classroomStats.Persons,
			ClassroomID: classroom.ID,
		})
		if err != nil {
			return
		}
	}

	return
}

func sendClassroomStats(c *gin.Context) (err error) {
	classroomID, isExist := c.GetQuery("classroom_id")
	if isExist {
		return fmt.Errorf("no classroom_name")
	}

	var id int
	id, err = strconv.Atoi(classroomID)
	if err != nil {
		return
	}
	classroomStatsItem, err := getClassroomStatsItem(id)
	if err != nil {
		return
	}

	stats := SingleClassroomStats{
		UpdateTime: classroomStatsItem.UpdateTime,
		ClassroomStats: ClassroomStats{
			ClassroomID: classroomStatsItem.ClassroomID,
			PersonCount: classroomStatsItem.PersonCount,
			Persons: classroomStatsItem.Persons,
		},
	}

	c.JSON(http.StatusOK, stats)
	return
}

func createDevice(c *gin.Context) (err error) {
	devicePath := c.PostForm("device_path")
	devicePort := c.PostForm("device_port")

	device := Device{
		DevicePath: devicePath,
		DevicePort: devicePort,
	}
	err = createTableItem(&device)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, DeviceResponse{
		DeviceID: device.ID,
		DevicePath: devicePath,
		DevicePort: devicePort,
	})
	return
}

func sendDevices(c *gin.Context) (err error) {
	deviceID, isNotMulti := c.GetQuery("device_id")

	if !isNotMulti {
		var devices []Device
		devices, err = getAllDevices()
		if err != nil {
			return
		}

		devicesResponse := make([]DeviceResponse, len(devices))
		for i := 0; i < len(devicesResponse); i++ {
			devicesResponse[i].DeviceID = devices[i].ID
			devicesResponse[i].DevicePath = devices[i].DevicePath
			devicesResponse[i].DevicePort = devices[i].DevicePort
		}
		c.JSON(http.StatusOK, DevicesResponse{Devices: devicesResponse})
	} else {
		var id int
		id, err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}

		var device *Device
		device, err = getDevice(id)

		c.JSON(http.StatusOK, device)
	}
	return
}

func createStudent(c *gin.Context) (err error) {
	studentNo := c.PostForm("student_no")
	faceToken := c.PostForm("face_token")
	classIDs := c.PostFormArray("[]class_ids")

	var ids []int
	ids, err = stringArrayToIntArray(classIDs)
	if err != nil {
		return
	}

	classes, err := getClasses(ids)
	if err != nil {
		return
	}

	classPointers := make([]*Class, len(classes))
	for k, v := range classes {
		classPointers[k] = &v
	}
	student := Student{
		StudentNo: &studentNo,
		FaceToken: faceToken,
		Classes: classPointers,
	}

	err = createTableItem(&student)
	if err != nil {
		return
	}

	classUintIDs := make([]uint, len(classes))
	for k, v := range classes {
		classUintIDs[k] = v.ID
	}
	c.JSON(http.StatusOK, StudentResponse{
		StudentNo: *student.StudentNo,
		FaceToken: student.FaceToken,
		ClassIDs: classUintIDs,
	})
	return
}

func sendStudents(c *gin.Context) (err error) {
	studentNo, byStuNo := c.GetQuery("student_no")
	classID , byClassID := c.GetQuery("class_id")

	if byStuNo {
		var student *Student
		student, err = getStudent(studentNo)
		if err != nil {
			return
		}

		classUintIDs := make([]uint, len(student.Classes))
		for k, v := range student.Classes {
			classUintIDs[k] = v.ID
		}
		c.JSON(http.StatusOK, StudentResponse{
			StudentNo: *student.StudentNo,
			FaceToken: student.FaceToken,
			ClassIDs: classUintIDs,
		})
	} else if byClassID {
		var id int
		id ,err = strconv.Atoi(classID)
		if err != nil {
			return
		}

		var students []Student
		students, err = getStudentsByClass(id)
		if err != nil {
			return
		}

		studentsResponse := make([]StudentResponse, len(students))
		for k, v := range students {
			studentsResponse[k].FaceToken = v.FaceToken
			studentsResponse[k].StudentNo = *v.StudentNo
		}

		c.JSON(http.StatusOK, StudentsResponse{
			Students: studentsResponse,
		})
	} else {
		err = fmt.Errorf("no param provide")
		return
	}

	return
}

func stringArrayToIntArray(strArray []string) (intArray []int, err error) {
	intArray = make([]int, len(strArray))

	for i := 0; i < len(intArray); i++ {
		intArray[i], err = strconv.Atoi(strArray[i])
		if err != nil {
			return
		}
	}
	return
}

func createCamera(c *gin.Context) (err error) {
	camPath := c.PostForm("cam_path")
	deviceID := c.PostForm("device_id")

	var id int
	id ,err = strconv.Atoi(deviceID)
	if err != nil {
		return
	}

	camera := Camera{
		CamPath: camPath,
		DeviceID: uint(id),
	}

	err = createTableItem(&camera)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, CameraResponse{
		CameraID: camera.ID,
		DeviceID: camera.DeviceID,
		CamPath: camera.CamPath,
	})
	return
}

func sendCameras(c *gin.Context) (err error) {
	cameraID, isNotMulti := c.GetQuery("camera_id")

	if !isNotMulti {
		var cameras []Camera

		cameras, err = getCameras()
		if err != nil {
			return
		}

		camerasResponse := make([]CameraResponse, len(cameras))
		for k, v := range cameras {
			camerasResponse[k].CamPath = v.CamPath
			camerasResponse[k].CameraID = v.ID
			camerasResponse[k].DeviceID = v.DeviceID
		}

		c.JSON(http.StatusOK, CamerasResponse{
			Cameras: camerasResponse,
		})
	} else {
		var camera *Camera

		var id int
		id ,err = strconv.Atoi(cameraID)
		if err != nil {
			return
		}
		camera, err = getCamera(id)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, CameraResponse{
			CameraID: camera.ID,
			CamPath: camera.CamPath,
			DeviceID: camera.DeviceID,
		})
	}
	return
}

func createClassroom(c *gin.Context) (err error) {
	classroomName := c.PostForm("classroom_name")
	cameraID := c.PostForm("camera_id")

	var id int
	id ,err = strconv.Atoi(cameraID)
	if err != nil {
		return
	}

	classroom := Classroom{
		Name: classroomName,
		CameraID: uint(id),
	}

	err = createTableItem(&classroom)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ClassroomResponse{
		ClassroomID: classroom.ID,
		ClassroomName: classroom.Name,
		CameraID: classroom.CameraID,
	})
	return
}

func sendClassrooms(c *gin.Context) (err error) {
	classroomID, isNotMulti := c.GetQuery("classroom_id")

	if !isNotMulti {
		var classrooms []Classroom

		classrooms, err = getClassrooms()
		if err != nil {
			return
		}

		classroomsResponse := make([]ClassroomResponse, len(classrooms))
		for k, v := range classrooms {
			classroomsResponse[k].CameraID = v.CameraID
			classroomsResponse[k].ClassroomName = v.Name
			classroomsResponse[k].ClassroomID = v.ID
		}

		c.JSON(http.StatusOK, ClassroomsResponse{
			Classrooms: classroomsResponse,
		})
	} else {
		var classroom *Classroom

		var id int
		id ,err = strconv.Atoi(classroomID)
		if err != nil {
			return
		}
		classroom, err = getClassroom(id)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, ClassroomResponse{
			ClassroomID: classroom.ID,
			CameraID: classroom.ID,
			ClassroomName: classroom.Name,
		})
	}
	return
}