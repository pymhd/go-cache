package main

import (
	"flag"
	logging "github.com/pymhd/go-logging"
	"os"
)

var (
	config *Config
	log    *logging.Logger
	cache  *Cache
)

func main() {
	flag.Parse()
	config = ParseConfig("./config.yaml")
	logFileWriter, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log = logging.New(logFileWriter)
	cache = NewCache(config.SyncFile, config.MaxItems, config.SyncTimeSec)

	runHTTPServer(config.HTTP.Ip, config.HTTP.Port, config.HTTP.SSL, config.HTTP.Crt, config.HTTP.Key)
}
