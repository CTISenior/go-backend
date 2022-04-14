package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	Pkg "iotwin/main/pkg"
)

func main() {
	fmt.Println("***Main***")

	// Init
	Pkg.InitLogger()
	Pkg.InitHelper()

	// MQTT
	Pkg.CreateMQTTClient()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}
