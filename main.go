package main

import (
	"fmt"
	"remote_monitor/monitor"
	"time"
)

func main() {
	fmt.Printf("%+v\n", monitor.GetStaticPCData())

	for true {
		fmt.Printf("%+v\n", monitor.GetDynamicPCData())
		time.Sleep(1000 * time.Millisecond)
	}
}
