package main

import "time"

type DetectedData struct {
	X0 int `json:"x0"`
	Y0 int `json:"y0"`
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
}

type PersonData struct {
	DetectedData DetectedData `json:"detected_data"`
	Token string `json:"face_token"`
}

type DetectedImage struct {
	DetectedData DetectedData
	Data []byte
}

type SearchFaceResult struct {
	Confidence float32 `json:"confidence"`
	FaceToken string `json:"face_token"`
}

type SearchFaceResults struct {
	Results []SearchFaceResult `json:"results"`
}

type JsonError struct {
	Error string `json:"error"`
}

type PersonDataFaces struct {
	Faces []PersonData `json:"faces"`
}

type Config struct {
	LocalPort string `json:"local_port"`
	ServerAddr string `json:"server_addr"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	SearchFaceUrl string `json:"search_face_url"`
	DetectSkeletonUrl string `json:"detect_skeleton_url"`
	DetectInterval int `json:"detect_interval"`
	Classrooms []Classroom `json:"classrooms"`
}

type Classroom struct {
	ClassroomID uint `json:"classroom_id"`
	ClassroomName string `json:"classroom_name"`
	CamPath string `json:"cam_path"`
}

type ClassroomStat struct {
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
	Stats SystemStats `json:"system_stats"`
	Classrooms []ClassroomStat `json:"classrooms"`
}