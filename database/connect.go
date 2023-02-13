package database

import (
	"fmt"
	"ginfo/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var err error

func ConnectDB() {
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&model.Product{}, &model.User{})
	fmt.Println("Database Migrated")
}
