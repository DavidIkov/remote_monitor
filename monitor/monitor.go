package monitor

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"sync"
	"time"

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

type ramInfo struct {
	Total uint64
	Used  float32
}

type DynamicPCData struct {
	CPUUsage []float32
	RAM      ramInfo
}

func GetAverageDynamicPCData(dataRange []DynamicPCData) DynamicPCData {
	if len(dataRange) == 0 {
		return DynamicPCData{}
	}
	avgData := dataRange[0]
	for i := 1; i < len(dataRange); i++ {
		for j, usage := range dataRange[i].CPUUsage {
			avgData.CPUUsage[j] = (avgData.CPUUsage[j]*float32(i) + usage) / float32(i+1)
		}
		avgData.RAM.Total = uint64((float32(avgData.RAM.Total)*float32(i) + float32(dataRange[i].RAM.Total)) / float32(i+1))
		avgData.RAM.Used = (avgData.RAM.Used*float32(i) + dataRange[i].RAM.Used) / float32(i+1)
	}
	return avgData
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
	cpu_usage64, cpu_usage_err := cpu.Percent(0, true)
	if cpu_usage_err != nil {
		fmt.Println("Error: ", cpu_usage_err)
	}
	cpu_usage := make([]float32, len(cpu_usage64))
	for i, v := range cpu_usage64 {
		cpu_usage[i] = float32(v)
	}

	ram, ram_err := mem.VirtualMemory()
	if ram_err != nil {
		fmt.Println("Error: ", ram_err)
	}

	return DynamicPCData{CPUUsage: cpu_usage, RAM: ramInfo{Total: ram.Total, Used: float32(ram.UsedPercent)}}
}

type DynamicDataUpdater struct {
	mutex sync.Mutex
	data  []DynamicPCData
	// In milliseconds.
	measureTime   uint32
	measureAmount uint32
}

func (updater *DynamicDataUpdater) Start(measureTime, measureAmount uint32) {
	updater.measureTime = measureTime
	updater.measureAmount = measureAmount
	updater.data = make([]DynamicPCData, measureAmount)

	go func() {
		for true {
			data := GetDynamicPCData()
			updater.mutex.Lock()
			updater.data = updater.data[1:]
			updater.data = append(updater.data, data)
			updater.mutex.Unlock()
			time.Sleep(time.Duration(measureTime) * time.Millisecond / time.Duration(measureAmount))
		}
	}()
}
func (updater *DynamicDataUpdater) GetData() DynamicPCData {
	updater.mutex.Lock()
	data := GetAverageDynamicPCData(updater.data)
	updater.mutex.Unlock()
	return data
}
