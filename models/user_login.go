package models

import (
	"errors"
	"sync"
)

type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

func (UserLogin) TableName() string {
	return `user_login`
}

type UserLoginDAO struct {
}

var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})
	return userLoginDao
}
func (u *UserLoginDAO) QueryUserLogin(username, password string, login *UserLogin) error {
	if login == nil {
		return errors.New("结构体指针为空")
	}
	return DB.Where("username=? and password=?", username, password).First(login).Error
}

func (u *UserLoginDAO) IsUserExistByUsername(username string) bool {
	var count int64
	DB.Model(&UserLogin{}).Where("username=?", username).Count(&count)
	return count > 0
}
