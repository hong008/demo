package main

import (
	"demo/2.gin-demo/models"
	"demo/2.gin-demo/pkg/gredis"
	"demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/routers"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/robfig/cron"
	"log"
	"syscall"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/17 23:23 
    @File: blog.go
*/

func main() {
	//TODO swagger、endless

	//通过http.Server - Shutdown() 实现热更新
	/*router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logging.Info(fmt.Sprintf("Listen: %s", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Info(fmt.Sprintf("Shutdown Server..."))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logging.Fatal(fmt.Sprintf("Server Shutdown: %v", err))
	}
	logging.Info(fmt.Sprintf("Server exiting"))*/

	//加载配置项
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	//第三方库endless实现热更新
	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		logging.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}

	/*go func() {
		clean()
	}()*/

	err := server.ListenAndServe()
	if err != nil {
		logging.Error(fmt.Sprintf("Server err:%v", err))
	}
}

func clean() {
	log.Println("Starting...")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
