package main

import (
	"fmt"
	"os"
	"time"

	"github.com/luisya22/downmonitor/monitor"
)

func main() {
	fmt.Println("Starting")
	go func() {

		for {
			file, err := os.Open("./downmonitor.log")
			if err != nil {
				fmt.Printf("error opening file: %v", err.Error())
				return
			}
			defer file.Close()

			err = monitor.QueryData(file)
			if err != nil {
				fmt.Printf("error querying: %s\n", err.Error())
			}

			time.Sleep(15 * time.Second)
		}
	}()

	monitor.Start()
}
