package main 


import (
        "os"
        "flag"
        "logging"
)

var (
    log   *logging.Logger
    cache *Cache
)


func main() {
    flag.Parse()    
    log = logging.New(os.Stdout)
    cache = NewCache("/tmp/cache.db", 1000000, 36000)
    runHTTPServer()
}
