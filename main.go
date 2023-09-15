package main

import (
	"fmt"

	"github.com/luisya22/downmonitor/monitor"
	"github.com/luisya22/downmonitor/server"
)

func main() {
	fmt.Println("Starting")
	// go func() {
	//
	// 	for {
	// 		file, err := os.Open("./downmonitor.log")
	// 		if err != nil {
	// 			fmt.Printf("error opening file: %v", err.Error())
	// 			return
	// 		}
	// 		defer file.Close()
	//
	// 		clock := &monitor.RealTime{}
	//
	// 		res, err := monitor.QueryData(file, clock)
	// 		if err != nil {
	// 			fmt.Printf("error querying: %s\n", err.Error())
	// 		}
	//
	// 		fmt.Printf("Today: Average downtime: %d seconds, Downtimes: %d\n", res.Today.Avg, res.Today.Amount)
	// 		fmt.Printf("Last 7 days: Average daily downtime: %d seconds, Downtimes: %d\n", res.Days7.Avg, res.Days7.Amount)
	// 		fmt.Printf("Last 30 days: Average daily downtime: %d seconds, Downtimes: %d\n", res.Days30.Avg, res.Days30.Amount)
	// 		fmt.Printf("---------\n\n----------\n")
	//
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }()

	go func() {
		server.Start()
	}()

	monitor.Start()
}
