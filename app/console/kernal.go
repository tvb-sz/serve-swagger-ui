package console

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jjonline/go-lib-backend/logger"
	"github.com/tvb-sz/serve-swagger-ui/client"
	"github.com/tvb-sz/serve-swagger-ui/conf"
	"github.com/tvb-sz/serve-swagger-ui/route"
	"github.com/tvb-sz/serve-swagger-ui/utils"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// region 服务引导启动器

// BootStrap 引导启动
func BootStrap() {
	// output base info
	client.Logger.Info(fmt.Sprintf("Golang Version      : %s", runtime.Version()))
	client.Logger.Info(fmt.Sprintf("MAX Cpu Num         : %d", runtime.GOMAXPROCS(-1)))
	client.Logger.Info(fmt.Sprintf("Command Args        : %s", os.Args))
	client.Logger.Info(fmt.Sprintf("Config file Path    : %s", conf.Config.ConfigFile))
	client.Logger.Info(fmt.Sprintf("Enable Google login : %s", conf.Config.EnableGoogle))
	client.Logger.Info(fmt.Sprintf("Server HOST         : %s", conf.Config.Server.Host))
	client.Logger.Info(fmt.Sprintf("Server PORT         : %s", conf.Config.Server.Port))
	client.Logger.Info(fmt.Sprintf("Swagger JSON Path   : %s", conf.Config.Swagger.Path))
	client.Logger.Info(fmt.Sprintf("Log Path            : %s", conf.Config.Log.Path))
	client.Logger.Info(fmt.Sprintf("Log Level           : %s", conf.Config.Log.Level))

	_, signalChan := quitCtx()

	startHTTPApp(signalChan)
}

// endregion

// region HttpServer启动

// startHTTPApp http api服务启动
func startHTTPApp(signalChan chan os.Signal) {
	// 设置gin启动模式
	gin.SetMode(utils.RunMode())

	// 设置gin路由注册日志输出func
	// 注意：gin路由注册日志输出仅dev模式才有
	gin.DebugPrintRouteFunc = logger.GinPrintInitRoute

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Config.Server.Port),
		Handler:        route.Bootstrap(),
		ReadTimeout:    time.Duration(conf.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(conf.Config.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, //1MB
	}

	// main主进程阻塞channel
	idleCloser := make(chan struct{})

	// http serv handle exit signal
	go func() {
		<-signalChan

		// 超时context
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer timeoutCancel()
		// We received an interrupt signal, shut down.
		if err := server.Shutdown(timeoutCtx); err != nil {
			// Error from closing listeners, or context timeout:
			client.Logger.Error("Http服务暴力停止：" + err.Error())
		} else {
			// time.Sleep(1 * time.Second)
			// successful shutdown process ok
			client.Logger.Info("Http服务优雅停止")
		}

		// closer quit chan
		close(idleCloser)
	}()

	// continue serv http service
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		client.Logger.Error("Http服务异常：" + err.Error())
		close(idleCloser)
	}

	// wait for stop main process
	<-idleCloser
	client.Logger.Info("进程已退出：服务已关闭")
}

// endregion

// region 全局进程控制信号捕获

// quitCtx 全局退出信号
func quitCtx() (context.Context, chan os.Signal) {
	ctx, cancel := context.WithCancel(context.Background())
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		<-quitChan // wait quit signal

		signal.Stop(quitChan)
		cancel()
		close(quitChan)
	}()

	return ctx, quitChan
}

// endregion