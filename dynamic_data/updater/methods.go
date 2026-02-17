package updater

import (
	"remote_monitor/monitor"
	"time"
)

func CreateUpdater(info UpdaterInfo) *Updater {
	updater := Updater{}

	updater.data = make([]monitor.DynamicPCData, info.MeasureAmount)

	go func() {
		for true {
			data := monitor.GetDynamicPCData()
			updater.mutex.Lock()
			updater.data = updater.data[1:]
			updater.data = append(updater.data, data)
			updater.mutex.Unlock()
			time.Sleep(time.Duration(info.MeasureTime) * time.Millisecond / time.Duration(info.MeasureAmount))
		}
	}()

	return &updater
}
func (updater *Updater) GetData() monitor.DynamicPCData {
	updater.mutex.Lock()
	data := monitor.GetAverageDynamicPCData(updater.data)
	updater.mutex.Unlock()
	return data
}
