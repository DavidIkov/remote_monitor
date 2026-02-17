package graph_creator

import (
	"remote_monitor/dynamic_data/updater"
	"remote_monitor/monitor"
	"sync"
)

type CreatorInfo struct {
	MeasureTime   uint32
	MeasureAmount uint32
}

type Creator struct {
	mutex   sync.Mutex
	data    []monitor.DynamicPCData
	updater *updater.Updater
}
