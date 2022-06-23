package unsideline

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tesrohit-developer/go-dmux/plugins"
	"log"
	"net/http"
)

var scanPlugin interface{}

func getScanPlugin() interface{} {
	s := plugins.NewManager("scan_plugin", "scan-*", "", &plugins.ScanImplPlugin{})
	//defer s.Dispose()
	err := s.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Launch()
	p, err := s.GetInterface("em-scan")
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
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(rows)
}

func unsideline(w http.ResponseWriter, r *http.Request) {

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	scanPlugin = getScanPlugin()
	http.HandleFunc("/scan/{startRow}/{endRow}", scan)
	http.HandleFunc("/unsideline", unsideline)
	http.HandleFunc("/healthCheck", healthCheck)
	log.Fatal(http.ListenAndServe(":9951", nil))
}
