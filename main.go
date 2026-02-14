package main

import (
	"fmt"
	"remote_monitor/monitor"
)

func main() {
	fmt.Printf("%+v\n", monitor.GetPCData())
}
