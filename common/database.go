package common

import (
	"fmt"
	"github.com/foolish06/gin-essential/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		db = initDB()
	}
	return db
}

func initDB() *gorm.DB {
	username := "root"
	password := "123456"
	host	 := "localhost"
	port 	 := "3306"
	database := "tmp"
	charset	 := "utf8"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("fail to connect to database, err:v", err.Error())
	}


	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalln("fail to migrate, err: ", err.Error())
	}

	return db
}


