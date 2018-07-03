package main

import (
	"fmt"
	"encoding/gob"
	"errors"
	"os"
	"time"
	"sync"
)

// main item
type Value interface{}
//type Value struct {
//	Data interface{}
//	Born time.Time
//	TTL  time.Duration
//}

// Key -> value data
type UnderlayDict map[string]Value

// key -> list of values (slice) data
type UnderlayList map[string][]interface{}

// key -> time when item will be expired data
type Curator map[string]time.Time

// Cache Underlay data type
type OrderedDict struct {
	Curator
	UnderlayDict
	UnderlayList
	Maxelem  int
	Filename string
	order    []string
	breakCh  chan bool
	wg	 sync.WaitGroup
}

func (o *OrderedDict) Put(k string, v interface{}, ttl int) error {
	log.Debugf("Backend PUT op: key -> (%s). Value type: %T, ttl = %d\n", k, v, ttl)
	if ttl == 0 {
		ttl =  config.DefaultTTLSec
	}
	//V := Value{Data: v, Born: time.Now(), TTL: time.Duration(ttl) * time.Second}
	// set up expiration time
	willDie := time.Now().Add(time.Duration(ttl) * time.Second)
	//o.wg.Add(1)
	o.Curator[k] = willDie
	//o.wg.Done()
	
	previousVal := o.UnderlayDict[k]
	o.UnderlayDict[k] = v

	// check if value with the same key existed before, if true, no need to check cache size
	if previousVal != nil {
		log.Debugf("Value with key (%s) alredy exist, just rewrite it\n", k)
		return nil
	}

	size := len(o.order)
	o.order = append(o.order, k)

	if size >= o.Maxelem {
		log.Debugf("Cache size is more then allowed (%d), deleteing key...\n", o.Maxelem)
		keyToDelete := o.order[0]
		log.Debugf("%s\n", keyToDelete)
		o.order = o.order[1:]
		delete(o.UnderlayDict, keyToDelete)
		delete(o.Curator, keyToDelete)
		log.Debug("Order rearranged, extra item deleted")
	}
	return nil
}

func (o *OrderedDict) Get(k string) (interface{}, error) {
	log.Debugf("Backend GET op: key -> (%s)\n", k)
	v := o.UnderlayDict[k]
	//emptyVal := Value{}
	if v != nil {
		if diff := time.Now().Sub(o.Curator[k]); diff >= 0 {
			return nil, errors.New("Expired Cache")
		}
		return v, nil
	} else {
		log.Debug("Empty Cache, error is returned")
		return nil, errors.New("Empty Cache")
	}
}

func (o *OrderedDict) Del(k ...string) {
	log.Debugf("Backend DEL op: key -> (%s)\n", k)
	for _, key := range k {
		delete(o.UnderlayDict, key)
		delete(o.Curator, key)
	}
}

func (o *OrderedDict) Add(k string, v interface{}) error {
	log.Debugf("Backend ADD op: key -> (%s)\n", k)
	o.UnderlayList[k] = append(o.UnderlayList[k], v)
	return nil
	//return errors.New("Wrong value type")
}

func (o *OrderedDict) Pop(k string, v interface{}) {
	log.Debugf("Backend POP op: key -> (%s)\n", k)
	index := getIndex(v, o.UnderlayList[k])
	if index >= 0 {
		o.UnderlayList[k] = append(o.UnderlayList[k][:index], o.UnderlayList[k][index+1:]...)
	}
}

func (o *OrderedDict) IsIn(k string, v interface{}) bool {
	index := getIndex(v, o.UnderlayList[k])
	return index >= 0
}

func getIndex(v interface{}, s []interface{}) int {
	index := -1
	for i, item := range s {
		if item == v {
			index = i
		}
	}
	return index
}

func (o *OrderedDict) Flush() {
        st := time.Now()
        defer func() {
                log.Info("Object deep copy took: ", time.Since(st))
        }()
        var no OrderedDict
	c := make(Curator, len(o.Curator))
        d := make(UnderlayDict, len(o.UnderlayDict))
        s := make(UnderlayList, len(o.UnderlayList))
        order := make([]string, len(o.order))
        no = OrderedDict{c, d, s, o.Maxelem, o.Filename, order, make(chan bool), sync.WaitGroup{}}
        for k, v := range o.Curator {
        	no.Curator[k] = v
        }
        for k, _ := range o.UnderlayList {
        	copy(no.UnderlayList[k], o.UnderlayList[k])
        }
        for k, _ := range o.UnderlayDict {
            no.UnderlayDict[k] = o.UnderlayDict[k]
        }
        copy(no.order, o.order) //fine
	go writeGob(no.Filename, no)
}

func (o *OrderedDict) BlockingFlush() {
        writeGob(o.Filename, o)   
}

func (o *OrderedDict) Len() int {
	dictElements  := len(o.UnderlayDict)
	var listElements int
	for _, l := range o.UnderlayList {
	    listElements += len(l)
	}
	return dictElements + listElements
}

func (o *OrderedDict) Size() int64 {
        fi, err := os.Stat(o.Filename)
        if err != nil {
            return 0
        }
        return fi.Size() / 1024
}

/*
func (o *OrderedDict) Break() {
	log.Debug("Checking procedure to break up")
	done := make(chan bool)
	go func() {
		o.wg.Wait()
		done <- true
	}()
	select {
	case <- done:
		//nothing is going on right now
		log.Debug("No active process detected, nothing to break...")
	case <- time.After(1 * time.Millisecond):
		log.Debug("Seems that there is smth going on, sendig break cmd to break channel")
		o.breakCh <- true
		time.Sleep(30 * time.Millisecond)
	}
}


func (o *OrderedDict) GetExpiredKeys() ([]string, error) {
	o.wg.Add(1)
	defer o.wg.Done()
	ret := make([]string, 0)
	resultC := make(chan []string, 0)
	sem := 0
	go func(){
		for k, v := range o.UnderlayDict {
			if sem == 1 {
				log.Debug("Found semaphore, exiting gouroutine")
				return // stop iteration
			}
			if diff := time.Now().Sub(v.Born.Add(v.TTL)); diff >= 0 {
				ret = append(ret, k)
			}
		}
		resultC <- ret
	}()
	select {
	case <- o.breakCh:
		sem = 1 
		return []string{}, errors.New("Operation was broke up")
	case ret = <- resultC:
	    	return ret, nil
	}
}
*/

func (o *OrderedDict) GetExpiredKeys() ([]string, error) { 
	ret := make([]string, 0)
	now := time.Now()
	for k, v := range o.Curator {
		if diff := now.Sub(v); diff >= 0 {
			ret = append(ret, k)
		}
	}
	log.Info("Number of iyems in cache: ", len(o.Curator))
	return ret, nil
}


func NewBackend(maxelem int, filename string) *OrderedDict {
	var ret *OrderedDict
	if err := readGob(filename, &ret); err == nil {
		fmt.Println("Successfully init form file")
		ret.Maxelem = maxelem //change in case of value changed
		return ret
	}
	c := make(Curator, 0)
	d := make(UnderlayDict, 0)
	s := make(UnderlayList, 0)
	//wg := sync.WaitGroup{}
	ret = &OrderedDict{c, d, s, maxelem, filename, make([]string, 0), make(chan bool), sync.WaitGroup{}}
	return ret
}

func writeGob(filePath string, object interface{}) error {
	st := time.Now()
	defer func() {
		log.Info("Disk sync took: : ", time.Since(st))
	}()
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}


func readGob(filePath string, object interface{}) error {
       file, err := os.Open(filePath)
       if err == nil {
              decoder := gob.NewDecoder(file)
              err = decoder.Decode(object)
       }
       file.Close()
       return err
}


