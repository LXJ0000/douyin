package video

import (
	"douyin/models"
	"douyin/service/video"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListResponse struct {
	models.CommonResponse
	*video.List
}

func QueryVideoListHandler(c *gin.Context) {
	p := NewProxyQueryVideoList(c)
	rawId, _ := c.Get("user_id")
	err := p.DoQueryVideoListByUserId(rawId)
	if err != nil {
		p.c.JSON(http.StatusOK, models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
}

// NewProxyQueryVideoList 封装一层
func NewProxyQueryVideoList(c *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{c: c}
}

// ProxyQueryVideoList  context代理类
type ProxyQueryVideoList struct {
	c *gin.Context
}

// DoQueryVideoListByUserId 根据userId字段进行查询
func (p *ProxyQueryVideoList) DoQueryVideoListByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	videoList, err := video.QueryVideoList(userId)
	if err != nil {
		return err
	}
	//json 返回到gin
	p.c.JSON(http.StatusOK, ListResponse{
		models.CommonResponse{StatusCode: 0},
		videoList,
	})
	return nil
}
