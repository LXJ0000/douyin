package router

import (
	"douyin/controller/user_info"
	"douyin/controller/user_login"
	"douyin/controller/video"
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

	baseGroup.GET("/feed/", video.FeedVideoListHandler)

	//用户信息
	baseGroup.GET("/user/", middleware.JWTMiddleWare(), user_info.UserInfoHandler)
	//登录
	baseGroup.POST("/user/login/", middleware.SHAMiddleWare(), user_login.UserLoginHandler)
	//注册
	baseGroup.POST("/user/register/", middleware.SHAMiddleWare(), user_login.UserRegisterHandler)
	//视频投稿
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare(), video.PublishVideoHandler)
	//视频列表
	baseGroup.GET("/publish/list/", middleware.JWTMiddleWare(), video.QueryVideoListHandler)

	return r
}
