package mysql

import (
	"backend/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 用于从数据库获取数据，返回给logic层
const secret = "tim.lbjfans"

// encryptPassword 对密码进行加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (error error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	// 获取过程中出错
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	// 已经存在，返回错误
	if count > 0 {
		return errors.New(ErrorUserExit)
	}
	// 没问题
	return
}

// InsertUser 注册业务-向数据库中插入一条新的用户
func InsertUser(user models.User) (error error) {
	// 1. 对密码加密
	user.Password = encryptPassword([]byte(user.Password))
	// 2. 放入数据库
	sqlStr := `insert into user(user_id,username,password,email,gender) values(?,?,?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.UserName, user.Password, user.Email, user.Gender)
	return err
}

func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录一下原始密码(用户登录的密码)
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	// 查询数据库出错
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	// 用户不存在：有的网站只显示用户/密码错误，防止对用户名的查找
	if err == sql.ErrNoRows {
		return errors.New(ErrorUserNotExit)
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return errors.New(ErrorPasswordWrong)
	}
	return nil
}
