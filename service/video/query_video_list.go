package video

import (
	"douyin/models"
	"errors"
)

// List 视频集合返回
type List struct {
	Videos []*models.Video `json:"video_list,omitempty"`
}

func QueryVideoList(userId int64) (*List, error) {
	return NewQueryVideoListFlow(userId).Do()
}

type QueryVideoListFlow struct {
	userId int64
	videos []*models.Video

	videoList *List
}

func NewQueryVideoListFlow(userId int64) *QueryVideoListFlow {
	return &QueryVideoListFlow{userId: userId}
}

func (q *QueryVideoListFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryVideoListFlow) checkNum() error {
	//检查userId是否存在
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return errors.New("用户不存在")
	}
	return nil
}

func (q *QueryVideoListFlow) packData() error {
	if err := models.NewVideoDAO().QueryVideoListByUserId(q.userId, &q.videos); err != nil {
		return err
	}

	var userinfo models.UserInfo
	if err := models.NewUserInfoDAO().QueryUserInfoById(q.userId, &userinfo); err != nil {
		return err
	}
	for i, _ := range q.videos {
		q.videos[i].Author = userinfo
	}
	q.videoList = &List{Videos: q.videos}
	return nil
}
