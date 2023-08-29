package models

import (
	"errors"
	"sync"
	"time"
)

// Video `gorm:"-"`  读写操作均会忽略该字段
type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	//Comments      []*Comment  `json:"-"` todo 添加评论功能后取消注释
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type VideoDAO struct {
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAO)
	})
	return videoDAO
}

func (v *VideoDAO) AddVideo(video *Video) error {
	if video == nil {
		return errors.New("AddVideo video 空指针")
	}
	return DB.Create(video).Error
}
func (v *VideoDAO) QueryVideoByVideoId(videoId int64, video *Video) error {
	if video == nil {
		return errors.New("QueryVideoByVideoId 空指针")
	}
	return DB.Where("id=?", videoId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		First(video).Error
	// todo modify select
}

func (v *VideoDAO) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return DB.Model(&Video{}).Where("user_info_id=?", userId).Count(count).Error
}
func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return DB.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error
	// todo
}

// QueryVideoListByLimitAndTime  返回按投稿时间倒序的视频列表，并限制为最多limit个
func (v *VideoDAO) QueryVideoListByLimitAndTime(limit int, latestTime time.Time, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByLimit videoList 空指针")
	}
	return DB.Model(&Video{}).Where("created_at<?", latestTime).
		Order("created_at ASC").Limit(limit).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).
		Find(videoList).Error
}

// IsVideoExistById 视频是否存在
func (v *VideoDAO) IsVideoExistById(id int64) bool {
	var count int64
	DB.Model(&Video{}).Where("id=?", id).Count(&count)
	return count > 0
}
