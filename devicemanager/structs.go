package main

import "github.com/jinzhu/gorm"

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

type FaceImageTokens struct {
	Faces []FaceImageToken `json:"faces"`
}

type FaceImageToken struct {
	FaceImage string `json:"face_image"`
	FaceToken string `json:"face_token"`
	StudentNo string `json:"student_no"`
}

type JsonError struct {
	Error string `json:"error"`
}

type FaceCountToken struct {
	FaceCount int `json:"face_count"`
	FaceSetToken string `json:"faceset_token"`
}

type Config struct {
	LocalPort string `json:"local_port"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	DetectFaceUrl string `json:"search_face_url"`
}

type Class struct {
	gorm.Model
	FaceSetToken string
	ClassName string
	Students []*Student `gorm:"many2many:student_class;"`
}

type Student struct {
	gorm.Model
	StudentNo string
	FaceToken string
	Classes []*Class `gorm:"many2many:student_class;"`
}

type Device struct {
	gorm.Model
	DevicePath string
	CamPath string
	Classroom string
}