package main

import (
	//"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"strings"
	"strconv"
)

//type IpAddr string
//type Net string

type Config struct {
	Host          string `yaml:"host"`
	MaxItems      int    `yaml:"MaxItems"`
	DefaultTTL    string `yaml:"DefaultTTL"`
	DefaultTTLSec int
	SyncFile      string `yaml:"SyncFile"`
	SyncTime      string `yaml:"SyncTime"`
	SyncTimeSec   int
	LogFile       string `yaml:"LogFile"`
	HTTP          struct {
		Ip        string   `yaml:"Ip"`
		Port      int      `yaml:"Port"`
		SSL       bool     `yaml:"SSL"`
		Crt       string   `yaml:"Crt"`
		Key       string   `yaml:"Key"`
		Allow     []string `yaml:"Allow"`
		Deny      []string `yaml:"Deny"`
		AllowNets []*net.IPNet
		DenyNets  []*net.IPNet
	} `yaml:"HTTP"`
}

func ParseConfig(filename string) *Config {
	cfg := Config{}
	fb, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(fb, &cfg); err != nil {
		panic(err)
	}
	for _, n := range cfg.HTTP.Allow {
		if !strings.Contains(n, "/") {
			n += "/32"
		}
		_, network, err := net.ParseCIDR(n)
		if err != nil {
			panic(err)
		}
		cfg.HTTP.AllowNets = append(cfg.HTTP.AllowNets, network)
	}
	cfg.SyncTimeSec = parsePeriod(cfg.SyncTime)
	cfg.DefaultTTLSec = parsePeriod(cfg.DefaultTTL)
	return &cfg
}


func parsePeriod(s string) int {
	unit := string(s[len(s) - 1])
	number, err := strconv.Atoi(s[:len(s) - 1])
	if err != nil {
		panic(err)
	}
	var multiplier int
	switch unit {
	case "s":
	    multiplier = 1
	case "m":
	    multiplier = 60
	case "h":
	    multiplier = 3600
	case "d":
	    multiplier = 24 * 3600
	case "w":
	    multiplier = 7 * 24 * 3600
	default:
	    panic("Unknown type of time")
	}
        return multiplier * number
}

//func main() {
//	fmt.Printf("%+v", ParseConfig("./config.yaml"))
//}
