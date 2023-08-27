package middleware

import (
	"douyin/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var MySecret = []byte("0xCAFEBABY")

type MyClaims struct {
	UserId int64
	jwt.RegisteredClaims
}

func CreateToken(user models.UserLogin) (string, error) {
	claim := MyClaims{
		UserId: user.UserInfoId,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 失效
			Issuer:    "LXJ",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	s, err := token.SignedString(MySecret)
	if err != nil {
		return "", err
	}
	return s, nil
}

func ParseToken(s string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(s, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("claim invalid")
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}
	return claims, nil
}

func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 401,
				StatusMsg:  "该用户不存在",
			})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if tokenStruck.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusOK, models.CommonResponse{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}
