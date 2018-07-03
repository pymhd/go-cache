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

type HealthResponse struct {
	Items    int   `json:"items"`
	FileSize int64 `json:"size_kb"`
	Max      int   `json:"max"`
}


