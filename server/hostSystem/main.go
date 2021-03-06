package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	var err error

	gin.SetMode(gin.ReleaseMode)
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
	router.MaxMultipartMemory = 20 << 20

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
	router.DELETE("/classes", func(c *gin.Context) {
		if err := deleteClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete class failed"})
		}
	})

	// 点名
	router.GET("/face_count", func(c *gin.Context) {
		if err := faceCount(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "face count error"})
		}
	})
	router.GET("/face_count_record", func(c *gin.Context) {
		if err := sendFaceCountRecord(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send face count error"})
		}
	})

	// 上课
	router.GET("/stand_up", func(c *gin.Context) {
		if err := standUp(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "stand up start error"})
		}
	})
	router.GET("/stand_up_mobile", func(c *gin.Context) {
		if err := standUpMobile(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "stand up mobile start error"})
		}
	})
	router.GET("/current_stand_up", func(c *gin.Context) {
		if err := currentStandUp(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send current stand up error"})
		}
	})
	router.GET("/student_status_records", func(c *gin.Context) {
		if err := sendStudentStatusRecords(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send student status error"})
		}
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
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send classroom stats error"})
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
		if err := updateDevices(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update device error"})
		}
	})
	router.DELETE("/devices", func(c *gin.Context) {
		if err := deleteDevices(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete device error"})
		}
	})
	router.GET("/device_stats", func(c *gin.Context) {
		if err := sendDeviceStats(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send device stats error"})
		}
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
	router.PATCH("/students", func(c *gin.Context) {
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

	// 教师
	router.POST("/teachers", func(c *gin.Context) {
		if err := createTeacher(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "create teacher error"})
		}
	})
	router.GET("/teachers", func(c *gin.Context) {
		if err := sendTeachers(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "send teacher error"})
		}
	})
	router.PATCH("/teachers", func(c *gin.Context) {
		if err := updateTeacher(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update teacher error"})
		}
	})
	router.DELETE("/teachers", func(c *gin.Context) {
		if err := deleteTeacher(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete teacher error"})
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
		if err := updateCameras(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "update cameras error"})
		}
	})
	router.DELETE("cameras", func(c *gin.Context) {
		if err := deleteCameras(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete cameras error"})
		}
	})

	// 教室
	router.POST("/classrooms", handlerFunc(createClassroom, "create classroom error"))
	router.GET("/classrooms", handlerFunc(sendClassrooms, "get classrooms error"))
	router.PATCH("/classrooms", handlerFunc(updateClassrooms, "update classrooms error"))
	router.DELETE("/classrooms", handlerFunc(deleteClassrooms, "delete classrooms error"))

	// 设置
	//router.GET("/config", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, config)
	//})
	//router.POST("/config", func(c *gin.Context) {
	//	if err := setConfig(c); err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, JsonMessage{Message: "can not set config"})
	//	}
	//})

	// 用户
	router.POST("/user/login", handlerFunc(userLogin, "user login error"))
	router.GET("/user", handlerFunc(sendUserInfo, "get userInfo error"))
	router.GET("/user/logout", handlerFunc(userLogout, "user logout error"))
	router.POST("/user/mobile_login", handlerFunc(mobileUserLogin, "mobile user login error"))

	// 仪表盘
	router.GET("/dashboard",handlerFunc(sendDashBoard, "get dashboard error"))

	// 图片
	router.GET("/images/:name", handlerFunc(sendImage, "get image error"))
	router.POST("/images", handlerFunc(saveImage, "save image error"))

	//PDF
	router.GET("/pdf/:name", handlerFunc(sendPDF, "get pdf error"))
	router.POST("/pdf", handlerFunc(savePDF, "save pdf error"))

	//router.GET("/clear", func(c *gin.Context) {
	//	if err := clearDBError(c); err != nil {
	//		log.Println(err)
	//		c.JSON(http.StatusInternalServerError, JsonMessage{Message: "clear"})
	//	}
	//})

	return router
}

func handlerFunc(handler func(c *gin.Context) (err error), jsonMessage string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: jsonMessage})
		}
	}
}