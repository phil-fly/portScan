package work

import (
	"sync"
	"portScan/utils/ping"
)

type Workdist struct {
	Host	string
	Port    string
	Results *ResultsSet
}

const (
	taskload		    = 60000
)
var wg sync.WaitGroup

func Task(ips []string,PortMap *PortSet,tasknum int) *ResultsSet {
	tasks := make(chan Workdist,taskload)

	Results := InitResults()
	wg.Add(tasknum)
	//创建chan消费者worker
	for gr:=1;gr<=tasknum;gr++ {
		go worker(tasks)
	}

	//创建chan生产者
	for _,host := range ips {
		if ping.Ping(host) {
			for Port,_ := range PortMap.Port {
				task := Workdist{
					Host:host,
					Port:Port,
					Results:Results,
				}
				tasks <- task
			}
		}
	}
	close(tasks)
	wg.Wait()
	return Results
}

func worker(tasks chan Workdist){
	defer wg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}
		if IsOpenTCP(task.Host,task.Port) {
			task.Results.Set(task.Host,task.Port)
			//log.Print("Host -> %s  Port -> %s  opend",task.Host,task.Port)
		}
	}
}


type ResultsSet struct {
	Results  map[string][]string
	ResultsMutex  *sync.RWMutex
}

func ResultsSetNew() *ResultsSet {
	return &ResultsSet{
		Results:make(map[string][]string),
		ResultsMutex:new(sync.RWMutex),
	}
}

func (this *ResultsSet)Set(host string,value interface{}) {
	this.ResultsMutex.Lock()
	defer this.ResultsMutex.Unlock()
	switch value.(type) {
	case string:
		this.Results[host] = append(this.Results[host], value.(string))
	case []string:
		this.Results[host] = append(this.Results[host], value.([]string)...)
	}
}

func (this *ResultsSet)Get(key string) []string{
	this.ResultsMutex.RLock()
	defer this.ResultsMutex.RUnlock()
	return  this.Results[key]
}

func InitResults() *ResultsSet {
	return ResultsSetNew()
}

var  ResultsMap *ResultsSet