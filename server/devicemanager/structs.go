package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

const (
	Admin = "admin"
	Developer = "developer"
)

type FaceCountRecord struct {
	gorm.Model
	ClassID int
	FaceRectTokens string
}

type FaceRectTokens struct {
	Faces []FaceRectToken `json:"faces"`
}

type FaceRectToken struct {
	FaceRectangle FaceRectangle `json:"face_rectangle"`
	FaceToken string `json:"face_token"`
}

type FaceRectangle struct {
	Width int `json:"width"`
	Top int `json:"top"`
	Left int `json:"left"`
	Height int `json:"height"`
}

type PersonData struct {
	Face FaceAnalyzeResult `json:"face"`
	GlobalWidth float64 `json:"global_width"`
	GlobalHeight float64 `json:"global_height"`
	PersonCount int `json:"person_count"`
	ImageUrl string `json:"image_url"`
}

type FaceAnalyzeResult struct {
	Attributes Attributes `json:"attributes"`
	FaceRectangle FaceRectangle `json:"face_rectangle"`
	FaceToken string `json:"face_token"`
}

type Attributes struct {
	Emotion Emotion `json:"emotion"`
	EyesStatus EyesStatus `json:"eyes_status"`
	HeadPose HeadPose `json:"head_pose"`
}

type Emotion struct {
	Sadness float32 `json:"sadness"`
	Neutral float32 `json:"neutral"`
	Disgust float32 `json:"disgust"`
	Anger float32 `json:"anger"`
	Surprise float32 `json:"surprise"`
	Fear float32 `json:"fear"`
	Happiness float32 `json:"happiness"`
}

type EyeStatus struct {
	NoGlassEyeClose float32 `json:"no_glass_eye_close"`
	NormalGlassEyeClose float32 `json:"normal_glass_eye_close"`
}

type EyesStatus struct {
	LeftEyeStatus EyeStatus `json:"left_eye_status"`
	RightEyeStatus EyeStatus `json:"right_eye_status"`
}

type HeadPose struct {
	YawAngle float32 `json:"yaw_angle"`
	PitchAngle float32 `json:"pitch_angle"`
	RollAngle float32 `json:"roll_angle"`
}

type JsonMessage struct {
	Message string `json:"message"`
}

type StandUpPacket struct {

}

type StandUpStatusTable struct {
	gorm.Model
	ClassID int
	WReadMWriteIndex int
	WWriteMReadIndex int
}

type Config struct {
	LocalPort string `json:"local_port"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	DetectFaceUrl string `json:"detect_face_url"`
	CreateFaceSetUrl string `json:"create_face_set_url"`
	DeleteFaceSetUrl string `json:"delete_face_set_url"`
	AddFaceUrl string `json:"add_face_url"`
	DeleteFaceUrl string `json:"delete_face_url"`
}

type DeviceResponse struct {
	DeviceID uint `json:"device_id"`
	DevicePath string `json:"device_path"`
	DevicePort string `json:"device_port"`
	CameraIDs []uint `json:"camera_ids"`
}
type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Total int `json:"total"`
}

type ClassResponse struct {
	ClassID uint `json:"class_id"`
	ClassName string `json:"class_name"`
	ClassImage string `json:"class_image"`
	FaceCount int `json:"face_count"`
	FaceSetToken string `json:"faceset_token"`
	ClassroomNo string `json:"classroom_no"`
	StudentNos []string `json:"student_nos"`
	TeacherNos []string `json:"teacher_nos"`
}
type ClassesResponse struct {
	Classes []ClassResponse `json:"classes"`
	Total int `json:"total"`
}

type StudentResponse struct {
	StudentNo string `json:"student_no"`
	StudentName string `json:"student_name"`
	FaceToken string `json:"face_token"`
	StudentImage string `json:"student_image"`
	StudentPassword string `json:"student_password"`
	ClassIDs []uint `json:"class_ids"`
}
type StudentsResponse struct {
	Students []StudentResponse `json:"students"`
	Total int `json:"total"`
}

type CameraResponse struct {
	CameraID uint `json:"camera_id"`
	CamStreamPath string `json:"cam_stream_path"`
	CamONVIFPath string `json:"cam_onvif_path"`
	CamAuthName string `json:"cam_auth_name"`
	CamAuthPassword string `json:"cam_auth_password"`
	ClassroomNo string `json:"classroom_no"`
	DeviceID string `json:"device_id"`
}
type CamerasResponse struct {
	Cameras []CameraResponse `json:"cameras"`
	Total int `json:"total"`
}

type ClassroomResponse struct {
	ClassroomNo string `json:"classroom_no"`
	CameraID string `json:"camera_id"`
}
type ClassroomsResponse struct {
	Classrooms []ClassroomResponse `json:"classrooms"`
	Total int `json:"total"`
}

type TeacherResponse struct {
	TeacherNo string `json:"teacher_no"`
	TeacherName string `json:"teacher_name"`
	TeacherImage string `json:"teacher_image"`
	TeacherPassword string `json:"teacher_password"`
	ClassIDs []uint `json:"class_ids"`
}
type TeachersResponse struct {
	Teachers []TeacherResponse `json:"teachers"`
	Total int `json:"total"`
}

type Class struct {
	gorm.Model
	FaceSetToken string
	ClassName string
	ClassImage string
	ClassroomNo string
	Students []*Student `gorm:"many2many:student_class;"`
	Teachers []*Teacher `gorm:"many2many:teacher_class;"`
}

type Student struct {
	gorm.Model
	StudentNo *string `gorm:"unique;not null"`
	StudentName string
	StudentImage string
	StudentPassword string
	FaceToken string
	Classes []*Class `gorm:"many2many:student_class;"`
}

type Teacher struct {
	gorm.Model
	TeacherNo *string `gorm:"unique;not null"`
	TeacherName string
	TeacherImage string
	TeacherPassword string
	Classes []*Class `gorm:"many2many:teacher_class;"`
}

type Device struct {
	gorm.Model
	DevicePath string
	DevicePort string
	Cameras []*Camera `gorm:"many2many:device_camera;"`
}

type Camera struct {
	gorm.Model
	CamStreamPath string
	CamONVIFPath string
	CamAuthName string
	CamAuthPassword string
	Devices []*Device `gorm:"many2many:device_camera;"`
	Classrooms []*Classroom `gorm:"many2many:classroom_camera;"`
}

type Classroom struct {
	gorm.Model
	ClassroomNo *string `gorm:"unique;not null"`
	Cameras []*Camera `gorm:"many2many:classroom_camera;"`
}

type ClassroomStats struct {
	ClassroomNo string `json:"classroom_no"`
	PersonCount int `json:"person_count"`
	Persons []FaceRectangle `json:"persons"`
}

type SystemStats struct {
	CpuUsed float64 `json:"cpu_used"`
	MemUsed float64 `json:"mem_used"`
}

type Stats struct {
	UpdateTime time.Time `json:"update_time"`
	SystemStats SystemStats `json:"system_stats"`
	Classrooms []ClassroomStats `json:"classrooms"`
}

type DeviceStatsTable struct {
	gorm.Model
	UpdateTime time.Time
	CpuUsed float64
	MemUsed float64
	DeviceID uint
}

type ClassroomStatsTable struct {
	gorm.Model
	UpdateTime time.Time
	PersonCount int
	Persons string
	ClassroomNo string
}

type DeviceManagerSystemStats struct {
	gorm.Model
	CpuUsed float64 `json:"cpu_used"`
	MemUsed float64 `json:"mem_used"`
}

type SingleClassroomStats struct {
	UpdateTime time.Time `json:"update_time"`
	ClassroomStats ClassroomStats `json:"classroom_stats"`
}

type UserInfo struct {
	Username string `json:"username"`
	Permissions string `json:"permissions"`
}

type UserInfoResp struct {
	User UserInfo `json:"user"`
}

type NumberCard struct {
	Devices int `json:"devices"`
	Cameras int `json:"cameras"`
}

type DashBoardResp struct {
	SystemStats []DeviceManagerSystemStats `json:"system_stats"`
	NumberCard NumberCard `json:"number_card"`
}

type UpdateFaceResp struct {
	FaceAdded int `json:"face_added"`
	FaceRemoved int `json:"face_removed"`
}

func newClassesResponse(classes []Class, page, pageSize string) (classesResp *ClassesResponse, err error) {
	classesResp = &ClassesResponse{}

	classesResp.Total = len(classes)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(classes) - start * size) > size {
			classes = classes[(start - 1) * size: start * size - 1]
		} else {
			classes = classes[(start - 1) * size:]
		}
	}

	_classesResp := make([]ClassResponse, len(classes))
	for i := 0; i < len(classes); i++ {
		_classesResp[i].ClassID = classes[i].ID
		_classesResp[i].ClassName = classes[i].ClassName
		_classesResp[i].FaceSetToken = classes[i].FaceSetToken
		_classesResp[i].ClassImage = classes[i].ClassImage
		_classesResp[i].ClassroomNo = classes[i].ClassroomNo

		var students []Student
		students, err = getStudentsByClass(int(classes[i].ID))
		if err != nil {
			return
		}
		studentNos := make([]string, len(students))
		for k, v := range students {
			if v.ID == 0 {
				continue
			}
			studentNos[k] = *v.StudentNo
		}

		var teachers []Teacher
		teachers, err = getTeachersByClass(int(classes[i].ID))
		if err != nil {
			return
		}
		teacherNos := make([]string, len(teachers))
		for k, v := range teachers {
			if v.ID == 0 {
				continue
			}
			teacherNos[k] = *v.TeacherNo
		}

		_classesResp[i].StudentNos = studentNos
		_classesResp[i].TeacherNos = teacherNos
		_classesResp[i].FaceCount = len(students)
	}

	classesResp.Classes = _classesResp
	return
}

func newStudentsResponse(students []Student, page, pageSize string) (studentsResp *StudentsResponse, err error) {
	studentsResp = &StudentsResponse{}

	studentsResp.Total = len(students)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(students) - start * size) > size {
			students = students[(start - 1) * size: start * size - 1]
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

		var classes []Class
		classes, err = getClassesByStudentNo(*v.StudentNo)
		if err != nil {
			return
		}
		classUintIDs := make([]uint, len(classes))
		for k, v := range classes {
			if v.ID == 0 {
				continue
			}
			classUintIDs[k] = v.ID
		}
		studentsResponse[k].ClassIDs = classUintIDs
	}

	studentsResp.Students = studentsResponse
	return
}

func newDevicesResponse(devices []Device, page, pageSize string) (deviceResp *DevicesResponse, err error) {
	deviceResp = &DevicesResponse{}
	deviceResp.Total = len(devices)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(devices) - start * size) > size {
			devices = devices[(start - 1) * size: start * size - 1]
		} else {
			devices = devices[(start - 1) * size:]
		}
	}

	devicesResponse := make([]DeviceResponse, len(devices))
	for i := 0; i < len(devicesResponse); i++ {
		devicesResponse[i].DeviceID = devices[i].ID
		devicesResponse[i].DevicePath = devices[i].DevicePath
		devicesResponse[i].DevicePort = devices[i].DevicePort

		var cameras []Camera
		cameras, err = getCamerasByDevice(int(devices[i].ID))
		if err != nil {
			return
		}

		cameraIDs := make([]uint, len(cameras))
		for k, v := range cameras {
			if v.ID == 0 {
				continue
			}
			cameraIDs[k] = v.ID
		}

		devicesResponse[i].CameraIDs = cameraIDs
	}

	deviceResp.Devices = devicesResponse
	return
}

func newCamerasResponse(cameras []Camera, page, pageSize string) (camerasResp *CamerasResponse, err error) {
	camerasResp = &CamerasResponse{}
	camerasResp.Total = len(cameras)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(cameras) - start * size) > size {
			cameras = cameras[(start - 1) * size: start * size - 1]
		} else {
			cameras = cameras[(start - 1) * size:]
		}
	}

	camerasResponse := make([]CameraResponse, len(cameras))
	for k, v := range cameras {
		camerasResponse[k].CamStreamPath = v.CamStreamPath
		camerasResponse[k].CameraID = v.ID
		camerasResponse[k].CamAuthPassword = v.CamAuthPassword
		camerasResponse[k].CamAuthName = v.CamAuthName
		camerasResponse[k].CamONVIFPath = v.CamONVIFPath

		var devices []Device
		devices, err = getDevicesByCamera(int(v.ID))
		if err != nil {
			return
		}
		if len(devices) > 0 && devices[0].ID != 0 {
			camerasResponse[k].DeviceID = fmt.Sprint(devices[0].ID)
		}

		var classrooms []Classroom
		classrooms, err = getClassroomsByCamera(int(v.ID))
		if err != nil {
			return
		}
		if len(classrooms) > 0 && classrooms[0].ID != 0 {
			camerasResponse[k].ClassroomNo = *classrooms[0].ClassroomNo
		}
	}

	camerasResp.Cameras = camerasResponse
	return
}

func newClassroomsResponse(classrooms []Classroom, page, pageSize string) (classroomsResp *ClassroomsResponse, err error) {
	classroomsResp = &ClassroomsResponse{}
	classroomsResp.Total = len(classrooms)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(classrooms) - start * size) > size {
			classrooms = classrooms[(start - 1) * size: start * size - 1]
		} else {
			classrooms = classrooms[(start - 1) * size:]
		}
	}

	classroomsResponse := make([]ClassroomResponse, len(classrooms))
	for k, v := range classrooms {
		classroomsResponse[k].ClassroomNo = *v.ClassroomNo

		var cameras []Camera
		cameras, err = getCamerasByClassroom(*v.ClassroomNo)
		if err != nil {
			return
		}
		if len(cameras) > 0 && cameras[0].ID != 0 {
			classroomsResponse[k].CameraID = fmt.Sprint(cameras[0].ID)
		}
	}

	classroomsResp.Classrooms = classroomsResponse
	return
}

func newTeacherResponse(teachers []Teacher, page, pageSize string) (teachersResp *TeachersResponse, err error) {
	teachersResp = &TeachersResponse{}
	teachersResp.Total = len(teachers)

	if page != "" && pageSize != "" {
		var start, size int
		start, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		size, err = strconv.Atoi(pageSize)
		if err != nil {
			return
		}

		if (len(teachers) - start * size) > size {
			teachers = teachers[(start - 1) * size: start * size - 1]
		} else {
			teachers = teachers[(start - 1) * size:]
		}
	}

	_teacherResp := make([]TeacherResponse, len(teachers))
	for k, v := range teachers {
		_teacherResp[k].TeacherPassword = v.TeacherPassword
		_teacherResp[k].TeacherName = v.TeacherName
		_teacherResp[k].TeacherImage = v.TeacherImage
		_teacherResp[k].TeacherNo = *v.TeacherNo

		var classes []Class
		classes, err = getClassesByTeacherNo(*v.TeacherNo)
		if err != nil {
			return
		}
		classUintIDs := make([]uint, len(classes))
		for k, v := range classes {
			if v.ID == 0 {
				continue
			}
			classUintIDs[k] = v.ID
		}
		_teacherResp[k].ClassIDs = classUintIDs
	}

	teachersResp.Teachers = _teacherResp
	return
}