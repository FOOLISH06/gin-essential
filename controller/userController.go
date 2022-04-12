package controller

import (
	"github.com/foolish06/gin-essential/common"
	_ "github.com/foolish06/gin-essential/common"
	"github.com/foolish06/gin-essential/model"
	"github.com/foolish06/gin-essential/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()

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
		name = util.RandString(10)
	}

	// 判断手机号是否存在
	if isTelephoneExit(db, telephone) {
		ctx.JSON(422, gin.H{
			"msg": "手机号已存在",
		})
		return
	}

	// 创建用户（需要加密）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
	}

	newUser := &model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hashedPassword),
	}
	db.Create(&newUser)

	// 返回结果
	ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	db := common.GetDB()

	// 获取参数
	telephone := ctx.PostForm("telephone")
	password  := ctx.PostForm("password")

	// 参数验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能小于6位"})
		return
	}

	// 判断用户（手机号）是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}


	// 密码验证
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error: %v\n", err)
		return
	}

	// 返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg": "登录成功",
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": user},
	})
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user *model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}