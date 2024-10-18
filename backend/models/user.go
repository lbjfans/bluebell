package models

// User 定义请求参数结构体
type User struct {
	UserID       uint64 `json:"user_id,string" db:"user_id"` // 指定json序列化/反序列化时使用小写user_id
	UserName     string `json:"username" db:"username"`      // 如db.Get(&user, "SELECT * FROM users WHERE user_id=?", id)会用到db
	Password     string `json:"password" db:"password"`
	Email        string `json:"email" db:"email"`   // 邮箱
	Gender       int    `json:"gender" db:"gender"` // 性别
	AccessToken  string
	RefreshToken string
}

// RegisterForm 注册请求参数
type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`                          // 用户名
	Email           string `json:"email" binding:"required"`                             // 邮箱
	Gender          int    `json:"gender" binding:"oneof=0 1 2"`                         // 性别 0:未知 1:男 2:女
	Password        string `json:"password" binding:"required"`                          // 密码
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"` // eqfield查询validator
}

// LoginForm 登录请求参数
type LoginForm struct {
	UserName string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}
