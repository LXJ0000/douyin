package models

import (
	"errors"
	"sync"
)

type UserInfo struct {
	Id            int64  `json:"id" gorm:"id,omitempty"`
	Name          string `json:"name" gorm:"name,omitempty"`
	User          *UserLogin
	FollowCount   int64 `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64 `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool  `json:"is_follow" gorm:"is_follow,omitempty"`

	//user_relations
	Follows []*UserInfo `json:"-" gorm:"many2many:user_relations;"` //用户之间的多对多
}

func (UserInfo) TableName() string {
	return `user_info`
}

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

// QueryUserInfoById 查询用户信息通过id
func (u *UserInfoDAO) QueryUserInfoById(userId int64, userinfo *UserInfo) error {
	if userinfo == nil {
		return errors.New("userinfo is nil")
	}
	//DB.Where("id=?",userId).First(userinfo)
	err := DB.Where("id=?", userId).Select([]string{"id", "name"}).First(userinfo).Error
	//id为零值，说明sql执行失败
	if err != nil {
		return errors.New("该用户不存在")
	}
	return nil
}

// AddUserInfo 添加用户并返回用户信息
func (u *UserInfoDAO) AddUserInfo(userinfo *UserInfo) error {
	if userinfo == nil {
		return errors.New("userinfo is nil")
	}
	return DB.Create(userinfo).Error
}

// IsUserExistById 判断用户是否存在
func (u *UserInfoDAO) IsUserExistById(id int64) bool {
	var count int64
	DB.Model(&UserInfo{}).Where("id=?", id).Count(&count)
	return count > 0
}
