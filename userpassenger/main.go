package main

import (
	"gotrains/userpassenger/initialize"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	go func() {
		// run server
		// initialize.InitServer()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// deregister service
}
