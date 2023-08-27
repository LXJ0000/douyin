package user_login

import (
	middleware "douyin/middlewares"
	"douyin/models"
	"errors"
)

// 常量
const (
	MaxUsernameLength = 100
	MaxPasswordLength = 20
	MinPasswordLength = 8
)

func PostUserLogin(username, password string) (*models.LoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Do()
}

type PostUserLoginFlow struct {
	username string
	password string
	userid   int64
	token    string
	data     *models.LoginResponse
}

func NewPostUserLoginFlow(username, password string) *PostUserLoginFlow {
	return &PostUserLoginFlow{username: username, password: password}
}

func (q *PostUserLoginFlow) Do() (*models.LoginResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.updateData(); err != nil {
		return nil, err
	}
	if err := q.packResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *PostUserLoginFlow) checkNum() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if len(q.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (q *PostUserLoginFlow) updateData() error {
	userLogin := models.UserLogin{Username: q.username, Password: q.password}
	userinfo := models.UserInfo{User: &userLogin, Name: q.username}

	userLoginDao := models.NewUserLoginDao()
	if userLoginDao.IsUserExistByUsername(q.username) {
		return errors.New("用户名已存在")
	}

	userInfoDao := models.NewUserInfoDAO()

	if err := userInfoDao.AddUserInfo(&userinfo); err != nil {
		return err
	}
	token, err := middleware.CreateToken(userLogin)
	if err != nil {
		return err
	}
	q.token = token
	q.userid = userinfo.Id
	return nil
}

func (q *PostUserLoginFlow) packResponse() error {
	q.data = &models.LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil

}
