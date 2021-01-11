package work

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

type NbnsScan struct {
	ipList	[]string	`json:"ipList"`
	port string	`json:"port"`
	tasknum	int	`json:"tasknum"`
	timeOut	int `json:"timeOut"`
	resultsOutput 	string 	`json:"resultsOutput"`
	resultChan	chan string
	wg sync.WaitGroup
	nbnsPayload	[]byte
}

type NbnsScanTaskInfo struct {
	Host	string
	Port    string
}

func (t *NbnsScan)SetIpList(ips []string){
	t.ipList = ips
}

func (t *NbnsScan)SetPort(){
	t.port = "137"
}

func (t *NbnsScan)SetTasknum(tasknum int){
	t.tasknum = tasknum
}

func (t *NbnsScan)SetTimeOut(timeOut int){
	t.timeOut = timeOut
}

func (t *NbnsScan)setNbnsPayload(){
	bf := NewBuffer()
	nbns(bf)
	t.nbnsPayload = bf.data
}

func (t *NbnsScan)SetResultsOutput(resultsOutput string) {
	t.resultsOutput = resultsOutput
	file,err :=os.Open(resultsOutput)
	defer file.Close()
	if err !=nil && os.IsNotExist(err) {
		file ,_= os.Create(resultsOutput)
	}
}

func (t *NbnsScan)Validate() error {
	switch {
	case len(t.ipList)< 1 :
		return errors.New("TcpScan ipList is nil.")
	case t.tasknum ==  0:
		t.tasknum = 1000
	case t.port == "" :
		t.SetPort()
	case t.nbnsPayload == nil:
		t.setNbnsPayload()
	}
	return nil
}

func (t *NbnsScan)RunScan() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	//结果写入文件
	t.resultChan = make(chan string,t.tasknum)
	go t.writeResultToFile()


	tasks := make(chan NbnsScanTaskInfo,taskload)

	t.wg.Add(t.tasknum)
	for gr:=1;gr<= t.tasknum;gr++ {
		go t.worker(tasks)
	}
	//创建chan生产者
	for _,host := range t.ipList {
		task := NbnsScanTaskInfo{
			Host:host,
			Port:t.port,
		}
		tasks <- task
	}
	close(tasks)
	t.wg.Wait()
	close(t.resultChan)
	return nil
}

func (t *NbnsScan)worker(tasks chan NbnsScanTaskInfo){
	defer t.wg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}

		hostname := t.sendNbns(task.Host,task.Port)
		if  hostname != ""{
			result := fmt.Sprintf("[NBNS]\t%s:%s\t%s\n",task.Host,task.Port,hostname)
			t.resultChan <- result
		}
	}
}

// 写文件
func (t *NbnsScan)writeResultToFile() {
	var f *os.File
	var err error
	f, err = os.OpenFile(t.resultsOutput, os.O_RDWR|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for  {
		res,ok :=<- t.resultChan
		if !ok {
			return
		}
		_, _ = f.WriteString(res)
	}

}



//ParseIP(arp.SourceProtAddress)
func (t *NbnsScan)sendNbns(SERVER_IP string, SERVER_PORT string) string {
	serverAddr := SERVER_IP + ":" + SERVER_PORT

	udpConn, err := net.DialTimeout("udp", serverAddr,1*time.Second)
	if err !=nil{
		return ""
	}

	defer udpConn.Close()
	len, err := udpConn.Write(t.nbnsPayload)
	if err != nil{
		return ""
	}
	readTimeout := 2*time.Second
	buf := make([]byte, 2048)
	err = udpConn.SetReadDeadline(time.Now().Add(readTimeout)) // timeout
	if err != nil {
		return ""
	}
	
	len, _ = udpConn.Read(buf)
	if len == 0 {
		return ""
	}
	hostname := ParseNBNS(buf)
	return hostname
}

func ParseNBNS(data []byte) string {
	var buf bytes.Buffer
	i := bytes.Index(data, []byte{0x20, 0x43, 0x4b, 0x41, 0x41})
	if i < 0 || len(data) < 32 {
		return ""
	}
	index := i + 1 + 0x20 + 12
	// data[index-1]是在 number of names 的索引上，如果number of names 为0，退出
	if data[index-1] == 0x00 {
		return ""
	}
	for t:= index; ; t++ {
		// 0x20 和 0x00 是终止符
		if data[t] == 0x20 || data[t] == 0x00 {
			break
		}
		buf.WriteByte(data[t])
	}
	return buf.String()
}

// 根据ip生成含mdns请求包，包存储在 buffer里
func nbns(buffer *Buffer) {
	rand.Seed(time.Now().UnixNano())
	tid := rand.Intn(0x7fff)
	b := buffer.PrependBytes(12)
	binary.BigEndian.PutUint16(b, uint16(tid)) // 0x0000 标识
	binary.BigEndian.PutUint16(b[2:], uint16(0x0010)) // 标识
	binary.BigEndian.PutUint16(b[4:], uint16(1)) // 问题数
	binary.BigEndian.PutUint16(b[6:], uint16(0)) // 资源数
	binary.BigEndian.PutUint16(b[8:], uint16(0)) // 授权资源记录数
	binary.BigEndian.PutUint16(b[10:], uint16(0)) // 额外资源记录数
	// 查询问题
	b = buffer.PrependBytes(1)
	b[0] = 0x20
	b = buffer.PrependBytes(32)
	copy(b, []byte{0x43, 0x4b})
	for i:=2; i<32; i++ {
		b[i] = 0x41
	}

	b = buffer.PrependBytes(1)
	// terminator
	b[0] = 0
	// type 和 classIn
	b = buffer.PrependBytes(4)
	binary.BigEndian.PutUint16(b, uint16(33))
	binary.BigEndian.PutUint16(b[2:], 1)
}

type Buffer struct {
	data  []byte
	start int
}

func (b *Buffer) PrependBytes(n int) []byte {
	length := cap(b.data) + n
	newData := make([]byte, length)
	copy(newData, b.data)
	b.start = cap(b.data)
	b.data = newData
	return b.data[b.start:]
}

func NewBuffer() *Buffer {
	return &Buffer{

	}
}