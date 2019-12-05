package main

import (
	"fmt"
	"gin_demo/pkg/setting"
	"gin_demo/routers"
	"net/http"
)

func main() {
	r := routers.InitRouter()

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		TLSConfig:      nil,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
