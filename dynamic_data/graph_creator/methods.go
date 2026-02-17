package graph_creator

import (
	"remote_monitor/dynamic_data/updater"
	"remote_monitor/monitor"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func CreateCreator(graphTime uint32, updaterInfo updater.UpdaterInfo) *Creator {
	measurementsAmount := uint32(float32(graphTime) / float32(updaterInfo.MeasureTime))
	creator := Creator{data: make([]monitor.DynamicPCData, measurementsAmount), updater: updater.CreateUpdater(updaterInfo)}

	go func() {
		for true {
			data := monitor.GetDynamicPCData()
			creator.mutex.Lock()

			creator.data = creator.data[1:]
			creator.data = append(creator.data, data)

			creator.SaveImages()

			creator.mutex.Unlock()
			time.Sleep(time.Duration(updaterInfo.MeasureTime) * time.Millisecond)
		}
	}()

	return &creator
}

func (creator *Creator) GetLastData() monitor.DynamicPCData {
	creator.mutex.Lock()
	lastData := creator.data[len(creator.data)-1]
	creator.mutex.Unlock()
	return lastData
}

func (creator *Creator) SaveImages() {
	cpuCores := 0
	for _, data := range creator.data {
		cpuCores = max(len(data.CPUUsage), cpuCores)
	}

	cpuUsageGraphData := make([]plotter.XYs, cpuCores)
	for i := range len(cpuUsageGraphData) {
		cpuUsageGraphData[i] = make(plotter.XYs, len(creator.data))
	}

	for dataInd, data := range creator.data {
		for cpuInd, cpuUsage := range data.CPUUsage {
			cpuUsageGraphData[cpuInd][dataInd].X = float64(dataInd)
			cpuUsageGraphData[cpuInd][dataInd].Y = float64(cpuUsage)
		}
	}

	for cpuInd, cpuUsageData := range cpuUsageGraphData {
		p := plot.New()

		p.Title.Text = ""
		p.X.Label.Text = "t"
		p.Y.Label.Text = "%"

		err := plotutil.AddLinePoints(p,
			"", cpuUsageData)
		if err != nil {
			panic(err)
		}

		if err := p.Save(4*vg.Centimeter, 4*vg.Centimeter, "graphs/cpu"+strconv.Itoa(cpuInd)+".png"); err != nil {
			panic(err)
		}
	}

}
