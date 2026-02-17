package updater

import (
	"remote_monitor/monitor"
	"sync"
)

type UpdaterInfo struct {
	// In milliseconds.
	MeasureTime   uint32
	MeasureAmount uint32
}

type Updater struct {
	mutex sync.Mutex
	data  []monitor.DynamicPCData
}
