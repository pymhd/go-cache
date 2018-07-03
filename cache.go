package main

import (
	"sync"
	"time"
        "os"
        "os/signal"
        "syscall"
)

type Backend interface {
	//map[string]interface{} methods
	Put(k string, v interface{}, ttl int) error
	Get(s string) (interface{}, error)
	Del(s ...string)
	//map[string][]interface{} methods
	Add(s string, v interface{}) error
	Pop(s string, v interface{})
	IsIn(s string, v interface{}) bool
	//Get out of date keys
	GetExpiredKeys() ([]string, error)
	//Flush in memory obj to disk
	Flush()
	BlockingFlush()
	Len() int
	Size() int64
	// interrupt expired keys search
	//Break()
}

type Cache struct {
	Backend
	sync.Mutex
        
}

func (c *Cache) Insert(k string, v interface{}, ttl int) error {
	c.Lock()
	defer c.Unlock()
	//c.Break() // need to break expired keys search
	return c.Put(k, v, ttl)
}

func (c *Cache) Fetch(k string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	return c.Get(k)
}

func (c *Cache) Remove(k ...string) {
	c.Lock()
	defer c.Unlock()
	//c.Break() // need to break expired keys search
	c.Del(k...)
}

func (c *Cache) Append(k string, v interface{}) error {
	c.Lock()
	defer c.Unlock()
	//c.Break() // need to break expired keys search
	return c.Add(k, v)
}

func (c *Cache) Reduce(k string, v interface{}) {
	c.Lock()
	defer c.Unlock()
	//c.Break() // need to break expired keys search
	c.Pop(k, v)
}

func (c *Cache) IsMember(k string, v interface{}) bool {
	c.Lock()
	defer c.Unlock()
	return c.IsIn(k, v)
}

func (c *Cache) Save() {
	c.Lock()
	defer c.Unlock()
	//c.Break() // need to break expired keys search
	c.Flush()
}

func (c *Cache) BlockingSave() {
    c.Lock()
    defer c.Unlock()
    c.BlockingFlush()
}

func (c *Cache) CleanUp() {
	c.Lock()
	defer c.Unlock()
	
	//start := time.Now()
	
	expired, err  := c.GetExpiredKeys()
	if err != nil {
		log.Debug(err)
		//someone broke up process
		return 
	}
	if len(expired) > 0 {
		c.Del(expired...) 
	}
	//log.Info("cleanup func took: ", time.Since(start))
}

func (c *Cache) Stats() *HealthResponse {
	c.Lock() 
	defer c.Unlock() 
	ret := HealthResponse{}
	ret.Max =  100
	ret.Items = c.Len()
	ret.FileSize = c.Size()
	return &ret
}


func NewCache(filename string, maxElem, syncTime int) *Cache {
	backend := NewBackend(maxElem, filename)
	cache := Cache{}
	cache.Backend = backend
	flushTicker := time.NewTicker(time.Duration(syncTime) * time.Second)
	cleanTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for _ = range flushTicker.C {
			cache.Save()  //
		}
	}()
	go func() {
		for _ = range cleanTicker.C {
			cache.CleanUp()
		}
	}()
        go func() {
                sigs := make(chan os.Signal, 1)
                signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
                sig := <- sigs
                log.Infof("Signal: '%s' caught, saving cache to disk and exit with status 0\n", sig)
                cache.BlockingSave()
                os.Exit(0)
        }()
	return &cache
}



/*
func (c *Cache) handleData(d Request) {
	r := Reply{}
	switch d.Action {
	case "insert":
		if err := c.Insert(d.Key, d.Value, d.TTL); err != nil {
			r.Error = true
			r.Reason = fmt.Sprintf("%s", err)
		}
		d.Pipe <- r
	case "fetch":
		v, err := c.Fetch(d.Key)
		if err != nil {
			r.Error = true
			r.Reason = fmt.Sprintf("%s", err)
		}
		r.Result = v
		d.Pipe <- r
	case "append":
		err := c.Add(d.Key, d.Value)
		if err != nil {
			r.Error = true
			r.Reason = fmt.Sprintf("%s", err)
		}
		d.Pipe <- r
	}
}

func (c *Cache) Start(port, syncTime int) {
	ch := make(chan Request, 50)
	go RunSocketServer(port, ch)
	for i := 0; i < 5; i++ {
		go func() {
			for data := range ch {
				c.handleData(data)
			}
		}()
	}
	ticker := time.NewTicker(time.Duration(syncTime) * time.Second)
	go func() {
		for _ = range ticker.C {
			c.Lock()
			c.Flush()
			c.Unlock()
		}
	}()
	cleanUpTicker := time.NewTicker(time.Duration(60) * time.Second)
	go func() {
		for _ = range cleanUpTicker.C {
		        expired := c.GetExpiredKeys()
		        c.Remove(expired...)
		}
	}()

}

func NewServer(filename string, port, maxElem, syncTime int) *Cache {
	backend := NewBackend(maxElem, filename)
	cache := Cache{}
	cache.Backend = backend
	cache.Start(port, syncTime)
	return &cache
}

//func NewClient(ip string, port int) *CacheClient {
//}

*/
