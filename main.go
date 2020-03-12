package main

import (
	"github.com/xiaka53/DeployAndLog/lib"
	"github.com/xiaka53/Short-address/public"
	"github.com/xiaka53/Short-address/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := lib.InitModule("./conf/dev/", []string{"base", "redis"}); err != nil {
		log.Fatal(err)
	}
	type at struct {
		Ac map[string]string
	}
	defer lib.Destroy()
	public.InitValidate()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
