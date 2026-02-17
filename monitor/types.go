package monitor

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
