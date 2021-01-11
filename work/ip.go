package work

import (
	"portScan/utils/ping"
	"sync"
)

type IpScan struct {
	ipList	[]string	`json:"ipList"`
	aliveIpList	[]string	`json:"aliveIpList"`
	tasknum	int	`json:"tasknum"`
	wg sync.WaitGroup
}

func NewIpScan() *IpScan {
	return &IpScan{}
}

func (t *IpScan)SetIpList(ips []string){
	t.ipList = ips
}

func (t *IpScan)SetTasknum(tasknum int){
	t.tasknum = tasknum
}

func (t *IpScan)RunScan() error {

	tasks := make(chan IpScanTaskInfo,taskload)
	t.wg.Add(t.tasknum)
	for gr:=1;gr<= t.tasknum;gr++ {
		go t.worker(tasks)
	}

	//创建chan生产者
	for _,host := range t.ipList {
		task := IpScanTaskInfo{
			Host:host,
		}
		tasks <- task
	}
	close(tasks)

	t.wg.Wait()
	return nil
}

type IpScanTaskInfo struct {
	Host	string
}

func (t *IpScan)worker(tasks chan IpScanTaskInfo){
	defer t.wg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}

		if ping.Ping(task.Host) {
			t.aliveIpList = append(t.aliveIpList,task.Host)
		}
	}
}