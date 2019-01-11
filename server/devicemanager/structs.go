package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	Admin = "admin"
	Developer = "developer"
)

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

type FaceNoToken struct {
	FaceToken string `json:"face_token"`
	StudentNo string `json:"student_no"`
}

type JsonMessage struct {
	Message string `json:"message"`
}

type Config struct {
	LocalPort string `json:"local_port"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	DetectFaceUrl string `json:"detect_face_url"`
	CreateFaceSetUrl string `json:"create_face_set_url"`
	DeleteFaceSetUrl string `json:"delete_face_set_url"`
}

type DeviceResponse struct {
	DeviceID uint `json:"device_id"`
	DevicePath string `json:"device_path"`
	DevicePort string `json:"device_port"`
}
type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

type ClassResponse struct {
	ClassID uint `json:"class_id"`
	ClassName string `json:"class_name"`
	ClassImage string `json:"class_image"`
	FaceCount int `json:"face_count"`
	FaceSetToken string `json:"faceset_token"`
	StudentNos []string `json:"student_nos"`
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
	CamPath string `json:"cam_path"`
	DeviceID uint `json:"device_id"`
}
type CamerasResponse struct {
	Cameras []CameraResponse `json:"cameras"`
}

type ClassroomResponse struct {
	ClassroomID uint `json:"classroom_id"`
	ClassroomName string `json:"classroom_name"`
	CameraID uint `json:"camera_id"`
}
type ClassroomsResponse struct {
	Classrooms []ClassroomResponse `json:"classrooms"`
}

type Class struct {
	gorm.Model
	FaceSetToken string
	ClassName string
	ClassImage string
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
	StudentNo *string `gorm:"unique;not null"`
	TeacherName string `json:"teacher_name"`
	TeacherImage string `json:"teacher_image"`
	TeacherPassword string
	Classes []*Class `gorm:"many2many:teacher_class;"`
}

// DevicePort: eg.(":8000")
type Device struct {
	gorm.Model
	DevicePath string
	DevicePort string
}

type Camera struct {
	gorm.Model
	CamPath string
	DeviceID uint
}

type Classroom struct {
	gorm.Model
	Name string
	CameraID uint
}

type ClassroomStats struct {
	ClassroomID uint `json:"classroom_id"`
	ClassroomName string `json:"classroom_name"`
	PersonCount int `json:"person_count"`
	Persons []DetectedData `json:"persons"`
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

type DetectedData struct {
	X0 int `json:"x0"`
	Y0 int `json:"y0"`
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
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
	Persons []DetectedData
	ClassroomID uint
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

func newClassesResponse(classes []Class) (classesResp *ClassesResponse) {
	classesResp = &ClassesResponse{}

	_classesResp := make([]ClassResponse, len(classes))
	for i := 0; i < len(classes); i++ {
		_classesResp[i].ClassID = classes[i].ID
		_classesResp[i].ClassName = classes[i].ClassName
		_classesResp[i].FaceCount = len(classes[i].Students)
		_classesResp[i].FaceSetToken = classes[i].FaceSetToken
		_classesResp[i].ClassImage = classes[i].ClassImage
	}

	classesResp.Classes = _classesResp
	classesResp.Total = len(_classesResp)
	return
}

func newStudentsResponse(students []Student) (studentsResp *StudentsResponse) {
	studentsResp = &StudentsResponse{}

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

	studentsResp.Students = studentsResponse
	studentsResp.Total = len(studentsResponse)
	return
}