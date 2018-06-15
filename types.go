package main

type Request struct {
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

type Response struct {
	Value  interface{} `json:"value,omitempty"`
	Error  bool        `json:"error,omitempty"`
	Reason string      `json:"reason,omitempty"`
}

/*

type Conductor chan interface{}


type Request struct {
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
	Action string      `json:"method"` // put, get, in
	TTL    int  	   `json:"ttl"`
	Pipe   Conductor
}

type Reply struct {
	Error  bool
	Result interface{}
	Reason string
}

*/
