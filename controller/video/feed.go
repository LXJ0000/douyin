package video

import (
	middleware "douyin/middlewares"
	"douyin/models"
	"douyin/service/video"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	models.CommonResponse
	*video.FeedVideoList
}

func FeedVideoListHandler(c *gin.Context) {
	p := NewProxyFeedVideoList(c)
	token, ok := c.GetQuery("token")
	if !ok {

		if err := p.NoToken(); err != nil {
			p.FeedVideoListError(err.Error())
		}
		return
	}
	if err := p.HasToken(token); err != nil {
		p.FeedVideoListError(err.Error())
	}

}

type ProxyFeedVideoList struct {
	*gin.Context
}

func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Context: c}
}

func (p *ProxyFeedVideoList) NoToken() error {
	rawTimeStamp := p.Query("latest_time")
	var latestTime time.Time

	if intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64); err != nil {
		latestTime = time.Unix(0, intTime*1e6)
	}

	videoList, err := video.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

func (p *ProxyFeedVideoList) HasToken(token string) error {
	if claim, ok := middleware.ParseToken(token); ok == nil {
		if claim.ExpiresAt.Time.Before(time.Now()) {
			return errors.New("token超时")
		}
		rawTimestamp := p.Query("latest_time")
		var latestTime time.Time

		if intTime, err := strconv.ParseInt(rawTimestamp, 10, 64); err != nil {
			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
		}
		//调用service层接口
		videoList, err := video.QueryFeedVideoList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		p.FeedVideoListOk(videoList)
		return nil
	}
	return errors.New("token error")
}

// FeedVideoListError 转为json返回给客户端
func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
	p.JSON(http.StatusOK, FeedResponse{CommonResponse: models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

// FeedVideoListOk 转为json返回给客户端
func (p *ProxyFeedVideoList) FeedVideoListOk(videoList *video.FeedVideoList) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	},
	)
}
