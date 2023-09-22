package main

import (
	"fmt"

	"github.com/luisya22/downmonitor/monitor"
	"github.com/luisya22/downmonitor/server"
)

func main() {
	fmt.Println("Starting")

	go func() {
		server.Start()
	}()

	monitor.Start()
}
