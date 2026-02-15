package main

import (
	"html/template"
	"log"
	"net/http"
	"remote_monitor/monitor"
)

type AllPCData struct {
	StaticData  monitor.StaticPCData
	DynamicData monitor.DynamicPCData
}

var DynamicData monitor.DynamicDataUpdater

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("panel/template.html")
	t.Execute(w, AllPCData{monitor.GetStaticPCData(), DynamicData.GetData()})
}

func main() {
	DynamicData.Start(1000, 10)
	fs := http.FileServer(http.Dir("./graphs"))
	http.Handle("/graphs/", http.StripPrefix("/graphs/", fs))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
