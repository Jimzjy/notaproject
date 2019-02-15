package main

import "time"

type FaceRectangle struct {
	Top int `json:"top"`
	Left int `json:"left"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type PersonData struct {
	Face FaceAnalyzeResult `json:"face"`
	GlobalWidth float64 `json:"global_width"`
	GlobalHeight float64 `json:"global_height"`
	PersonCount int `json:"person_count"`
	ImageUrl string `json:"image_url"`
}

type FaceDetectResults struct {
	Faces []FaceAnalyzeResult `json:"faces"`
}

type FaceAnalyzeResult struct {
	Attributes Attributes `json:"attributes"`
	FaceRectangle FaceRectangle `json:"face_rectangle"`
	FaceToken string `json:"face_token"`
}

type Attributes struct {
	Emotion Emotion `json:"emotion"`
	EyesStatus EyesStatus `json:"eyestatus"`
	HeadPose HeadPose `json:"headpose"`
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

//type DetectedImage struct {
//	FaceRectangle FaceRectangle
//	Data []byte
//}

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

//type PersonDataFaces struct {
//	Faces []PersonData `json:"faces"`
//}

type Config struct {
	LocalPort string `json:"local_port"`
	ServerAddr string `json:"server_addr"`
	ApiKey string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	SearchFaceUrl string `json:"search_face_url"`
	DetectFaceUrl string `json:"detect_face_url"`
	AnalyzeFaceUrl string `json:"analyze_face_url"`
	DetectInterval int `json:"detect_interval"`
	FaceDetectParam string `json:"face_detect_param"`
	FaceDetectBin string `json:"face_detect_bin"`
	BodyDetectParam string `json:"body_detect_param"`
	BodyDetectBin string `json:"body_detect_bin"`
	Qps int `json:"qps"`
	Classrooms []Classroom `json:"classrooms"`
}

type Classroom struct {
	ClassroomNo string `json:"classroom_no"`
	CamStreamPath string `json:"cam_stream_path"`
}

type ClassroomStat struct {
	ClassroomNo string `json:"classroom_no"`
	PersonCount int `json:"person_count"`
	GlobalWidth float64 `json:"global_width"`
	GlobalHeight float64 `json:"global_height"`
	Persons []FaceRectangle `json:"persons"`
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