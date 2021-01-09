package work

import (
	"strings"
	"sync"
	"strconv"
	"log"
	"os"
)

//  1,2,3,100-1000

type PortSet struct {
	Port  map[string]int
	PortMutex  *sync.RWMutex
}

func New() *PortSet {
	return &PortSet{
		Port:make(map[string]int),
		PortMutex:new(sync.RWMutex),
	}
}

func (this *PortSet)Set(key interface{},value interface{}) {
	this.PortMutex.Lock()
	defer this.PortMutex.Unlock()
	switch key.(type) {
	case int:
		keyStr := strconv.Itoa(key.(int))
		this.Port[keyStr] = value.(int)
	case string:
		this.Port[key.(string)] = value.(int)
	}
}

func (this *PortSet)PLGet(key string) (int,bool) {
	this.PortMutex.RLock()
	defer this.PortMutex.RUnlock()
	value,ok :=this.Port[key]
	return value,ok
}


func InitPort() *PortSet {
	return New()
}

func String2PortMap(portStr string) *PortSet{
	PortMap := InitPort()
	ports := strings.Split(portStr, ",")

	for _,v:= range ports {
		if strings.Contains(v, "-") {
			portslist := strings.Split(v, "-")
			PortMap.addPort(portslist[0],portslist[1])
		}else{
			PortMap.Set(v,1)
		}
	}
	return PortMap
}
func (this *PortSet)addPort(portStart,portEnd string){
	portStartInt, err := strconv.Atoi(portStart)
	if err != nil || portStartInt < 0 || portStartInt > 65535 {
		log.Fatal("Invalid Port :",portStart)
		os.Exit(-1)
	}

	portEndInt, err := strconv.Atoi(portEnd)
	if err != nil || portEndInt < 0 || portEndInt > 65535 {
		log.Fatal("Invalid Port :",portEnd)
		os.Exit(-1)
	}
	if portStartInt > portEndInt {
		for i := portEndInt;i <=portStartInt ;i++  {
			this.Set(i,1)
		}
	}else{
		for i := portStartInt;i <=portEndInt ;i++  {
			this.Set(i,1)
		}
	}

}
