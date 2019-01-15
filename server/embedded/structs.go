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

type JsonMessage struct {
	Message string `json:"message"`
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
	FaceDetectParam string `json:"face_detect_param"`
	FaceDetectBin string `json:"face_detect_bin"`
	BodyDetectParam string `json:"body_detect_param"`
	BodyDetectBin string `json:"body_detect_bin"`
	Classrooms []Classroom `json:"classrooms"`
}

type Classroom struct {
	ClassroomNo string `json:"classroom_no"`
	CamStreamPath string `json:"cam_stream_path"`
}

type ClassroomStat struct {
	ClassroomNo string `json:"classroom_no"`
	PersonCount int `json:"person_count"`
	Persons []DetectedData `json:"persons"`
}

type SystemStats struct {
	CpuUsed float64 `json:"cpu_used"`
	MemUsed float64 `json:"mem_used"`
}

type Stats struct {
	UpdateTime time.Time `json:"update_time"`
	Stats SystemStats `json:"system_stats"`
	Classrooms []ClassroomStat `json:"classrooms"`
}