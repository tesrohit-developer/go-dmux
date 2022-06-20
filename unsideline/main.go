package unsideline

import (
	"github.com/tesrohit-developer/go-dmux/plugins"
	"log"
	"net/http"
)

func getScanPlugin() interface{} {
	s := plugins.NewManager("scan_plugin", "scan-*", "./plugins/built", &plugins.ScanImplPlugin{})
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

}

func unsideline(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/scan", scan)
	http.HandleFunc("/unsideline", unsideline)
	log.Fatal(http.ListenAndServe(":9951", nil))
}
