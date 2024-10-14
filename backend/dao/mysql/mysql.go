package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 匿名import，调用init
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	// dsn := "root:666666@tcp(192.168.100.4:3306)/bubble"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	db, err = sqlx.Connect("mysql", dsn) // 返回*DB，没有真正连接
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
