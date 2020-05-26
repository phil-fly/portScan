package work

import (
	"sync"
	"portScan/utils/ping"
	"fmt"
)

type Workdist struct {
	Host	string
	Port    string
}

const (
	taskload		    = 60000
)
var wg sync.WaitGroup

func Task(ips []string,PortMap *PortSet,tasknum int) {
	tasks := make(chan Workdist,taskload)
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
				}
				tasks <- task
			}
		}
	}
	close(tasks)
	wg.Wait()
	return
}

func worker(tasks chan Workdist){
	defer wg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}

		if IsOpenTCP(task.Host,task.Port) {
			fmt.Printf("[TCP]\t%s:%s\topen\n",task.Host,task.Port)
		}

	}
}