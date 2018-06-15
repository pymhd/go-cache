package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

func runHTTPServer() {
	rolton := mux.NewRouter()
	
	rolton.HandleFunc("/map/{key}", DictOpHandler).Methods("GET", "DELETE", "POST")
	rolton.HandleFunc("/list/{key}", ListOpHandler).Methods("PUT", "DELETE", "GET")
	rolton.HandleFunc("/stats", ShowStatsHandler).Methods("GET")
	
	//rolton.HandleFunc("/", MwManager(novaHandler, Logging())).Methods("GET", "POST")
	//rolton.HandleFunc("/nova", MwManager(novaHandler, Logging())).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":9000", rolton))

}

func DictOpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//action := vars["action"]
	k := vars["key"]
	resp := Response{}
	req := Request{}
	enc := json.NewEncoder(w)
	switch r.Method {
	case http.MethodGet:
		v, err := cache.Fetch(k)
		if err != nil {
			resp.Error = true
			resp.Reason = fmt.Sprintf("%s", err)
			enc.Encode(resp) // write data 
			return 		//   and exit
		}
		resp.Value = v
		enc.Encode(resp) // write data to output and exit
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.Error = true
			resp.Reason = fmt.Sprintf("%s", err)
			enc.Encode(resp) // write data 
			return		 //  and exit
		}
		cache.Insert(k, req.Value, req.TTL)
		enc.Encode(resp) //write data and exit
	case http.MethodDelete:
		enc.Encode(resp) // FIXME
	}
}


func ListOpHandler(w http.ResponseWriter, r *http.Request) {

}


func ShowStatsHandler(w http.ResponseWriter, r *http.Request) {
}