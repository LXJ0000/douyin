package router

import (
	"douyin/controller/user_info"
	"douyin/controller/user_login"
	middleware "douyin/middlewares"
	"douyin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitDouyinRouter() *gin.Engine {
	models.InitDB()

	r := gin.Default()

	//设置静态文件夹 存储视频图片
	r.Static("static", "./static")

	r.GET("/ping/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"ping": "success",
		})
	})

	//设置统一入口
	baseGroup := r.Group("/douyin")

	//用户信息
	baseGroup.GET("/user/", middleware.JWTMiddleWare(), user_info.UserInfoHandler)
	//登录
	baseGroup.POST("/user/login/", middleware.SHAMiddleWare(), user_login.UserLoginHandler)
	//注册
	baseGroup.POST("/user/register/", middleware.SHAMiddleWare(), user_login.UserRegisterHandler)

	return r
}