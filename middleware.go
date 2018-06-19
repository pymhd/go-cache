package main

import (
        "net"
        "net/http"
        "strings"
)


var (

)

type Middleware func(http.HandlerFunc) http.HandlerFunc


func CheckAccess(nets []*net.IPNet) Middleware {
    return func(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            from := strings.Split(r.RemoteAddr, ":")[0]
            fromIP := net.ParseIP(from)
            if fromIP == nil {
                log.Error("Could not resolve IP addr ", r.RemoteAddr)
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
            }
            for _, network := range nets {
                if network.Contains(fromIP) {
                    f(w, r)
                    return
                }
            }
            log.Errorf("Ip addr: %s did not match any allowed network, raising 401 Forbidden response\n", from)
            http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        }
    }
}

func MwManager(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
        for _, m := range middlewares {
                f = m(f)
        }
        return f
}
