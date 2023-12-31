package user_login

import (
	middleware "douyin/middlewares"
	"douyin/models"
	"errors"
	"log"
)

type QueryUserLoginFlow struct {
	username string
	password string
	userid   int64
	token    string
	data     *models.LoginResponse
}

func QueryUserLogin(username, password string) (*models.LoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

func (q *QueryUserLoginFlow) Do() (*models.LoginResponse, error) {

	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//准备好数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包最终数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *QueryUserLoginFlow) checkNum() error {
	if q.username == "" {
		return errors.New("username is null")
	}
	if len(q.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码不能为空")
	}
	return nil
}
func (q *QueryUserLoginFlow) prepareData() error {
	userLoginDAO := models.NewUserLoginDao()
	var login models.UserLogin

	if err := userLoginDAO.QueryUserLogin(q.username, q.password, &login); err != nil {
		return err
	}
	q.userid = login.UserInfoId

	token, err := middleware.CreateToken(login)
	if err != nil {
		return err
	}
	q.token = token
	return nil
}
func (q *QueryUserLoginFlow) packData() error {
	q.data = &models.LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	log.Println(q.token)
	return nil
}
