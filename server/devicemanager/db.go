package main

import (
	"github.com/jinzhu/gorm"
	"log"
)

func init() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Println("failed to connect database")
		return
	}
	defer db.Close()

	db.AutoMigrate(&Class{})
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Device{})
	db.AutoMigrate(&Camera{})
	db.AutoMigrate(&Classroom{})
	db.AutoMigrate(&DeviceStatsTable{})
	db.AutoMigrate(&ClassroomStatsTable{})
}

func getAllClasses() ([]Class, error) {
	var err error

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var classes []Class
	db.Find(&classes)

	return classes, nil
}

//func getLastClass() (class *Class, err error) {
//	class = &Class{}
//
//	db, err := gorm.Open("sqlite3", "test.db")
//	if err != nil {
//		return
//	}
//	defer db.Close()
//
//	db.Last(class)
//	return
//}

func getClass(id int) (class *Class, err error) {
	class = &Class{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(class, "id = ?", id)
	return
}

func getClasses(ids []int) (classes []Class, err error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&classes, "id in (?)", ids)
	return
}

func getAllDevices() (devices []Device, err error) {
	devices = []Device{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.Find(&devices)
	return
}


func getDevice(id int) (device *Device, err error) {
	device = &Device{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(device, "id = ?", id)
	return
}

func getDeviceByPath(devicePath string) (device *Device, err error) {
	device = &Device{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(device, "device_path = ?", devicePath)
	return
}

func getClassroom(id int) (classroom *Classroom, err error) {
	classroom = &Classroom{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(classroom, "id = ?", id)
	return
}

func getStudent(studentNo string) (student *Student, err error) {
	student = &Student{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(student, "student_no = ?", studentNo)
	return
}

func getStudentsByClass(classID int) (students []Student, err error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	var class Class
	db.First(&class, "id = ?", classID)

	db.Model(&class).Related(&students, "Students")
	return
}

func getClassroomStatsItem(classroomID int) (stats *ClassroomStatsTable, err error) {
	stats = &ClassroomStatsTable{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.Last(stats, "classroom_id = ?", classroomID)
	return
}

func getCamera(cameraID int) (camera *Camera, err error) {
	camera = &Camera{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(camera, "id = ?", cameraID)
	return
}

func getCameras() (cameras []Camera, err error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&cameras)
	return
}

func getClassrooms() (classrooms []Classroom, err error) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&classrooms)
	return
}

func createTableItem(v interface{}) error {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}
	defer db.Close()

	db.Create(v)

	return nil
}