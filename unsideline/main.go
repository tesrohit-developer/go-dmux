package main

import (
	"encoding/json"
	"github.com/flipkart-incubator/go-dmux/plugins"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var scanPlugin interface{}
var unsidelinePlugin interface{}

func getScanPlugin() interface{} {
	s := plugins.NewManager("scan_plugin", "scan-*", "", &plugins.ScanImplPlugin{})
	//defer s.Dispose()
	err := s.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Launch()
	p, err := s.GetInterface("sideline-em")
	if err != nil {
		log.Fatal(err.Error())
	}
	return p
}

func getUnsidelinePlugin() interface{} {
	s := plugins.NewManager("unsideline_plugin", "unsideline-*", "", &plugins.UnsidelineImplPlugin{})
	//defer s.Dispose()
	err := s.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Launch()
	p, err := s.GetInterface("em")
	if err != nil {
		log.Fatal(err.Error())
	}
	return p
}

func scan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var request = plugins.ScanWithStartRowEndRowRequest{
		StartKey: vars["startRow"],
		EndKey:   vars["endRow"],
	}
	rows, err := scanPlugin.(plugins.ScanImpl).ScanWithStartRowEndRow(request)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rows)
}

func unsideline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var request = plugins.UnsidelineByKeyRequest{
		Key: vars["key"],
	}
	rows, err := unsidelinePlugin.(plugins.UnsidelineImpl).UnsidelineByKey(request)
	if err != nil {
		log.Printf("Received error while executing plugin")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		w.Header().Set("Content-Type", "application/json")
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rows)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	log.Println("Hi starting the API")
	//scanPlugin = getScanPlugin()
	//unsidelinePlugin = getUnsidelinePlugin()
	//r := mux.NewRouter()
	//r.HandleFunc("/scan/{startRow}/{endRow}", scan)
	//r.HandleFunc("/unsideline/{key}", unsideline)
	//r.HandleFunc("/healthCheck", healthCheck)
	//log.Fatal(http.ListenAndServe(":9951", r))
}
