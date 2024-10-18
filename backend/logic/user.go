package logic

import (
	"backend/dao/mysql"
	"backend/models"
	"backend/pkg/jwt"
	"backend/pkg/snowflake"
)

func SignUp(p *models.RegisterForm) (error error) {
	// 1. 判断用户是否存在，存在就不能注册
	err := mysql.CheckUserExist(p.UserName)
	if err != nil {
		return err
	}
	// 2. 生成uid
	userId, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}
	// 3. 放到数据库
	// 构造一个User实例
	u := models.User{
		UserID:   userId,
		UserName: p.UserName,
		Password: p.Password,
		Email:    p.Email,
		Gender:   p.Gender,
	}
	return mysql.InsertUser(u)
}

func Login(p *models.LoginForm) (user *models.User, error error) {
	// 1. 从数据库查找用户，如果不存在，返回错误
	// 2. 如果密码不正确，返回错误
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 3. JWT
	accessToken, refreshToken, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return
}
