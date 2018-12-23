package main

import (
	"github.com/jinzhu/gorm"
	"time"
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

type JsonError struct {
	Error string `json:"error"`
}

type Config struct {
	LocalPort string `json:"local_port"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	DetectFaceUrl string `json:"search_face_url"`
}

type ClassResponse struct {
	ClassID uint `json:"class_id"`
	ClassName string `json:"class_name"`
	FaceCount int `json:"face_count"`
	FaceSetToken string `json:"faceset_token"`
}
type ClassesResponse struct {
	Classes []ClassResponse `json:"classes"`
}

type StudentResponse struct {
	StudentNo string `json:"student_no"`
	FaceToken string `json:"face_token"`
	ClassIDs []uint `json:"class_ids"`
}
type StudentsResponse struct {
	Students []StudentResponse `json:"students"`
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
	ClassName *string `gorm:"unique;not null" json:"class_name"`
	Students []*Student `gorm:"many2many:student_class;"`
}

type Student struct {
	gorm.Model
	StudentNo *string `gorm:"unique;not null"`
	FaceToken string
	Classes []*Class `gorm:"many2many:student_class;"`
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
	CpuUsed string `json:"cpu_used"`
	MemUsed string `json:"mem_used"`
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
	CpuUsed string
	MemUsed string
	DeviceID uint
}

type ClassroomStatsTable struct {
	gorm.Model
	UpdateTime time.Time
	PersonCount int
	Persons []DetectedData
	ClassroomID uint
}

type SingleClassroomStats struct {
	UpdateTime time.Time `json:"update_time"`
	ClassroomStats ClassroomStats `json:"classroom_stats"`
}