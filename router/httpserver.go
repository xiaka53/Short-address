package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/DeployAndLog/lib"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

//启动服务
func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		log.Printf(" [INFO] HttpServerRun%s\n", lib.GetStringConf("base.http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun%s err:%v\n", lib.GetStringConf("base.http.addr"), err)
		}
	}()
}

//关闭服务
func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stoppend\n")
}
