package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"varchar(20);not null"`
	Telephone string `gorm:"varchar(20);not null; unique"`
	Password string `gorm:"varchar(20);not null"`
}

func main() {
	db := InitDB()

	router := gin.Default()
	_ = router.SetTrustedProxies([]string{"localhost"})
	router.POST("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		// 数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能小于6位"})
			return
		}

		// 如果用户名为空，返回一个10位的随机字符串
		if len(name) == 0 {
			name = RandString(10)
		}

		// 判断手机号是否存在
		if isTelephoneExit(db, telephone) {
			ctx.JSON(422, gin.H{
				"msg": "手机号已存在",
			})
			return
		}

		// 创建用户
		newUser := &User{Name: name, Telephone: telephone, Password: password}
		db.Create(&newUser)

		// 返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})

	if err := router.Run(); err != nil {
		log.Fatalln(err.Error())
	} // listen and serve on 0.0.0.0:8080
}

func RandString(n int) string {
	letters := "qwfeuhfouidnviwueqhvupiqiducvhduovqwhewuie"
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user *User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func InitDB() *gorm.DB {
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


	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalln("fail to migrate, err: ", err.Error())
	}

	return db
}