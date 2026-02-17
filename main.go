package main

import (
	"html/template"
	"log"
	"net/http"
	data_graph "remote_monitor/dynamic_data/graph_creator"
	data_updater "remote_monitor/dynamic_data/updater"
	"remote_monitor/monitor"
)

type AllPCData struct {
	StaticData  monitor.StaticPCData
	DynamicData monitor.DynamicPCData
}

var GraphCreator *data_graph.Creator = nil

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("panel/template.html")
	t.Execute(w, AllPCData{monitor.GetStaticPCData(), GraphCreator.GetLastData()})
}

func main() {
	GraphCreator = data_graph.CreateCreator(20000, data_updater.UpdaterInfo{MeasureTime: 1000, MeasureAmount: 10})
	fs := http.FileServer(http.Dir("./graphs"))
	http.Handle("/graphs/", http.StripPrefix("/graphs/", fs))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
