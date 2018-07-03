package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func runHTTPServer(host string, port int, ssl bool, crt, key string) {
	rolton := mux.NewRouter()

	rolton.HandleFunc("/map/{key}", MwManager(DictOpHandler, CheckAccess(config.HTTP.AllowNets))).Methods("GET", "DELETE", "POST")
	rolton.HandleFunc("/list/{key}", MwManager(ListOpHandler, CheckAccess(config.HTTP.AllowNets))).Methods("PUT", "DELETE", "GET")
	rolton.HandleFunc("/healthz", MwManager(ShowStatsHandler, CheckAccess(config.HTTP.AllowNets))).Methods("GET")
	if !ssl {
		log.Fatal(http.ListenAndServe(host+":"+strconv.Itoa(port), rolton))
	} else {
		log.Fatal(http.ListenAndServeTLS(host+":"+strconv.Itoa(port), crt, key, rolton))
	}
}

func DictOpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
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
			return           //   and exit
		}
		resp.Value = v
		enc.Encode(resp) // write data to output and exit
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.Error = true
			resp.Reason = fmt.Sprintf("%s", err)
			enc.Encode(resp) // write data
			return           //  and exit
		}
		cache.Insert(k, req.Value, req.TTL)
		enc.Encode(resp) //write data and exit
	case http.MethodDelete:
		cache.Remove(k)
		enc.Encode(resp)
	}
}

func ListOpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["key"]
	resp := Response{}
	req := Request{}
	enc := json.NewEncoder(w)
	switch r.Method {
	case http.MethodGet:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.Error = true
			resp.Reason = fmt.Sprintf("%s", err)
			enc.Encode(resp) // write data
			return
		}
		v := cache.IsMember(k, req.Value) //true or false
		resp.Value = v
		enc.Encode(resp) // write data to output and exit
	case http.MethodPut:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.Error = true
			resp.Reason = fmt.Sprintf("%s", err)
			enc.Encode(resp) // write data
			return
		}
		cache.Append(k, req.Value)
		enc.Encode(resp)
	case http.MethodDelete:
		cache.Reduce(k, req.Value)
		enc.Encode(resp)
	}
}

func ShowStatsHandler(w http.ResponseWriter, r *http.Request) {
	resp := cache.Stats()
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
