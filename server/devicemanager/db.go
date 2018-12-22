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