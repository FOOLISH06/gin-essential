package controller

import (
	"github.com/foolish06/gin-essential/common"
	_ "github.com/foolish06/gin-essential/common"
	"github.com/foolish06/gin-essential/model"
	"github.com/foolish06/gin-essential/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// 创建用户
	newUser := &model.User{Name: name, Telephone: telephone, Password: password}
	db.Create(&newUser)

	// 返回结果
	ctx.JSON(200, gin.H{
		"message": "注册成功",
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