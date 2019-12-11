package main

import (
	"fmt"
	"gin_demo/models"
	"gin_demo/pkg/logging"
	"gin_demo/pkg/setting"
	"gin_demo/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {
	//r := routers.InitRouter()
	//
	//s := http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
	//	Handler:        r,
	//	TLSConfig:      nil,
	//	ReadTimeout:    setting.ReadTimeout,
	//	WriteTimeout:   setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}

	setting.Setup()
	models.Setup()
	logging.Setup()

	//fmt.Println(	setting.AppSetting)
	//
	//fmt.Println(setting.ServerSetting)
	//
	//fmt.Println(setting.DatabaseSetting)

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeOut
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeOut
	endless.DefaultMaxHeaderBytes = 1 << 20

	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())

	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
