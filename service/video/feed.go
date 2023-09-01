package video

import (
	"douyin/models"
	"douyin/util"
	"time"
)

var MaxVideoNum = 30

func QueryFeedVideoList(userId int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

type FeedVideoList struct {
	Videos   []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

type QueryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time

	videos    []*models.Video
	nextTime  int64
	feedVideo *FeedVideoList
}

func NewQueryFeedVideoListFlow(userId int64, latestTime time.Time) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

func (q *QueryFeedVideoListFlow) Do() (*FeedVideoList, error) {
	q.checkNum()

	if err := q.prepareData(); err != nil {
		return nil, err
	}
	q.packData()
	return q.feedVideo, nil
}

func (q *QueryFeedVideoListFlow) checkNum() {
	if q.userId > 0 {
		//	todo
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}

func (q *QueryFeedVideoListFlow) prepareData() error {

	if err := models.NewVideoDAO().QueryVideoListByLimitAndTime(MaxVideoNum, q.latestTime, &q.videos); err != nil {
		return err
	}
	latestTime, _ := util.FillVideoListField(q.userId, &q.videos)

	if latestTime != nil {
		q.nextTime = (*latestTime).UnixNano()
		return nil
	}
	q.nextTime = time.Now().Unix() / 1e6
	return nil
}

func (q *QueryFeedVideoListFlow) packData() {
	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
}
