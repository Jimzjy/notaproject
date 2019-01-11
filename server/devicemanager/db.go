package main

import (
	"github.com/jinzhu/gorm"
	"log"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	DBName  = "deviceManager.db"
	DB = "sqlite3"
)

func init() {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		log.Println("failed to connect database")
		return
	}
	defer db.Close()

	db.AutoMigrate(&Class{})
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Teacher{})
	db.AutoMigrate(&Device{})
	db.AutoMigrate(&Camera{})
	db.AutoMigrate(&Classroom{})
	db.AutoMigrate(&DeviceStatsTable{})
	db.AutoMigrate(&ClassroomStatsTable{})
	db.AutoMigrate(&DeviceManagerSystemStats{})
}

func getAllClasses() ([]Class, error) {
	var err error

	db, err := gorm.Open(DB, DBName)
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
//	db, err := gorm.Open(DB, DBName)
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

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(class, "id = ?", id)
	return
}

func getClasses(ids []int) (classes []Class, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&classes, "id in (?)", ids)
	return
}

func getClassesByName(className string) (classes []Class, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&classes, "class_name = ?", className)
	return
}

func getClassesByStudentNo(studentNo string) (classes []Class, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var student Student
	db.First(&student, "student_no = ?", studentNo)

	db.Model(&student).Related(&classes, "Students")
	return
}

func getClassesByTeacherNo(teacherNo string) (classes []Class, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var teacher Teacher
	db.First(&teacher, "teacher_no = ?", teacherNo)

	db.Model(&teacher).Related(&classes, "Teachers")
	return
}

func getAllDevices() (devices []Device, err error) {
	devices = []Device{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.Find(&devices)
	return
}


func getDevice(id int) (device *Device, err error) {
	device = &Device{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(device, "id = ?", id)
	return
}

func getDeviceByPath(devicePath string) (device *Device, err error) {
	device = &Device{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(device, "device_path = ?", devicePath)
	return
}

func getClassroom(id int) (classroom *Classroom, err error) {
	classroom = &Classroom{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(classroom, "id = ?", id)
	return
}

func getStudent(studentNo string) (student *Student, err error) {
	student = &Student{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(student, "student_no = ?", studentNo)
	return
}

func getStudentsByClass(classID int) (students []Student, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var class Class
	db.First(&class, "id = ?", classID)

	db.Model(&class).Related(&students, "Classes")
	return
}

func getStudents(studentNos []string) (students []Student, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&students, "student_no in (?)", studentNos)
	return
}

func getAllStudents() (students []Student, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Order("created_at desc").Find(&students)
	return
}

func getTeachers(teachersNos []string) (teachers []Teacher, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&teachers, "teacher_no in (?)", teachersNos)
	return
}

func getClassroomStatsItem(classroomID int) (stats *ClassroomStatsTable, err error) {
	stats = &ClassroomStatsTable{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Last(stats, "classroom_id = ?", classroomID)
	return
}

func getCamera(cameraID int) (camera *Camera, err error) {
	camera = &Camera{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(camera, "id = ?", cameraID)
	return
}

func getCameras() (cameras []Camera, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&cameras)
	return
}

func getClassrooms() (classrooms []Classroom, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&classrooms)
	return
}

func getDeviceManagerSystemStats() (systemStats []DeviceManagerSystemStats, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Order("created_at desc").Limit(10).Find(&systemStats)
	return
}

func getDeviceCameraCount() (devices int, cameras int, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Table("devices").Count(&devices)
	db.Table("cameras").Count(&cameras)
	return
}

func getTableCount(tableName string) (count int, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Table(tableName).Count(&count)
	return
}

func createTableItem(v interface{}) error {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Create(v)

	return nil
}

func updateTableItem(old interface{}, new map[string]interface{}) (err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Model(old).Update(new)
	return
}

func deleteTableItem(v interface{}) (err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Unscoped().Delete(v)
	return
}

func deleteTableItems(v interface{}, formatString string, value interface{}) (err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Unscoped().Where(formatString, value).Delete(v)
	return
}