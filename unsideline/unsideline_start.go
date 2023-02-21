package unsideline

import (
	"encoding/json"
	unsideline_models "github.com/flipkart-incubator/go-dmux/unsideline-models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var scanImpl Scan
var unsidelineImpl Unsideline

func scan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var request = unsideline_models.ScanWithStartRowEndRowRequest{
		StartKey: vars["startRow"],
		EndKey:   vars["endRow"],
	}
	rows, err := scanImpl.ScanWithStartRowEndRow(request)
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
	var request = unsideline_models.UnsidelineByKeyRequest{
		Key: vars["key"],
	}
	rows, err := unsidelineImpl.UnsidelineByKey(request)
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

func UnsidelineStart(scanImplArg Scan, unsidelineImplArg Unsideline, configPath string) {
	log.Println("Hi starting the API")

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	var conf unsideline_models.UnsidelineContainerConfig
	confSerdeErr := json.Unmarshal(raw, &conf)
	if confSerdeErr != nil {
		log.Fatal(confSerdeErr.Error())
	}
	scanImpl = scanImplArg
	unsidelineImpl = unsidelineImplArg
	r := mux.NewRouter()
	r.HandleFunc("/scan/{startRow}/{endRow}", scan)
	r.HandleFunc("/unsideline/{key}", unsideline)
	r.HandleFunc("/healthCheck", healthCheck)
	log.Fatal(http.ListenAndServe(":"+strconv.FormatInt(conf.Port, 10), r))
}
