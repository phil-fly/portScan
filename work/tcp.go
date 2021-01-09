package work

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"portScan/utils/ping"
	"sync"
	"time"
)

type TcpScan struct {
	ipList	[]string	`json:"ipList"`
	portMap *PortSet	`json:"portMap"`
	tasknum	int	`json:"tasknum"`
	wg sync.WaitGroup
	resultsOutput 	string 	`json:"resultsOutput"`
	resultChan	chan string
}

type TcpScanTaskInfo struct {
	Host	string
	Port    string

}

func (t *TcpScan)SetIpList(ips []string){
	t.ipList = ips
}

func (t *TcpScan)SetPortMap(portMap *PortSet){
	t.portMap = portMap
}

func (t *TcpScan)SetTasknum(tasknum int){
	t.tasknum = tasknum
}

func (t *TcpScan)SetResultsOutput(resultsOutput string) {
	t.resultsOutput = resultsOutput
	file,err :=os.Open(resultsOutput)
	defer file.Close()
	if err !=nil && os.IsNotExist(err) {
		file ,_= os.Create(resultsOutput)
	}
}

func (t *TcpScan)Validate() error {
	switch {
	case len(t.ipList)< 1 :
		return errors.New("TcpScan ipList is nil.")
	case t.tasknum ==  0:
		t.tasknum = 1000
		return nil
	case t.portMap == nil :
		return errors.New("TcpScan portMap is nil.")
	}
	return nil
}

func (t *TcpScan)RunScan() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	//结果写入文件
	t.resultChan = make(chan string)
	defer close(t.resultChan)
	t.writeResultToFile()


	tasks := make(chan TcpScanTaskInfo,taskload)

	for gr:=1;gr<= t.tasknum;gr++ {
		t.wg.Add(1)
		go t.worker(tasks)
	}

	//创建chan生产者
	for _,host := range t.ipList {
		if ping.Ping(host) {
			for Port,_ := range PortMap.Port {
				task := TcpScanTaskInfo{
					Host:host,
					Port:Port,
				}
				tasks <- task
			}
		}
	}
	close(tasks)
	t.wg.Wait()
	return nil
}

func (t *TcpScan)worker(tasks chan TcpScanTaskInfo){
	defer t.wg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}

		if t.IsOpenTCP(task.Host,task.Port) {
			result := fmt.Sprintf("[TCP]\t%s:%s\topen\n",task.Host,task.Port)
			t.resultChan <- result
		}
	}
}

// 写文件
func (t *TcpScan)writeResultToFile() {
	var f *os.File
	var err error
	f, err = os.OpenFile(t.resultsOutput, os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for {
		select{
		case res,ok:=<- t.resultChan:
			if ok{
				_, _ = w.WriteString(res)
			}
		default:
			w.Flush()
			return
		}
	}
}


func (t *TcpScan)IsOpenTCP(IpAddr,Port string) bool {
	conn, err := net.DialTimeout("tcp", IpAddr+":"+Port, time.Second*1)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}