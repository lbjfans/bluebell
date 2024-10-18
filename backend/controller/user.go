package controller

import (
	"backend/dao/mysql"
	"backend/logic"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 注册业务
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数：用户名，密码，再次输入的密码；参数检验
	// 前端会发送json，用结构体接收(models/user)
	var registerUser *models.RegisterForm // 用指针接收，因为结构体是值类型，传递开销大
	// 检测输入是否为json格式，以及每个字段的数据类型是否正确
	if err := c.ShouldBindJSON(&registerUser); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors) // 判断输入是否为空，密码和确认密码是否相同；前后端都需要校验
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译，回给前端
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	fmt.Printf("registerUser: %v\n", registerUser)
	// 2. 业务处理
	if err := logic.SignUp(registerUser); err != nil {
		zap.L().Error("logic.signup failed", zap.Error(err))
		if err.Error() == mysql.ErrorUserExit { // 用户已经存在
			ResponseError(c, CodeUserExist)
			return
		}
		// 其他错误
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应，可以封装一个返回中间件
	ResponseSuccess(c, nil)
}

// LoginHandler 登录业务
func LoginHandler(c *gin.Context) {
	// 1. 前端获取参数（结构体）：用户名，密码；参数校验
	var loginUser *models.LoginForm
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors) // 判断输入是否为空；前后端都需要校验
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译，回给前端
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	fmt.Printf("loginUser: %v\n", loginUser)
	// 2. 业务处理
	user, err := logic.Login(loginUser)
	if err != nil {
		zap.L().Error("logic.login failed", zap.String("username", loginUser.UserName), zap.Error(err))
		// 2个错误：用户不存在，密码错误
		if err.Error() == mysql.ErrorUserNotExit || err.Error() == mysql.ErrorPasswordWrong {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		// 其他错误
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应，可以封装一个返回中间件
	ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID), //js识别的最大值：id值大于1<<53-1  int64: i<<63-1
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}
