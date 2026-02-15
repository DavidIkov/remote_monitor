package monitor

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type StaticCPUData struct {
	Name        string
	CoresAmount int
	Arch        string
}

type StaticPCData struct {
	CPU      StaticCPUData
	Hostname string
	Username string
}

type DynamicPCData struct {
	CPUUsage []float64
	RAM      mem.VirtualMemoryStat
}

func GetStaticPCData() StaticPCData {
	hostname, hostname_err := os.Hostname()
	if hostname_err != nil {
		fmt.Println("Error: ", hostname_err)
	}
	username, username_err := user.Current()
	if username_err != nil {
		fmt.Println("Error: ", username_err)
	}
	cpuInfo, cpu_err := cpu.Info()
	if cpu_err != nil {
		fmt.Println("Error: ", cpu_err)
	}

	return StaticPCData{
		CPU: StaticCPUData{
			Name:        cpuInfo[0].ModelName,
			CoresAmount: runtime.NumCPU(),
			Arch:        runtime.GOARCH},
		Hostname: hostname,
		Username: username.Username}
}

func GetDynamicPCData() DynamicPCData {
	cpu_usage, cpu_usage_err := cpu.Percent(0, true)
	if cpu_usage_err != nil {
		fmt.Println("Error: ", cpu_usage_err)
	}

	ram, ram_err := mem.VirtualMemory()
	if ram_err != nil {
		fmt.Println("Error: ", ram_err)
	}

	return DynamicPCData{CPUUsage: cpu_usage, RAM: *ram}

}
