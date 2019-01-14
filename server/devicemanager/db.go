package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"reflect"

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
	db.Order("created_at desc").Find(&classes)

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
	if class.ID == 0 {
		err = fmt.Errorf("can not find class for id %v", id)
		return
	}
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

	db.Model(&student).Related(&classes, "Classes")
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

	db.Model(&teacher).Related(&classes, "Classes")
	return
}

func getAllDevices() (devices []Device, err error) {
	devices = []Device{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.Order("created_at desc").Find(&devices)
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
	if device.ID == 0 {
		err = fmt.Errorf("can not find device for id %v", id)
		return
	}
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
	if device.ID == 0 {
		err = fmt.Errorf("can not find device for path %v", devicePath)
		return
	}
	return
}

func getDevicesByCamera(cameraID int) (devices []Device, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var camera Camera
	db.First(&camera, "id = ?", cameraID)

	db.Model(&camera).Related(&devices, "Devices")
	return
}

func getClassroom(no string) (classroom *Classroom, err error) {
	classroom = &Classroom{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(classroom, "classroom_no = ?", no)
	if classroom.ID == 0 {
		err = fmt.Errorf("can not find classroom for no %v", no)
	}
	return
}

func getClassroomByID(id int) (classroom *Classroom, err error) {
	classroom = &Classroom{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(classroom, "id = ?", id)
	if classroom.ID == 0 {
		err = fmt.Errorf("can not find classroom for no %v", id)
	}
	return
}

func getAllClassrooms() (classrooms []Classroom, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Order("created_at desc").Find(&classrooms)
	return
}

func getClassroomsByCamera(cameraID int) (classrooms []Classroom, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var camera Camera
	db.First(&camera, "id = ?", cameraID)

	db.Model(&camera).Related(&classrooms, "Classrooms")
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
	if student.ID == 0 {
		err = fmt.Errorf("can not find student for no %v", studentNo)
	}
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

	db.Model(&class).Related(&students, "Students")
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

	db.Order("created_at desc").Order("created_at desc").Find(&students)
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

func getTeachersByClass(classID int) (teachers []Teacher, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var class Class
	db.First(&class, "id = ?", classID)

	db.Model(&class).Related(&teachers, "Teachers")
	return
}

func getTeacher(no string) (teacher *Teacher, err error) {
	teacher = &Teacher{}

	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.First(teacher, "teacher_no = ?", no)
	if teacher.ID == 0 {
		err = fmt.Errorf("can not find teacher for no %v", no)
	}
	return
}

func getAllTeachers() (teachers []Teacher, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Order("created_at desc").Find(&teachers)
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
	if stats.ID == 0 {
		err = fmt.Errorf("can not find classroom stats for classroom id %v", classroomID)
	}
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
	if camera.ID == 0 {
		err = fmt.Errorf("can not find camera for id %v", cameraID)
	}
	return
}

func getCameras(cameraIDs []int) (cameras []Camera, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&cameras, "id in (?)", cameraIDs)
	return
}

func getCamerasByDevice(deviceID int) (cameras []Camera, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var device Device
	db.First(&device, "id = ?", deviceID)

	db.Model(&device).Related(&cameras, "Cameras")
	return
}

func getCamerasByClassroom(classroomNo string) (cameras []Camera, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	var classroom Classroom
	db.First(&classroom, "classroom_no = ?", classroomNo)

	db.Model(&classroom).Related(&cameras, "Cameras")
	return
}

func getAllCameras() (cameras []Camera, err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	db.Find(&cameras)
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

func updateTableItem(old, new interface{}) (err error) {
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

func updateAssociation(value interface{}, foreignKey string, replace interface{}) (err error) {
	db, err := gorm.Open(DB, DBName)
	if err != nil {
		return
	}
	defer db.Close()

	_replace := reflect.ValueOf(replace)
	if _replace.Kind() != reflect.Slice {
		err = errors.New("not slice")
		return
	}

	replaceSlice := make([]interface{}, _replace.Len())
	for i := 0; i < _replace.Len(); i++ {
		replaceSlice[i] = _replace.Index(i).Interface()
	}

	db.Model(value).Association(foreignKey).Clear()
	for _, v := range replaceSlice {
		db.Model(value).Association(foreignKey).Append(v)
	}
	return
}