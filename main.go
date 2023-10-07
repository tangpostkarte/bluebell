package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// Go web开发通用脚手架

func main() {
	// 1, 加载配置

	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err: %v", err)
		return
	}

	// 2, 初始化日志

	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err: %v", err)
		return
	}

	defer zap.L().Sync() // 把缓冲区的日志追加到日志文件

	zap.L().Debug("logger init success")

	// 3, 初始化Mysql

	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err: %v", err)
		return
	}

	// 关闭数据库
	defer mysql.Close()

	// 4, 初始化Redis

	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err: %v", err)
		return
	}

	defer redis.Close()

	// 5, 注册路由

	r := routes.SetUp()

	// 6, 启动服务（优雅关机）

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen failed", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)

	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行

	zap.L().Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown", zap.Error(err))
	}
}
