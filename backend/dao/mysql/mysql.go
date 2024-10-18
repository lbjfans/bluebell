package mysql

import (
	"backend/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 匿名import，调用init(在driver.go)，调用mysql的驱动，用于连接mysql
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	// dsn := "root:666666@tcp(192.168.100.4:3306)/bubble"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB,
	)
	db, err = sqlx.Connect("mysql", dsn) // 返回*DB，先有驱动才能连接
	if err != nil {
		zap.L().Error("connect mysql fail", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_cons"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_cons"))
	return
}

func Close() {
	db.Close()
}
