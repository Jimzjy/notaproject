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

func getDevice(devicePath string) (device *Device, err error) {
	device = &Device{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(device, "device_path = ?", devicePath)
	return
}

func getClassroom(name string) (classroom *Classroom, err error) {
	classroom = &Classroom{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	db.First(classroom, "name = ?", name)
	return
}

func getClassroomStatsItem(classroomName string) (stats *ClassroomStatsTable, err error) {
	stats = &ClassroomStatsTable{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return
	}
	defer db.Close()

	classroom, err := getClassroom(classroomName)
	if err != nil {
		return
	}

	db.Last(stats, "classroom_id = ?", classroom.ID)
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