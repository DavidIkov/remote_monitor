package monitor

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUData struct {
	Name        string
	CoresAmount int
	Arch        string
}

type PCData struct {
	CPU      CPUData
	Hostname string
	Username string
}

func GetPCData() PCData {
	hostname, hostname_err := os.Hostname()
	if hostname_err != nil {
		fmt.Println("Error:", hostname_err)
	}
	username, username_err := user.Current()
	if username_err != nil {
		fmt.Println("Error:", username_err)
	}
	cpuInfo, cpu_err := cpu.Info()
	if cpu_err != nil {
		fmt.Println("Error:", cpu_err)
	}

	return PCData{CPU: CPUData{Name: cpuInfo[0].ModelName, CoresAmount: runtime.NumCPU(), Arch: runtime.GOARCH}, Hostname: hostname, Username: username.Username}
}
