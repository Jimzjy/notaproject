package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RootAddress = "/home/pi/hostSystem/"
	ConfigFileName = RootAddress + "config.json"
	UpdateStatsTime = 30
	ImageFileDir = RootAddress + "images/"
	PDFFileDir = RootAddress + "pdf/"
	Domain = "192.168.*"
	PitchAngle = 30
	EyeClose = 50
)

var config Config
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var standUpChannels sync.Map
var standUpStatusRecord sync.Map
var client = &http.Client{
	Timeout: time.Second * 5,
}
var longTimeoutClient = &http.Client{
	Timeout: time.Second * 25,
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

//func setConfig(c *gin.Context) error {
//	var err error
//
//	var reqConfig Config
//	err = c.ShouldBindJSON(&reqConfig)
//	if err != nil {
//		return err
//	}
//
//	tmp, err := json.Marshal(reqConfig)
//	if err != nil {
//		return err
//	}
//	err = ioutil.WriteFile(ConfigFileName, tmp, 0644)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func createClass(c *gin.Context) (err error) {
	className := c.PostForm("class_name")
	classImage := c.PostForm("class_image")
	classroomNo := c.PostForm("classroom_no")
	teacherNos := c.PostFormArray("teacher_nos")

	teachers, err := getTeachers(teacherNos)
	if err != nil {
		return
	}

	teacherPointers := make([]*Teacher, len(teachers))
	for k, v := range teachers {
		if v.ID == 0 {
			continue
		}
		teacherPointers[k] = &v
	}

	classroom, err := getClassroom(classroomNo)
	if err != nil {
		return
	}

	body, err := sendPostForm(url.Values{
		"api_key": {config.ApiKey},
		"api_secret": {config.ApiSecret},
		"display_name": {className},
	}, config.CreateFaceSetUrl)

	var classResponse ClassResponse
	err = json.Unmarshal(body, &classResponse)
	if err != nil {
		return
	}

	class := Class{
		FaceSetToken: classResponse.FaceSetToken,
		ClassName: className,
		ClassImage: classImage,
		Teachers: teacherPointers,
		ClassroomNo: *classroom.ClassroomNo,
	}
	err = createTableItem(&class)
	if err != nil {
		return
	}

	//classResponse.ClassName = className
	//classResponse.ClassID = class.ID
	//c.JSON(http.StatusOK, classResponse)
	c.JSON(http.StatusOK, JsonMessage{Message: "create class successful"})
	return
}

func sendClasses(c *gin.Context) (err error) {
	classID, byID := c.GetQuery("class_id")
	classname, byName := c.GetQuery("class_name")
	studentNo, byStudentNo := c.GetQuery("student_no")
	teacherNo, byTeacherNo := c.GetQuery("teacher_no")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byID {
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

		var classesResp *ClassesResponse
		classesResp, err = newClassesResponse([]Class{*class}, page, pageSize)

		c.JSON(http.StatusOK, classesResp)
	} else if byName {
		var classes []Class
		classes, err = getClassesByName(classname)
		if err != nil {
			return
		}

		var classesResp *ClassesResponse
		classesResp, err = newClassesResponse(classes, page, pageSize)

		c.JSON(http.StatusOK, classesResp)
	} else if byStudentNo {
		var classes []Class

		classes, err = getClassesByStudentNo(studentNo)
		if err != nil {
			return
		}

		var classesResp *ClassesResponse
		classesResp, err = newClassesResponse(classes, page, pageSize)

		c.JSON(http.StatusOK, classesResp)
	} else if byTeacherNo {
		var classes []Class

		classes, err = getClassesByTeacherNo(teacherNo)
		if err != nil {
			return
		}

		var classesResp *ClassesResponse
		classesResp, err = newClassesResponse(classes, page, pageSize)

		c.JSON(http.StatusOK, classesResp)
	} else {
		var classes []Class
		classes, err = getAllClasses()
		if err != nil {
			return
		}

		var classesResp *ClassesResponse
		classesResp, err = newClassesResponse(classes, page, pageSize)

		c.JSON(http.StatusOK, classesResp)
	}

	return
}

func updateClass(c *gin.Context) (err error) {
	classID := c.PostForm("class_id")
	className := c.PostForm("class_name")
	classImage := c.PostForm("class_image")
	classroomNo := c.PostForm("classroom_no")
	studentNos := c.PostFormArray("student_nos")
	teacherNos := c.PostFormArray("teacher_nos")

	var id int
	id, err = strconv.Atoi(classID)
	if err != nil {
		return
	}

	oldClass, err := getClass(id)
	if err != nil {
		return
	}

	var newClass Class
	if oldClass.ClassImage != classImage {
		newClass.ClassImage = classImage
	}
	if oldClass.ClassName != className {
		newClass.ClassName = className
	}

	classroom, err := getClassroom(classroomNo)
	if err != nil {
		return
	}
	newClass.ClassroomNo = *classroom.ClassroomNo

	err = updateTableItem(oldClass, newClass)
	if err != nil {
		return
	}

	oldStudents, err := getStudentsByClass(id)
	if err != nil {
		return
	}
	newStudents, err := updateFaceSetByStudentNos(oldStudents, studentNos, oldClass.FaceSetToken)
	if err != nil {
		return
	}
	err = updateAssociation(oldClass, "Students", newStudents)
	if err != nil {
		return
	}

	teachers, err := getTeachers(teacherNos)
	if err != nil {
		return
	}
	teacherPointers := make([]*Teacher, len(teachers))
	for k, v := range teachers {
		if v.ID == 0 {
			continue
		}
		teacherPointers[k] = &v
	}
	err = updateAssociation(oldClass, "Teachers", teacherPointers)

	c.JSON(http.StatusOK, JsonMessage{Message: "update class successful"})
	return
}

func deleteClass(c *gin.Context) (err error) {
	classID := c.PostForm("class_id")
	classIDs := c.PostFormArray("class_ids")

	if classID != "" {
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

		_, err = sendPostForm(url.Values{
			"api_key": {config.ApiKey},
			"api_secret": {config.ApiSecret},
			"faceset_token": {class.FaceSetToken},
			"check_empty": {"0"},
		}, config.DeleteFaceSetUrl)
		if err != nil {
			return
		}

		err = deleteTableItem(class)
		if err != nil {
			return
		}
	} else if len(classIDs) > 0 {
		for _, v := range classIDs {
			var id int
			id, err = strconv.Atoi(v)
			if err != nil {
				return
			}

			var class *Class
			class, err = getClass(id)
			if err != nil {
				return
			}

			_, err = sendPostForm(url.Values{
				"api_key": {config.ApiKey},
				"api_secret": {config.ApiSecret},
				"faceset_token": {class.FaceSetToken},
				"check_empty": {"0"},
			}, config.DeleteFaceSetUrl)
			if err != nil {
				return
			}

			err = deleteTableItem(class)
			if err != nil {
				return
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "delete class successful"})
	return
}

//func detectFace(c *gin.Context) error {
//	var err error
//
//	fileHeader, err := c.FormFile("image_file")
//	if err != nil {
//		return err
//	}
//	src, err := fileHeader.Open()
//	if err != nil {
//		return err
//	}
//	defer src.Close()
//	data, err := ioutil.ReadAll(src)
//	if err != nil {
//		return err
//	}
//
//	params := map[string]string{
//		"api_key": config.ApiKey,
//		"api_secret": config.ApiSecret,
//	}
//	body, err := fileUploadRequest(config.DetectFaceUrl, params,
//		"image_file", data, fileHeader.Filename)
//	if err != nil {
//		return err
//	}
//
//	var faceRectTokens FaceRectTokens
//	err = json.Unmarshal(body, &faceRectTokens)
//	if err != nil {
//		return err
//	}
//	if len(faceRectTokens.Faces) < 1 {
//		c.JSON(http.StatusBadRequest, "no person in image")
//		return nil
//	}
//
//	stuFace := FaceNoToken{
//		FaceToken: faceRectTokens.Faces[0].FaceToken,
//		StudentNo: strings.TrimSuffix(fileHeader.Filename, path.Ext(fileHeader.Filename)),
//	}
//	c.JSON(http.StatusOK, stuFace)
//
//	return nil
//}

func updateClassroomStats(c *gin.Context) (err error) {
	var stats Stats
	err = c.ShouldBindJSON(&stats)
	if err != nil {
		return
	}

	devicePath := strings.Split(c.Request.RemoteAddr, ":")[0]
	device, err := getDeviceByPath(devicePath)
	if err != nil {
		return
	}

	err = createTableItem(&DeviceStatsTable{
		UpdateTime: stats.UpdateTime.Unix(),
		CpuUsed: stats.SystemStats.CpuUsed,
		MemUsed: stats.SystemStats.MemUsed,
		DeviceID: device.ID,
	})
	if err != nil {
		return
	}

	var classroom *Classroom
	for _, classroomStats := range stats.Classrooms {
		classroom, err = getClassroom(classroomStats.ClassroomNo)
		if err != nil {
			return
		}

		var persons []byte
		persons, err = json.Marshal(classroomStats.Persons)
		if err != nil {
			return
		}

		err = createTableItem(&ClassroomStatsTable{
			UpdateTime: stats.UpdateTime.Unix(),
			PersonCount: classroomStats.PersonCount,
			Persons: string(persons),
			ClassroomNo: *classroom.ClassroomNo,
			GlobalWidth: classroomStats.GlobalWidth,
			GlobalHeight: classroomStats.GlobalHeight,
		})

		if err != nil {
			return
		}
	}

	return
}

func sendClassroomStats(c *gin.Context) (err error) {
	classroomNo, isExist := c.GetQuery("classroom_no")
	if !isExist {
		return fmt.Errorf("no classroom_no")
	}

	classroomStatsItem, err := getClassroomStatsItem(classroomNo)
	if err != nil {
		return
	}

	var persons []FaceRectangle
	err = json.Unmarshal([]byte(classroomStatsItem.Persons), &persons)
	if err != nil {
		return
	}

	stats := SingleClassroomStats{
		UpdateTime: classroomStatsItem.UpdateTime,
		ClassroomStats: ClassroomStats{
			ClassroomNo: classroomStatsItem.ClassroomNo,
			PersonCount: classroomStatsItem.PersonCount,
			GlobalWidth: classroomStatsItem.GlobalWidth,
			GlobalHeight: classroomStatsItem.GlobalHeight,
			Persons: persons,
		},
	}

	c.JSON(http.StatusOK, stats)
	return
}

func sendDeviceStats(c *gin.Context) (err error) {
	deviceID, isExist := c.GetQuery("device_id")
	if !isExist {
		return fmt.Errorf("no device_id")
	}

	id, err := strconv.Atoi(deviceID)
	if err != nil {
		return
	}

	deviceStatsItem, err := getDeviceStats(id)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, DeviceStats{
		UpdateTime: deviceStatsItem[0].UpdateTime,
		DeviceStats: deviceStatsItem,
	})
	return
}

func createDevice(c *gin.Context) (err error) {
	devicePath := c.PostForm("device_path")
	devicePort := c.PostForm("device_port")
	cameraIDs := c.PostFormArray("camera_ids")

	_cameraIDs, err := stringArrayToIntArray(cameraIDs)
	cameras, err := getCameras(_cameraIDs)
	if err != nil {
		return
	}

	cameraPointers := make([]*Camera, len(cameras))
	for k, v := range cameras {
		if v.ID == 0 {
			continue
		}
		cameraPointers[k] = &v
	}

	device := Device{
		DevicePath: devicePath,
		DevicePort: devicePort,
		Cameras: cameraPointers,
	}
	err = createTableItem(&device)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "create device successful"})
	return
}

func sendDevices(c *gin.Context) (err error) {
	deviceID, byID := c.GetQuery("device_id")
	cameraID, byCamera := c.GetQuery("camera_id")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byID {
		var id int
		id, err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}

		var device *Device
		device, err = getDevice(id)

		var devicesResponse *DevicesResponse
		devicesResponse, err = newDevicesResponse([]Device{*device}, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, devicesResponse)
	} else if byCamera {
		var id int
		id, err = strconv.Atoi(cameraID)
		if err != nil {
			return
		}

		var devices []Device
		devices, err = getDevicesByCamera(id)
		if err != nil {
			return
		}

		var devicesResponse *DevicesResponse
		devicesResponse, err = newDevicesResponse(devices, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, devicesResponse)
	} else {
		var devices []Device
		devices, err = getAllDevices()
		if err != nil {
			return
		}

		var devicesResponse *DevicesResponse
		devicesResponse, err = newDevicesResponse(devices, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, devicesResponse)
	}

	return
}

func updateDevices(c *gin.Context) (err error) {
	deviceID := c.PostForm("device_id")
	devicePath := c.PostForm("device_path")
	devicePort := c.PostForm("device_port")
	cameraIDs := c.PostFormArray("camera_ids")

	var id int
	id, err = strconv.Atoi(deviceID)
	if err != nil {
		return
	}

	oldDevice, err := getDevice(id)
	if err != nil {
		return
	}

	var newDevice Device
	if oldDevice.DevicePath != devicePath {
		newDevice.DevicePath = devicePath
	}
	if oldDevice.DevicePort != devicePort {
		newDevice.DevicePort = devicePort
	}

	err = updateTableItem(oldDevice, newDevice)
	if err != nil {
		return
	}

	_cameraIDs, err := stringArrayToIntArray(cameraIDs)
	if err != nil {
		return
	}
	newCameras, err := getCameras(_cameraIDs)
	if err != nil {
		return
	}
	err = updateAssociation(oldDevice, "Cameras", newCameras)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "update device successful"})
	return
}

func deleteDevices(c *gin.Context) (err error) {
	deviceID := c.PostForm("device_id")
	deviceIDs := c.PostFormArray("device_ids")

	if deviceID != "" {
		var id int
		id, err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}

		var device *Device
		device, err = getDevice(id)

		err = deleteTableItem(device)
		if err != nil {
			return
		}
	} else if len(deviceIDs) > 0 {
		var ids []int
		ids, err = stringArrayToIntArray(deviceIDs)
		if err != nil {
			return
		}

		err = deleteTableItems(Device{}, "id in (?)", ids)
		if err != nil {
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "delete device successful"})
	return
}

func createStudent(c *gin.Context) (err error) {
	studentNo := c.PostForm("student_no")
	studentImage := c.PostForm("student_image")
	studentName := c.PostForm("student_name")
	studentPassword := c.PostForm("student_password")

	data, err := ioutil.ReadFile(fmt.Sprintf(ImageFileDir + "%v", studentImage))
	if err != nil {
		return
	}

	params := map[string]string{
		"api_key": config.ApiKey,
		"api_secret": config.ApiSecret,
	}
	body, err := fileUploadRequest(config.DetectFaceUrl, params,
		"image_file", data, studentImage)
	if err != nil {
		return
	}

	var faceRectTokens FaceRectTokens
	err = json.Unmarshal(body, &faceRectTokens)
	if err != nil {
		return err
	}
	if len(faceRectTokens.Faces) < 1 {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "get face from student image failed"})
		return nil
	}

	student := Student{
		StudentNo: &studentNo,
		StudentImage: studentImage,
		StudentName: studentName,
		StudentPassword: fmt.Sprintf("%x", sha256.Sum256([]byte(studentPassword))),
		FaceToken: faceRectTokens.Faces[0].FaceToken,
	}

	err = createTableItem(&student)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "create student successful"})
	return
}

func sendStudents(c *gin.Context) (err error) {
	studentNo, byStuNo := c.GetQuery("student_no")
	classID, byClassID := c.GetQuery("class_id")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byStuNo {
		var student *Student
		student, err = getStudent(studentNo)
		if err != nil {
			return
		}

		var studentsResp *StudentsResponse
		studentsResp, err = newStudentsResponse([]Student{*student}, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, studentsResp)
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

		var studentsResp *StudentsResponse
		studentsResp, err = newStudentsResponse(students, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, studentsResp)
	} else {
		var students []Student
		students, err = getAllStudents()
		if err != nil {
			return
		}

		var studentsResp *StudentsResponse
		studentsResp, err = newStudentsResponse(students, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, studentsResp)
	}

	return
}

func updateStudent(c *gin.Context) (err error) {
	studentNo := c.PostForm("student_no")
	studentImage := c.PostForm("student_image")
	studentName := c.PostForm("student_name")
	studentPassword := c.PostForm("student_password")

	oldStudent, err := getStudent(studentNo)
	if err != nil {
		return
	}

	var newStudent Student

	if studentImage != oldStudent.StudentImage {
		newStudent.StudentImage = studentImage
	}
	if studentName != oldStudent.StudentName {
		newStudent.StudentName = studentName
	}

	if len(studentPassword) > 0 {
		newPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(studentPassword)))
		if newPassword != oldStudent.StudentPassword {
			newStudent.StudentPassword = newPassword
		}
	}

	err = updateTableItem(oldStudent, newStudent)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "update student successful"})
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
	} else if len(studentNos) > 0 {
		err = deleteTableItems(Student{}, "student_no in (?)", studentNos)
		if err != nil {
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "student deleted"})
	return
}

func createTeacher(c *gin.Context) (err error) {
	teacherNo := c.PostForm("teacher_no")
	teacherImage := c.PostForm("teacher_image")
	teacherName := c.PostForm("teacher_name")
	teacherPassword := c.PostForm("teacher_password")

	teacher := Teacher{
		TeacherNo: &teacherNo,
		TeacherImage: teacherImage,
		TeacherName: teacherName,
		TeacherPassword: fmt.Sprintf("%x", sha256.Sum256([]byte(teacherPassword))),
	}

	err = createTableItem(&teacher)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "create teacher successful"})
	return
}

func sendTeachers(c *gin.Context) (err error) {
	teacherNo, byNo := c.GetQuery("teacher_no")
	classID, byClass := c.GetQuery("class_id")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byNo {
		var teacher *Teacher
		teacher, err = getTeacher(teacherNo)
		if err != nil {
			return
		}

		var teachersResp *TeachersResponse
		teachersResp, err = newTeacherResponse([]Teacher{*teacher}, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, teachersResp)
	} else if byClass {
		var id int
		id ,err = strconv.Atoi(classID)
		if err != nil {
			return
		}

		var teachers []Teacher
		teachers, err = getTeachersByClass(id)
		if err != nil {
			return
		}

		var teachersResp *TeachersResponse
		teachersResp, err = newTeacherResponse(teachers, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, teachersResp)
	} else {
		var teachers []Teacher
		teachers, err = getAllTeachers()
		if err != nil {
			return
		}

		var teachersResp *TeachersResponse
		teachersResp, err = newTeacherResponse(teachers, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, teachersResp)
	}
	return
}

func updateTeacher(c *gin.Context) (err error) {
	teacherNo := c.PostForm("teacher_no")
	teacherImage := c.PostForm("teacher_image")
	teacherName := c.PostForm("teacher_name")
	teacherPassword := c.PostForm("teacher_password")

	oldTeacher, err := getTeacher(teacherNo)
	if err != nil {
		return
	}

	var newTeacher Teacher

	newTeacher.TeacherImage = teacherImage
	newTeacher.TeacherName = teacherName

	if len(teacherPassword) > 0 {
		newPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(teacherPassword)))
		if newPassword != oldTeacher.TeacherPassword {
			newTeacher.TeacherPassword = newPassword
		}
	}

	err = updateTableItem(oldTeacher, newTeacher)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "update teacher successful"})
	return
}

func deleteTeacher(c *gin.Context) (err error) {
	teacherNo := c.PostForm("teacher_no")
	teacherNos := c.PostForm("teacher_nos")

	if teacherNo != "" {
		var teacher *Teacher
		teacher, err = getTeacher(teacherNo)
		if err != nil {
			return
		}

		err = deleteTableItem(teacher)
		if err != nil {
			return
		}
	} else if len(teacherNos) > 0 {
		err = deleteTableItems(Teacher{}, "teacher_no in (?)", teacherNos)
		if err != nil {
			return 
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "delete teacher successful"})
	return
}

func createCamera(c *gin.Context) (err error) {
	camStreamPath := c.PostForm("cam_stream_path")
	camONVIFPath := c.PostForm("cam_onvif_path")
	camAuthName := c.PostForm("cam_auth_name")
	camAuthPassword := c.PostForm("cam_auth_password")
	deviceID := c.PostForm("device_id")
	classroomNo := c.PostForm("classroom_no")

	var devices []*Device
	var classrooms []*Classroom
	var id int
	if deviceID != "" {
		id ,err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}

		var device *Device
		device, err = getDevice(id)
		if err != nil {
			return
		}

		if device.ID != 0 {
			devices = append(devices, device)
		}
	}
	if classroomNo != "" {
		var classroom *Classroom
		classroom, err = getClassroom(classroomNo)
		if err != nil {
			return
		}

		if classroom.ID != 0 {
			classrooms = append(classrooms, classroom)
		}
	}

	camera := Camera{
		CamStreamPath: camStreamPath,
		CamONVIFPath: camONVIFPath,
		CamAuthName: camAuthName,
		CamAuthPassword: camAuthPassword,
		Devices: devices,
		Classrooms: classrooms,
	}

	err = createTableItem(&camera)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "create camera successful"})
	return
}

func sendCameras(c *gin.Context) (err error) {
	cameraID, byID := c.GetQuery("camera_id")
	deviceID, byDevice := c.GetQuery("device_id")
	classroomNo, byClassroom := c.GetQuery("classroom_no")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byID {
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

		var camerasResponse *CamerasResponse
		camerasResponse, err = newCamerasResponse([]Camera{*camera}, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, camerasResponse)
	} else if byDevice {
		var cameras []Camera

		var id int
		id ,err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}
		cameras, err = getCamerasByDevice(id)
		if err != nil {
			return
		}

		var camerasResponse *CamerasResponse
		camerasResponse, err = newCamerasResponse(cameras, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, camerasResponse)
	} else if byClassroom {
		var cameras []Camera

		cameras, err = getCamerasByClassroom(classroomNo)
		if err != nil {
			return
		}

		var camerasResponse *CamerasResponse
		camerasResponse, err = newCamerasResponse(cameras, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, camerasResponse)
	} else {
		var cameras []Camera

		cameras, err = getAllCameras()
		if err != nil {
			return
		}

		var camerasResponse *CamerasResponse
		camerasResponse, err = newCamerasResponse(cameras, page, pageSize)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, camerasResponse)
	}
	return
}

func updateCameras(c *gin.Context) (err error) {
	cameraID := c.PostForm("camera_id")
	camStreamPath := c.PostForm("cam_stream_path")
	camONVIFPath := c.PostForm("cam_onvif_path")
	camAuthName := c.PostForm("cam_auth_name")
	camAuthPassword := c.PostForm("cam_auth_password")
	deviceID := c.PostForm("device_id")
	classroomNo := c.PostForm("classroom_no")

	var id int
	id, err = strconv.Atoi(cameraID)
	if err != nil {
		return
	}

	oldCamera, err := getCamera(id)
	if err != nil {
		return
	}

	var newCamera Camera
	newCamera.CamStreamPath = camStreamPath
	newCamera.CamONVIFPath = camONVIFPath
	newCamera.CamAuthName = camAuthName
	newCamera.CamAuthPassword = camAuthPassword

	err = updateTableItem(oldCamera, newCamera)
	if err != nil {
		return
	}

	var devices []*Device
	var classrooms []*Classroom
	if deviceID != "" {
		id ,err = strconv.Atoi(deviceID)
		if err != nil {
			return
		}

		var device *Device
		device, err = getDevice(id)
		if err != nil {
			return
		}

		if device.ID != 0 {
			devices = append(devices, device)
		}
	}
	if classroomNo != "" {
		var classroom *Classroom
		classroom, err = getClassroom(classroomNo)
		if err != nil {
			return
		}

		if classroom.ID != 0 {
			classrooms = append(classrooms, classroom)
		}
	}

	err = updateAssociation(oldCamera, "Devices", devices)
	if err != nil {
		return
	}
	err = updateAssociation(oldCamera, "Classrooms", classrooms)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "update camera successful"})
	return
}

func deleteCameras(c *gin.Context) (err error) {
	cameraID := c.PostForm("camera_id")
	cameraIDs := c.PostFormArray("camera_ids")

	if cameraID != "" {
		var id int
		id ,err = strconv.Atoi(cameraID)
		if err != nil {
			return
		}

		var camera *Camera
		camera, err = getCamera(id)
		if err != nil {
			return
		}

		err = deleteTableItem(camera)
		if err != nil {
			return
		}
	} else if len(cameraIDs) > 0 {
		var ids []int
		ids, err = stringArrayToIntArray(cameraIDs)
		if err != nil {
			return
		}

		err = deleteTableItems(Camera{}, "id in (?)", ids)
		if err != nil {
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "delete camera successful"})
	return
}

func createClassroom(c *gin.Context) (err error) {
	classroomNo := c.PostForm("classroom_no")
	cameraID := c.PostForm("camera_id")

	var cameras []*Camera
	if cameraID != "" {
		var id int
		id ,err = strconv.Atoi(cameraID)
		if err != nil {
			return
		}

		var camera *Camera
		camera, err = getCamera(id)
		if err != nil {
			return
		}

		if camera.ID != 0 {
			cameras = append(cameras, camera)
		}
	}

	classroom := Classroom{
		ClassroomNo: &classroomNo,
		Cameras: cameras,
	}

	err = createTableItem(&classroom)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "create classroom successful"})
	return
}

func sendClassrooms(c *gin.Context) (err error) {
	classroomNo, byNo := c.GetQuery("classroom_no")
	cameraID, byCamera := c.GetQuery("camera_id")

	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if byNo {
		var classroom *Classroom

		classroom, err = getClassroom(classroomNo)
		if err != nil {
			return
		}

		var classroomsResponse *ClassroomsResponse
		classroomsResponse, err = newClassroomsResponse([]Classroom{*classroom}, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, classroomsResponse)
	} else if byCamera {
		var id int
		id, err = strconv.Atoi(cameraID)
		if err != nil {
			return
		}

		var classrooms []Classroom
		classrooms, err = getClassroomsByCamera(id)
		if err != nil {
			return
		}

		var classroomsResponse *ClassroomsResponse
		classroomsResponse, err = newClassroomsResponse(classrooms, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, classroomsResponse)
	} else {
		var classrooms []Classroom

		classrooms, err = getAllClassrooms()
		if err != nil {
			return
		}

		var classroomsResponse *ClassroomsResponse
		classroomsResponse, err = newClassroomsResponse(classrooms, page, pageSize)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, classroomsResponse)
	}
	return
}

func updateClassrooms(c *gin.Context) (err error) {
	classroomNo := c.PostForm("classroom_no")
	cameraID := c.PostForm("camera_id")

	oldClassroom, err := getClassroom(classroomNo)
	if err != nil {
		return
	}

	var id int
	id, err = strconv.Atoi(cameraID)
	if err != nil {
		return
	}
	camera, err := getCamera(id)
	if err != nil {
		return
	}

	err = updateAssociation(oldClassroom, "Cameras", []Camera{*camera})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "update classroom successful"})
	return
}

func deleteClassrooms(c *gin.Context) (err error) {
	classroomNo := c.PostForm("classroom_no")
	classroomNos := c.PostFormArray("classroom_nos")

	if classroomNo != "" {
		var classroom *Classroom
		classroom, err = getClassroom(classroomNo)
		if err != nil {
			return
		}

		err = deleteTableItem(classroom)
		if err != nil {
			return
		}
	} else if len(classroomNos) > 0 {
		err = deleteTableItems(Classroom{}, "classroom_no in (?)", classroomNos)
		if err != nil {
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no params provide"})
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: "delete classroom successful"})
	return
}

func sendFaceCountRecord(c *gin.Context) (err error) {
	recordID := c.Query("record_id")

	id, err := strconv.Atoi(recordID)
	if err != nil {
		return
	}
	faceCountRecord, err := getFaceCountRecordByID(id)
	if err != nil {
		return
	}

	var faceRectTokens []FaceRectToken
	err = json.Unmarshal([]byte(faceCountRecord.FaceRectTokens), &faceRectTokens)
	if err != nil {
		return
	}

	students, err := getStudentsByClass(faceCountRecord.ClassID)
	if err != nil {
		return
	}
	var studentInClass []string
	for _, v1 := range faceRectTokens {
		for _, v2 := range students {
			if checkIfInStringSlice(studentInClass, v1.FaceToken) {
				break
			}

			if v1.FaceToken == *v2.StudentNo {
				studentInClass = append(studentInClass, v1.FaceToken)
				break
			}
		}
	}
	var studentNotIn []string
	for _, v := range students {
		if !checkIfInStringSlice(studentInClass, *v.StudentNo) && !checkIfInStringSlice(studentNotIn, *v.StudentNo) {
			studentNotIn = append(studentNotIn, *v.StudentNo)
		}
	}

	c.JSON(http.StatusOK, FaceCountRecordResponse{
		FaceRectTokens: faceRectTokens,
		StudentCount: len(students),
		StudentInClassCount: len(studentInClass),
		StudentNotInClass: studentNotIn,
	})
	return
}

func faceCount(c *gin.Context) (err error) {
	classID := c.Query("class_id")
	id, err := strconv.Atoi(classID)
	if err != nil {
		return
	}
	class, err := getClass(id)
	if err != nil {
		return
	}

	camera, err := getCamerasByClassroom(class.ClassroomNo)
	if err != nil {
		return
	}
	if len(camera) < 1 {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no camera in classroom"})
		return
	}

	device, err := getDevicesByCamera(int(camera[0].ID))
	if err != nil {
		return
	}
	if len(device) < 1 {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no device"})
		return
	}

	chanResp := make(chan string)
	chanRequest := make(chan string)
	chanPersonData := make(chan []byte)
	var faceRecords []FaceRectToken

	go handleFaceCountResp(c, chanResp, chanPersonData)
	go handleFaceCountRequest(camera[0].CamStreamPath, class.FaceSetToken,
		device[0].DevicePath, device[0].DevicePort, chanRequest, chanPersonData, chanResp, &faceRecords)

	<- chanRequest
	<- chanResp

	var err2 error
	faceRecordsBytes, err2 := json.Marshal(faceRecords)
	if err2 != nil {
		log.Println(err2)
		return
	}

	faceCountRecord := FaceCountRecord{
		FaceRectTokens: string(faceRecordsBytes),
		ClassID: id,
	}
	err2 = createTableItem(&faceCountRecord)
	if err2 != nil {
		log.Println(err2)
		return
	}
	return
}

func handleFaceCountResp(c *gin.Context, chanResp chan string, chanPersonData chan []byte) {
	defer close(chanResp)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("resp: ", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("read from user:", err)
				break
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case personData, ok := <-chanPersonData:
			if ok {
				err = conn.WriteMessage(websocket.TextMessage, personData)
				if err != nil {
					log.Println("write to user:", err)
					return
				}
			} else {
				err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write to user", err)
					return
				}
				return
			}
		}
	}
}

func handleFaceCountRequest(camSteamPath, faceSetToken, devicePath, devicePort string, chanRequest chan string, chanPersonData chan []byte, chanResp chan string, faceRecords *[]FaceRectToken) {
	defer close(chanRequest)
	defer close(chanPersonData)

	urlString := fmt.Sprintf("ws://%v:%v/face_search?camStreamPath=%v&faceSetToken=%v",
		devicePath, devicePort, camSteamPath, faceSetToken)

	conn, _, err := websocket.DefaultDialer.Dial(urlString, nil)
	if err != nil {
		log.Println("request: ", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})

	var personData PersonData
	var faceTokenRecords []string
	go func() {
		defer close(done)

		for {
			select {
			case <-chanResp:
				return
			default:
				_, data, err := conn.ReadMessage()
				if err != nil {
					log.Println("read from embedded:", err)
					return
				}

				err = json.Unmarshal(data, &personData)
				if err != nil {
					log.Println(err)
					continue
				}

				if personData.Face.FaceToken == "" {
					chanPersonData <- data
					continue
				}
				for _, v := range faceTokenRecords {
					if v == personData.Face.FaceToken {
						chanPersonData <- data
						continue
					}
				}
				faceTokenRecords = append(faceTokenRecords, personData.Face.FaceToken)

				var student *Student
				student, err = getStudentByFaceToken(personData.Face.FaceToken)
				if err != nil {
					log.Println(err)
					chanPersonData <- data
					continue
				}

				personData.Face.FaceToken = *student.StudentNo
				*faceRecords = append(*faceRecords, FaceRectToken{FaceRectangle: personData.Face.FaceRectangle, FaceToken: personData.Face.FaceToken})

				data2, err := json.Marshal(personData)
				if err != nil {
					log.Println(err)
					chanPersonData <- data
					continue
				}

				chanPersonData <- data2
			}
		}
	}()

	<- done
}

func standUp(c *gin.Context) (err error) {
	classID := c.Query("class_id")
	teacherNo := c.Query("teacher_no")
	pdfUrl := c.Query("pdf_url")
	pdfNumPages := c.Query("pdf_num_pages")

	if teacherNo == "" {
		err = fmt.Errorf("no teacher_no")
		return
	}

	id, err := strconv.Atoi(classID)
	if err != nil {
		return
	}

	if pdfUrl == "" || pdfNumPages == "" {
		err = fmt.Errorf("no pdf params")
		return
	}

	_pdfPages, err := strconv.Atoi(pdfNumPages)
	if err != nil {
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	//standUpChannels = append(standUpChannels, make(chan StandUpPacket, 2))
	standUpRecordID := fmt.Sprintf("%v", time.Now().UnixNano())
	readChannelIndex := standUpRecordID + "0"
	standUpStatusRecord.Store(standUpRecordID, &StandUpStatusRecord{
		ID: standUpRecordID,
		ClassID: id,
		TeacherNo: teacherNo,
		WReadMWriteIndex: readChannelIndex,
		//WWriteMReadIndex: "",
	})
	readChannel := make(chan StandUpPacket, 2)
	var writeChannel *chan StandUpPacket = nil
	standUpChannels.Store(readChannelIndex, &readChannel)

	//standUpStatusTable := StandUpStatusTable{
	//	ClassID: id,
	//	TeacherNo: teacherNo,
	//	WReadMWriteIndex: readChannelIndex,
	//	WWriteMReadIndex: writeChannelIndex,
	//}
	//err = createTableItem(&standUpStatusTable)
	//if err != nil {
	//	log.Println(err)
	//}

	var currentPDFPage int
	var faceRectNos []FaceRectToken
	var faceCountFinish = false
	var faceCountRecordID = 0
	done := make(chan string)
	go func() {
		defer close(done)

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("read from teacher web client:", err)
				break
			}

			var standUpPacket StandUpPacket
			err = json.Unmarshal(data, &standUpPacket)
			if err != nil {
				log.Println(err)
				continue
			}

			if standUpPacket.FaceCountClose {
				var _faceCountRecord *FaceCountRecord
				_faceCountRecord, err = getLastFaceCountRecord(id)
				if err != nil {
					log.Println(err)
					continue
				}

				err = json.Unmarshal([]byte(_faceCountRecord.FaceRectTokens), &faceRectNos)
				if err != nil {
					log.Println(err)
					continue
				}

				faceCountRecordID = int(_faceCountRecord.ID)
				standUpPacket.FaceCountRecordID = faceCountRecordID
				faceCountFinish = true
			}

			if standUpPacket.CurrentPDFPage > 0 {
				currentPDFPage = standUpPacket.CurrentPDFPage
			}

			if writeChannel != nil {
				select {
				case *writeChannel <- standUpPacket:
				default:
				}
			}
		}
	}()

	class, err := getClass(id)
	if err != nil {
		log.Println(err)
		return
	}
	students, err := getStudentsByClass(id)
	if err != nil {
		log.Println(err)
		return
	}
	studentWarningRecordList := make([]StudentWarningRecord, len(students))
	for i := 0; i < len(studentWarningRecordList); i++ {
		studentWarningRecordList[i].StudentNo = *students[i].StudentNo
	}

	studentsStatusWithPage := make([]StudentStatusWithPage, _pdfPages)
	for i := 0; i < _pdfPages; i++ {
		studentsStatusWithPage[i].PDFPage = i + 1
	}
	go func() {
		camera, err := getCamerasByClassroom(class.ClassroomNo)
		if err != nil {
			log.Println(err)
			return
		}
		if len(camera) <= 0 {
			err = fmt.Errorf("no camera in classroom %v", class.ClassroomNo)
			log.Println(err)
			return
		}

		device, err := getDevicesByCamera(int(camera[0].ID))
		if err != nil {
			log.Println(err)
			return
		}
		if len(device) <= 0 {
			err = fmt.Errorf("no device with camera %v", int(camera[0].ID))
			log.Println(err)
			return
		}

		ticker := time.NewTicker(time.Duration(config.DetectInterval * 1e9))

		for range ticker.C {
			select {
			case _, ok := <-done:
				if !ok {
					return
				}
			default:
				if !faceCountFinish {
					continue
				}
			}

			var _studentsStatus []StudentStatus
			devicePath := fmt.Sprintf("http://%v:%v/person_status", device[0].DevicePath, device[0].DevicePort)
			_studentsStatus, err = getStudentsStatus(camera[0].CamStreamPath, devicePath, faceRectNos)
			if err != nil {
				log.Println(err)
			}

			if writeChannel != nil {
				studentWarningList := ""
				for _, v := range _studentsStatus {
					headPose := v.Attributes.HeadPose
					eyesStatusL := v.Attributes.EyesStatus.LeftEyeStatus
					eyesStatusR := v.Attributes.EyesStatus.RightEyeStatus

					if math.Abs(headPose.PitchAngle) > PitchAngle {
						for i := 0; i < len(studentWarningRecordList); i++ {
							if studentWarningRecordList[i].StudentNo == v.StudentNo {
								if studentWarningRecordList[i].LastWarning {
									studentWarningList += v.StudentNo + ", "
									studentWarningRecordList[i].LastWarning = false
									studentWarningRecordList[i].Warning += 1
								} else {
									studentWarningRecordList[i].LastWarning = true
								}
								break
							}
						}
						continue
					}

					if eyesStatusL.NoGlassEyeClose > EyeClose ||
						eyesStatusL.NormalGlassEyeClose > EyeClose ||
						eyesStatusR.NoGlassEyeClose > EyeClose ||
						eyesStatusR.NormalGlassEyeClose > EyeClose {

						for i := 0; i < len(studentWarningRecordList); i++ {
							if studentWarningRecordList[i].StudentNo == v.StudentNo {
								if studentWarningRecordList[i].LastWarning {
									studentWarningList += v.StudentNo + "(闭眼), "
									studentWarningRecordList[i].LastWarning = false
									studentWarningRecordList[i].Warning += 1
								} else {
									studentWarningRecordList[i].LastWarning = true
								}
								break
							}
						}
						continue
					}
				}

				if len(studentWarningList) > 0 {
					_standUpPacket := StandUpPacket{
						StudentWarningList: strings.TrimSuffix(studentWarningList, ", "),
						StudentWarningRecordList: studentWarningToIntList(studentWarningRecordList),
					}
					select {
					case *writeChannel <- _standUpPacket:
					default:
					}
				}
			}

			if currentPDFPage > 0 {
				studentsStatusWithPage[currentPDFPage-1].StudentsStatus = append(studentsStatusWithPage[currentPDFPage-1].StudentsStatus, _studentsStatus...)
			}
		}
	}()

	var startData []byte
	startData, err = json.Marshal(StandUpPacket{
		WReadMWriteIndex: readChannelIndex,
		//WWriteMReadIndex: writeChannelIndex,
	})
	if err != nil {
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, startData)
	if err != nil {
		log.Println("write to teacher web client", err)
		return
	}

	for {
		select {
		case _, ok := <-done:
			if !ok {
				//err = deleteTableItem(&standUpStatusTable)
				//if err != nil {
				//	log.Println(err)
				//}
				standUpStatusRecord.Delete(standUpRecordID)
				close(readChannel)
				standUpChannels.Delete(readChannelIndex)
				if writeChannel != nil {
					closePack := StandUpPacket{
						SayGoodbye: true,
					}
					*writeChannel <- closePack
				}

				var finalStudentStatus, finalStudentWarningRecordList []byte
				//finalFaceCountResult, err = json.Marshal(faceRectNos)
				//if err != nil {
				//	log.Println(err)
				//}
				finalStudentStatus, err = json.Marshal(studentsStatusWithPage)
				if err != nil {
					log.Println(err)
				}

				finalStudentWarningRecordList, err = json.Marshal(studentWarningRecordList)
				if err != nil {
					log.Println(err)
				}

				err = createTableItem(&StudentStatusTable{
					ClassID: id,
					ClassName: class.ClassName,
					PDF: pdfUrl,
					PDFPageCount: _pdfPages,
					TeacherNo: teacherNo,
					FaceCountRecordID: faceCountRecordID,
					StudentStatus: string(finalStudentStatus),
					StudentWarningRecordList: string(finalStudentWarningRecordList),
				})
				if err != nil {
					log.Println(err)
				}
				return
			}
		case standUpPacket, ok := <-readChannel:
			if ok {
				if len(standUpPacket.WWriteMReadIndex) > 0 {
					//writeChannelIndex = standUpPacket.WWriteMReadIndex
					_v, existed := standUpChannels.Load(standUpPacket.WWriteMReadIndex)
					if existed {
						if _writeChannel, _typeOK := _v.(*chan StandUpPacket); _typeOK {
							writeChannel = _writeChannel
						} else {
							log.Println("WWriteMReadIndex in standUpPacket not type of chan StandUpPacket!")
						}
					} else {
						log.Println("WWriteMReadIndex not existed!")
					}
				}

				if standUpPacket.RequestStartPacket {
					startPacket := StandUpPacket{
						PDFUrl: pdfUrl,
						CurrentPDFPage: currentPDFPage,
						StudentWarningRecordList: studentWarningToIntList(studentWarningRecordList),
					}

					if faceCountRecordID != 0 {
						startPacket.FaceCountRecordID = faceCountRecordID
						startPacket.FaceCountClose = true
					}

					*writeChannel <- startPacket
				}

				var data []byte
				data, err = json.Marshal(standUpPacket)
				if err != nil {
					continue
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("write to teacher web client", err)
					return
				}
			}
		}
	}
}

func standUpMobile(c *gin.Context) (err error) {
	classID := c.Query("class_id")
	//writeChannelIndex := c.Query("write_channel_index")

	if classID == "" {
		err = fmt.Errorf("no class id")
		return
	}

	var id int
	id, err = strconv.Atoi(classID)
	if err != nil {
		return
	}
	var standUpStatus *StandUpStatusRecord
	standUpStatus, err = getStandUpStatusRecord(id, "")
	if err != nil {
		err = fmt.Errorf("not stand up")
		return
	}

	readChannelIndex := standUpStatus.WWriteMReadIndex
	writeChannelIndex := standUpStatus.WReadMWriteIndex
	var readChannel *chan StandUpPacket
	var writeChannel *chan StandUpPacket

	_v, existed := standUpChannels.Load(standUpStatus.WReadMWriteIndex)
	if existed {
		if _writeChannel, _typeOK := _v.(*chan StandUpPacket); _typeOK {
			writeChannel = _writeChannel
		} else {
			log.Println("standUpMobile WReadMWriteIndex in standUpPacket not type of chan StandUpPacket!")
		}
	} else {
		log.Println("standUpMobile WReadMWriteIndex not existed!")
	}

	if !(len(readChannelIndex) > 0) {
		_readChannel := make(chan StandUpPacket, 2)
		readChannel = &_readChannel
		readChannelIndex = writeChannelIndex + "0"
		standUpChannels.Store(readChannelIndex, &readChannel)

		standUpStatus.WWriteMReadIndex = readChannelIndex

		*writeChannel <- StandUpPacket{
			WWriteMReadIndex: readChannelIndex,
		}
	} else {
		_v, existed := standUpChannels.Load(standUpStatus.WWriteMReadIndex)
		if existed {
			if _readChannel, _typeOK := _v.(*chan StandUpPacket); _typeOK {
				readChannel = _readChannel
			} else {
				log.Println("standUpMobile WWriteMReadIndex in standUpPacket not type of chan StandUpPacket!")
			}
		} else {
			log.Println("standUpMobile WWriteMReadIndex not existed!")
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	done := make(chan string)
	deleteReadChannelFlag := false
	go func() {
		defer close(done)

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("read from teacher mobile client:", err)
				break
			}

			var standUpPacket StandUpPacket
			err = json.Unmarshal(data, &standUpPacket)
			if err != nil {
				log.Println(err)
				continue
			}

			select {
			case *writeChannel <- standUpPacket:
			default:
			}
		}
	}()

	startPacket := StandUpPacket{ RequestStartPacket: true }
	*writeChannel <- startPacket

	for {
		select {
		case <-done:
			close(*readChannel)
			if deleteReadChannelFlag {
				standUpChannels.Delete(readChannelIndex)
			}
			return
		case standUpPacket, ok := <-*readChannel:
			if ok {
				if standUpPacket.SayGoodbye {
					deleteReadChannelFlag = true
				}

				var data []byte
				data, err = json.Marshal(standUpPacket)
				if err != nil {
					continue
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("write to teacher mobile client", err)
					return
				}
			}
		}
	}
}

func getStudentsStatus(camStreamPath, devicePath string, faceRectNos []FaceRectToken) (studentsStatus []StudentStatus, err error) {
	body, err := sendLongTimeoutPostForm(url.Values{
		"cam_stream_path": {camStreamPath},
	}, devicePath)
	if err != nil {
		return
	}

	var faceDetectResults FaceDetectResults
	err = json.Unmarshal(body, &faceDetectResults)
	if err != nil {
		return
	}

	var studentNoRecords []string
	for _, v := range faceDetectResults.Faces {
		FaceRectNosLabel: for _, v2 := range faceRectNos {
			for _, v3 := range studentNoRecords {
				if v3 == v2.FaceToken {
					continue FaceRectNosLabel
				}
			}

			dis := math.Sqrt(math.Pow(math.Abs(float64(v.FaceRectangle.Left - v2.FaceRectangle.Left)), 2) +
				math.Pow(math.Abs(float64(v.FaceRectangle.Top - v2.FaceRectangle.Top)), 2))
			standard := math.Sqrt(float64(v2.FaceRectangle.Height * v2.FaceRectangle.Height + v2.FaceRectangle.Width * v2.FaceRectangle.Width))

			if dis < standard {
				studentsStatus = append(studentsStatus, StudentStatus{
					StudentNo: v2.FaceToken,
					Attributes: v.Attributes,
					UpdateTime: time.Now().Unix(),
				})
			}
		}
	}

	notGoodStatus := Attributes{
		Emotion: Emotion{Neutral: 100},
		HeadPose: HeadPose{
			YawAngle: 0,
			PitchAngle: 50,
			RollAngle: 0,
		},
		EyesStatus: EyesStatus{
			LeftEyeStatus: EyeStatus{
				NoGlassEyeClose: 99,
				NormalGlassEyeClose: 99,
			},
			RightEyeStatus: EyeStatus{
				NoGlassEyeClose: 99,
				NormalGlassEyeClose: 99,
			},
		},
	}
	var notDetectedStudents []StudentStatus
	for _, v := range faceRectNos {
		detected := false
		for _, v2 := range studentsStatus {
			if v.FaceToken == v2.StudentNo {
				detected = true
			}
		}

		if !detected {
			notDetectedStudents = append(notDetectedStudents, StudentStatus{
				StudentNo: v.FaceToken,
				Attributes: notGoodStatus,
				UpdateTime: time.Now().Unix(),
			})
		}
	}

	studentsStatus = append(studentsStatus, notDetectedStudents...)
	return
}

func currentStandUp(c *gin.Context) (err error) {
	teacherNo := c.Query("teacher_no")

	standUpStatus, err := getStandUpStatusRecord(-1, teacherNo)
	if err != nil {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no class"})
		return
	}

	c.JSON(http.StatusOK, standUpStatus)
	return
}

// TODO("Separate Response")
func sendStudentStatusRecords(c *gin.Context) (err error) {
	teacherNo, ByTeacher := c.GetQuery("teacher_no")
	classID, ByClass := c.GetQuery("class_id")

	var studentStatus []StudentStatusTable
	if ByTeacher {
		studentStatus, err = getStudentStatusRecordByTeacher(teacherNo)
		if err != nil {
			return
		}
	} else if ByClass {
		var id int
		id, err = strconv.Atoi(classID)
		if err != nil {
			return
		}

		studentStatus, err = getStudentStatusRecordByClass(id)
		if err != nil {
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, JsonMessage{Message: "no param provide"})
		return
	}

	studentStatusResponse := make([]StudentStatusResponse, len(studentStatus))
	for k, v := range studentStatus {
		studentStatusResponse[k].TeacherNo = v.TeacherNo
		studentStatusResponse[k].ClassID = v.ClassID
		studentStatusResponse[k].FaceCountRecordID = v.FaceCountRecordID
		studentStatusResponse[k].PDF = v.PDF
		studentStatusResponse[k].PDFPageCount = v.PDFPageCount
		studentStatusResponse[k].UpdateTime = v.UpdatedAt.Unix()
		studentStatusResponse[k].ClassName = v.ClassName

		var _studentStatusWithPage []StudentStatusWithPage
		err = json.Unmarshal([]byte(v.StudentStatus), &_studentStatusWithPage)
		if err != nil {
			return
		}
		studentStatusResponse[k].StudentStatus = _studentStatusWithPage

		var _studentWarningRecordList []StudentWarningRecord
		err = json.Unmarshal([]byte(v.StudentWarningRecordList), &_studentWarningRecordList)
		if err != nil {
			return
		}
		studentStatusResponse[k].StudentWarningRecordList = _studentWarningRecordList
	}

	c.JSON(http.StatusOK, StudentStatusListResponse{
		StudentStatus: studentStatusResponse,
		Total: len(studentStatusResponse),
	})
	return
}

func userLogin(c *gin.Context) (err error) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == AdminPermission {
		if password == AdminPermission {
			c.SetCookie("Permission", AdminPermission, 36000, "/", Domain, false, true)
			return
		}
	} else {
		var teacher *Teacher
		teacher, err = getTeacher(username)
		if err != nil {
			return
		}

		if fmt.Sprintf("%x", sha256.Sum256([]byte(password))) == teacher.TeacherPassword {
			c.SetCookie("Permission", TeacherPermission, 36000, "/", Domain, false, true)
			c.SetCookie("Username", username, 36000, "/", Domain, false, true)
			return
		}
	}

	c.JSON(http.StatusBadRequest, JsonMessage{Message: "Auth failed"})
	return
}

func mobileUserLogin(c *gin.Context) (err error) {
	userType := c.PostForm("user_type")
	username := c.PostForm("username")
	password := c.PostForm("password")

	switch userType {
	case TeacherPermission:
		var teacher *Teacher
		teacher, err = getTeacher(username)
		if err != nil {
			return
		}

		if fmt.Sprintf("%x", sha256.Sum256([]byte(password))) == teacher.TeacherPassword {
			c.JSON(http.StatusOK, JsonMessage{Message: "auth successful"})
			return
		}
	case NormalUserPermission:
		var student *Student
		student, err = getStudent(username)
		if err != nil {
			return
		}

		if fmt.Sprintf("%x", sha256.Sum256([]byte(password))) == student.StudentPassword {
			c.JSON(http.StatusOK, JsonMessage{Message: "auth successful"})
			return
		}
	}

	c.JSON(http.StatusBadRequest, JsonMessage{Message: "auth failed"})
	return
}

func sendUserInfo(c *gin.Context) (err error) {
	permission, err := c.Cookie("Permission")
	if err != nil {
		c.JSON(http.StatusOK, JsonMessage{Message: "Not Login"})
		return nil
	}

	if permission == AdminPermission {
		c.JSON(http.StatusOK, UserInfoResp{
			User: UserInfo{Username: AdminPermission, Permissions: AdminPermission},
		})
	} else {
		username := ""
		username, err = c.Cookie("Username")
		if err != nil {
			c.JSON(http.StatusOK, JsonMessage{Message: "Not Login"})
			return nil
		}

		c.JSON(http.StatusOK, UserInfoResp{
			User: UserInfo{Username: username, Permissions: TeacherPermission},
		})
	}

	return
}

func userLogout(c *gin.Context) (err error) {
	c.SetCookie("Permission", "", -1, "/", Domain, false, true)
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

func savePDF(c *gin.Context) (err error) {
	file, _ := c.FormFile("file")

	fileName := fmt.Sprintf("%v%v", time.Now().UnixNano(), path.Ext(file.Filename))

	err = c.SaveUploadedFile(file, PDFFileDir + fileName)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, JsonMessage{Message: fileName})
	return
}

func sendPDF(c *gin.Context) (err error) {
	pdfName := c.Param("name")

	c.File(PDFFileDir + pdfName)
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

func fileUploadRequest(url string, params map[string]string, fileParamName string, fileContent []byte, fileName string) (respBody []byte, err error) {
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
	if err != nil {
		return
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf(response.Status)
		return
	}

	respBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}

func checkIfInStringSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func updateFaceSetByStudentNos(oldStudents []Student, newStudentNos []string, faceSetToken string) (newStudents []*Student, err error) {
	var studentsAdd []string
	studentsDelete := ""
	oldStudentNos := make([]string, len(oldStudents))

	for k, v := range oldStudents {
		oldStudentNos[k] = *v.StudentNo
	}

	for _, v := range oldStudents {
		if !checkIfInStringSlice(newStudentNos, *v.StudentNo) {
			studentsDelete += fmt.Sprintf("%v,", v.FaceToken)
		}
	}

	for _, v := range newStudentNos {
		var student *Student
		student, err = getStudent(v)

		if student.ID == 0 {
			continue
		}
		newStudents = append(newStudents, student)
		if !checkIfInStringSlice(oldStudentNos, v) {
			studentsAdd = append(studentsAdd, student.FaceToken)
		}
	}

	studentsDelete = strings.TrimSuffix(studentsDelete, ",")

	removeNum := 0
	removed := 0
	added := 0

	if len(studentsDelete) > 0 {
		removeNum = len(strings.Split(studentsDelete, ","))
	}

	var body []byte
	var resp UpdateFaceResp
	if removeNum > 0 {
		body, err = sendPostForm(url.Values{
			"api_key": {config.ApiKey},
			"api_secret": {config.ApiSecret},
			"faceset_token": {faceSetToken},
			"face_tokens": {studentsDelete},
		}, config.DeleteFaceUrl)
		if err != nil {
			return
		}

		err = json.Unmarshal(body, &resp)
		if err != nil {
			return
		}

		removed = resp.FaceRemoved
	}

	if len(studentsAdd) > 0 {
		for i := 0; i < int(math.Ceil(float64(len(studentsAdd)) / 5)); i++ {
			addString := ""
			var addSlice []string
			if (len(studentsAdd) - i * 5) >= 5 {
				addSlice = studentsAdd[i*5:i*5+4]
			} else {
				addSlice = studentsAdd[i*5:]
			}

			for _, v := range addSlice {
				addString += fmt.Sprintf("%v,", v)
			}
			addString = strings.TrimSuffix(addString, ",")

			body, err = sendPostForm(url.Values{
				"api_key": {config.ApiKey},
				"api_secret": {config.ApiSecret},
				"faceset_token": {faceSetToken},
				"face_tokens": {addString},
			}, config.AddFaceUrl)
			if err != nil {
				return
			}

			err = json.Unmarshal(body, &resp)
			if err != nil {
				return
			}
			added += resp.FaceAdded
		}
	}

	if removed != removeNum || added != len(studentsAdd) {
		log.Println(studentsAdd)
		log.Println(studentsDelete)
		log.Printf("students not completely removed or added %v  %v  %v  %v\n", removed, removeNum, added, len(studentsAdd))
		return
	}
	return
}

func sendPostForm(params url.Values, url string) (body []byte, err error) {
	response, err := client.PostForm(url, params)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf(response.Status)
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func sendLongTimeoutPostForm(params url.Values, url string) (body []byte, err error) {
	response, err := longTimeoutClient.PostForm(url, params)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf(response.Status)
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func studentWarningToIntList(record []StudentWarningRecord) (list []int) {
	list = make([]int, len(record))

	for k, v := range record {
		list[k] = v.Warning
	}

	return
}

//func clearDBError(c *gin.Context) (err error) {
//	err = clearStandUpStatusTable()
//	if err != nil {
//		return
//	}
//
//	return
//}

func getStandUpStatusRecord(classID int, teacherNo string) (record *StandUpStatusRecord, err error) {
	if classID > 0 {
		standUpChannels.Range(func(key, value interface{}) bool {
			if _r, ok := value.(*StandUpStatusRecord); ok {
				if _r.ClassID == classID {
					record = _r
					return false
				}
			}
			return true
		})
		if record == nil {
			err = fmt.Errorf("no standUp status record classID %v or teacherNo %v", classID, teacherNo)
		}
		return
	} else if len(teacherNo) > 0 {
		standUpChannels.Range(func(key, value interface{}) bool {
			if _r, ok := value.(*StandUpStatusRecord); ok {
				if _r.TeacherNo == teacherNo {
					record = _r
					return false
				}
			}
			return true
		})
		if record == nil {
			err = fmt.Errorf("no standUp status record type mismatch classID %v or teacherNo %v", classID, teacherNo)
		}
		return
	}

	err = fmt.Errorf("no standUp status record for classID %v or teacherNo %v", classID, teacherNo)
	return
}