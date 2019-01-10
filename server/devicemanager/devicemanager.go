package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	ConfigFileName = "config.json"
	UpdateStatsTime = 30
	ImageFileDir = "images/"
)

var config Config

func main() {
	var err error

	err = getConfig(&config)
	if err != nil {
		log.Println(err)
		return
	}

	go updateStats()

	router := setupRouter()
	err = router.Run(config.LocalPort)
	if err != nil {
		log.Println(err)
		return
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.MaxMultipartMemory = 2 << 20

	// 班级
	router.POST("/classes", func(c *gin.Context) {
		if err := createClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create class failed"})
		}
	})
	router.GET("/classes", func(c *gin.Context) {
		if err := sendClasses(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send classes failed"})
		}
	})
	router.PATCH("/classes", func(c *gin.Context) {
		if err := updateClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update class failed"})
		}
	})

	// 人脸
	router.POST("/detect_face", func(c *gin.Context) {
		if err := detectFace(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "detect face failed"})
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
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update stats error"})
		}
	})
	router.GET("/classroom_stats", func(c *gin.Context) {
		if err := sendClassroomStats(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "get stats error"})
		}
	})

	// 设备
	router.POST("/devices", func(c *gin.Context) {
		if err := createDevice(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create device error"})
		}
	})
	router.GET("/devices", func(c *gin.Context) {
		if err := sendDevices(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send device error"})
		}
	})
	router.PATCH("/devices", func(c *gin.Context) {
		// TODO("patch device")
	})

	// 学生
	router.POST("/students", func(c *gin.Context) {
		if err := createStudent(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create student error"})
		}
	})
	router.GET("/students", func(c *gin.Context) {
		if err := sendStudents(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send students error"})
		}
	})
	router.PATCH("/students/:no", func(c *gin.Context) {
		if err := updateStudent(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update student error"})
		}
	})
	router.DELETE("/students", func(c *gin.Context) {
		if err := deleteStudent(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete students error"})
		}
	})

	// 摄像头
	router.POST("/cameras", func(c *gin.Context) {
		if err := createCamera(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create camera error"})
		}
	})
	router.GET("/cameras", func(c *gin.Context) {
		if err := sendCameras(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send cameras error"})
		}
	})
	router.PATCH("/cameras", func(c *gin.Context) {
		// TODO("patch cameras")
	})

	// 教室
	router.POST("/classrooms", func(c *gin.Context) {
		if err := createClassroom(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create classroom error"})
		}
	})
	router.GET("/classrooms", func(c *gin.Context) {
		if err := sendClassrooms(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send classrooms error"})
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
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "can not set config"})
		}
	})

	// 管理员
	router.POST("/admin/login", func(c *gin.Context) {
		if err := adminLogin(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "user login error"})
		}
	})
	router.GET("/admin", func(c *gin.Context) {
		if err := sendAdminInfo(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send userInfo error"})
		}
	})
	router.GET("/admin/logout", func(c *gin.Context) {
		if err := adminLogout(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "user logout error"})
		}
	})

	// 仪表盘
	router.GET("/dashboard", func(c *gin.Context) {
		if err := sendDashBoard(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send dashboard error"})
		}
	})

	// 图片
	router.GET("/images/:name", func(c *gin.Context) {
		if err := sendImage(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send image error"})
		}
	})
	router.POST("/images", func(c *gin.Context) {
		if err := saveImage(c); err != nil {
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "save image error"})
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
	classImage := c.PostForm("class_image")

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
		ClassImage: classImage,
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
	studentImage := c.PostForm("student_image")
	studentName := c.PostForm("student_name")
	studentPassword := c.PostForm("student_password")
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
		Classes: classPointers,
		StudentImage: studentImage,
		StudentName: studentName,
		StudentPassword: fmt.Sprintf("%x", sha256.Sum256([]byte(studentPassword))),
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
		StudentName: student.StudentName,
		StudentImage: student.StudentImage,
		StudentPassword: studentPassword,
	})
	return
}

func sendStudents(c *gin.Context) (err error) {
	studentNo, byStuNo := c.GetQuery("student_no")
	classID , byClassID := c.GetQuery("class_id")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

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
		c.JSON(http.StatusOK, StudentsResponse{
			Students: []StudentResponse{{
				StudentNo: *student.StudentNo,
				FaceToken: student.FaceToken,
				ClassIDs: classUintIDs,
				StudentName: student.StudentName,
				StudentImage: student.StudentImage,
				StudentPassword: student.StudentPassword,
			}},
			Total: 1,
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

		pageCount := len(students)
		if page != "" && pageSize != "" {
			start, _ := strconv.Atoi(page)
			size, _ := strconv.Atoi(pageSize)

			if (len(students) - start * size) > size {
				students = students[(start - 1) * size: start * size]
			} else {
				students = students[(start - 1) * size:]
			}
		}

		studentsResponse := make([]StudentResponse, len(students))
		for k, v := range students {
			studentsResponse[k].FaceToken = v.FaceToken
			studentsResponse[k].StudentNo = *v.StudentNo
			studentsResponse[k].StudentImage = v.StudentImage
			studentsResponse[k].StudentName = v.StudentName
			studentsResponse[k].StudentPassword = v.StudentPassword

			classUintIDs := make([]uint, len(v.Classes))
			for k, v := range v.Classes {
				classUintIDs[k] = v.ID
			}
			studentsResponse[k].ClassIDs = classUintIDs
		}

		c.JSON(http.StatusOK, StudentsResponse{
			Students: studentsResponse,
			Total: pageCount,
		})
	} else {
		var students []Student
		students, err = getAllStudents()
		if err != nil {
			return
		}

		pageCount := len(students)
		if page != "" && pageSize != "" {
			start, _ := strconv.Atoi(page)
			size, _ := strconv.Atoi(pageSize)

			if (len(students) - start * size) > size {
				students = students[(start - 1) * size: start * size]
			} else {
				students = students[(start - 1) * size:]
			}
		}

		studentsResponse := make([]StudentResponse, len(students))
		for k, v := range students {
			studentsResponse[k].FaceToken = v.FaceToken
			studentsResponse[k].StudentNo = *v.StudentNo
			studentsResponse[k].StudentImage = v.StudentImage
			studentsResponse[k].StudentName = v.StudentName
			studentsResponse[k].StudentPassword = v.StudentPassword

			classUintIDs := make([]uint, len(v.Classes))
			for k, v := range v.Classes {
				classUintIDs[k] = v.ID
			}
			studentsResponse[k].ClassIDs = classUintIDs
		}

		c.JSON(http.StatusOK, StudentsResponse{
			Students: studentsResponse,
			Total: pageCount,
		})
	}

	return
}

func updateStudent(c *gin.Context) (err error) {
	studentNo := c.Param("no")
	studentImage := c.PostForm("student_image")
	studentName := c.PostForm("student_name")
	studentPassword := c.PostForm("student_password")
	classIDs := c.PostFormArray("[]class_ids")

	oldStudent, err := getStudent(studentNo)
	if err != nil {
		return
	}

	newStudentMap := make(map[string]interface{})

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
	classUintIDs := make([]uint, len(classes))
	for k, v := range classes {
		classUintIDs[k] = v.ID
	}
	newStudentMap["classes"] = classPointers

	if studentImage != oldStudent.StudentImage {
		newStudentMap["student_image"] = studentImage
	}
	if studentName != oldStudent.StudentName {
		newStudentMap["student_name"] = studentName
	}
	newPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(studentPassword)))
	if newPassword != oldStudent.StudentPassword {
		newStudentMap["student_password"] = newPassword
	}

	err = updateTableItem(&oldStudent, newStudentMap)
	if err != nil {
		return
	}

	newStudent, err := getStudent(studentNo)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, StudentResponse{
		StudentNo: *newStudent.StudentNo,
		FaceToken: newStudent.FaceToken,
		ClassIDs: classUintIDs,
		StudentName: newStudent.StudentName,
		StudentImage: newStudent.StudentImage,
		StudentPassword: studentPassword,
	})
	return
}

func deleteStudent(c *gin.Context) (err error) {
	studentNo := c.PostForm("student_no")
	studentNos := c.PostFormArray("student_nos")

	if studentNo != "" {
		var student *Student
		student, err = getStudent(studentNo)
		if err != nil {
			return
		}

		err = deleteTableItem(&student)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, JsonMessage{Message: "student deleted"})
	} else {
		for _, v := range studentNos {
			var student *Student
			student, err = getStudent(v)
			if err != nil {
				return
			}

			err = deleteTableItem(&student)
			if err != nil {
				return
			}

			c.JSON(http.StatusOK, JsonMessage{Message: "students deleted"})
		}
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

func adminLogin(c *gin.Context) (err error) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == password && (username == Admin || username == Developer) {
		c.SetCookie("admin_token", username, 36000, "/", "localhost", false, true)
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "Auth failed"})
		return
	}

	return
}

func sendAdminInfo(c *gin.Context) (err error) {
	cookie, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusOK, JsonMessage{Message: "Not Login"})
		return nil
	}

	c.JSON(http.StatusOK, UserInfoResp{
		User: UserInfo{Username: cookie, Permissions: cookie},
	})

	return
}

func adminLogout(c *gin.Context) (err error) {
	c.SetCookie("admin_token", "", -1, "/", "localhost", false, true)
	return
}

func sendDashBoard(c *gin.Context) (err error) {
	systemStats, err := getDeviceManagerSystemStats()
	if err != nil {
		return
	}

	cameras, devices, err := getDeviceCameraCount()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, DashBoardResp{
		SystemStats: systemStats,
		NumberCard: NumberCard{cameras, devices},
	})

	return
}

func updateStats() {
	ticker := time.NewTicker(time.Duration(UpdateStatsTime * 1000000000))
	var err error

	for range ticker.C {
		err = updateCpuMemStats()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updateCpuMemStats() (err error) {
	cpuUsed, err := cpu.Percent(time.Second, false)
	if err != nil {
		return
	}
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	memUsed := (float64(memStats.Active) / float64(memStats.Total)) * 100

	systemStats := DeviceManagerSystemStats{
		CpuUsed: getFloatPrecision(cpuUsed[0], "1"),
		MemUsed: getFloatPrecision(memUsed, "1"),
	}

	err = createTableItem(&systemStats)
	if err != nil {
		return
	}

	return
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

func saveImage(c *gin.Context) (err error) {
	file, _ := c.FormFile("file")

	fileName := fmt.Sprintf("%v%v", time.Now().UnixNano(), path.Ext(file.Filename))

	err = c.SaveUploadedFile(file, ImageFileDir + fileName)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: fileName})
	return
}

func sendImage(c *gin.Context) (err error) {
	imageName := c.Param("name")

	c.File(ImageFileDir + imageName)

	return
}