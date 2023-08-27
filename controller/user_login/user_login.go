package user_login

import (
	"douyin/models"
	"douyin/service/user_login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	models.CommonResponse
	*models.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")

	password, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析出错",
			},
		})
	}
	userLoginResponse, err := user_login.QueryUserLogin(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		LoginResponse: userLoginResponse,
	})
}
