package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
	router.DELETE("/classes", func(c *gin.Context) {
		if err := deleteClass(c); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, JsonMessage{Message: "delete class failed"})
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
	router.PATCH("/teachers/:id", func(c *gin.Context) {
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
