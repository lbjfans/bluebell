package main

import (
	"backend/controller"
	"backend/dao/mysql"
	"backend/dao/redis"
	"backend/logger"
	"backend/pkg/snowflake"
	"backend/routers"
	"backend/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 用makefile运行
	//var confFile string
	//flag.StringVar(&confFile, "conf", "./conf/config.yaml", "配置文件")
	//flag.Parse()
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init settings fail", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init log fail", err)
		return
	}
	defer zap.L().Sync() // 内存 > 磁盘
	// 3. 初始化mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init mysql fail", err)
		return
	}
	defer mysql.Close()
	// 4. redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis fail", err)
		return
	}
	defer redis.Close()
	// 雪花算法生成分布式ID
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 注册使用：validator
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator Trans failed,err:%v\n", err)
		return
	}
	// 5. 注册路由
	r := routers.Setup(settings.Conf.Mode)
	// 6. 启动服务，优雅关机
	srv := &http.Server{
		//Addr:    ":8080",
		Addr:    fmt.Sprintf("localhost:%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
