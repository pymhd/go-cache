package main

import (
	"fmt"
	"net"
	"io/ioutil"
	"strings"
	//"encoding/json"
	"gopkg.in/yaml.v2"
)

type IpAddr string
type Net string

type Config struct {
	Host       string `yaml:"host"`
	MaxItems   int    `yaml"MaxItems"`
	DefaultTTL string `yaml:"DefaultTTL"`
	SyncFile   string `yaml:"SyncFile"`
	LogFile    string `yaml:"LogFile"`
	HTTP       struct {
	                    Ip     string `yaml:"Ip"`
		            Port   int    `yaml:"Port"`
		            SSL    bool   `yaml:"SSL"`
		            Crt    string `yaml:"Crt"`
		            Key    string `yaml:"Key"`
		            Allow  []string `yaml:"Allow"`
		            Deny   []string `yaml:"Deny"`
	} `yaml:"HTTP"`
}

func main() {
	f := "./config.yaml"
	fb, _ := ioutil.ReadFile(f)
	cfg := Config{}
	// get fikle extension
	ext := strings.Split(f, ".")[len(strings.Split(f, "."))-1]
	switch ext {
	case "yaml":
		fmt.Println("ямл")
		if err := yaml.Unmarshal(fb, &cfg); err != nil {
			fmt.Println(err)
		}
	case "json":
		fmt.Println("жесон")
	}
	//fmt.Printf("%+v\n", cfg)
	nets := make([]*net.IPNet, 0)
	for _, n := range cfg.HTTP.Allow {
	    if strings.Contains(n, "/") {
	       _, network, err := net.ParseCIDR(n)
	           if err != nil {
	           	panic(err)
	           }
	           nets = append(nets, network)
	    }
	}
	fmt.Println(nets)
	
}
