package work

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	//"portScan/utils/ping"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"github.com/saintfish/chardet"
)

type WebScan struct {
	ipList	[]string	`json:"ipList"`
	portMap *PortSet	`json:"portMap"`
	tasknum	int	`json:"tasknum"`
	scanWg sync.WaitGroup
	resultsOutput 	string 	`json:"resultsOutput"`
	resultChan	chan string
	timeOut	int	`json:"timeOut"`
	IsHttps	bool
}

type WebScanTaskInfo struct {
	Host	string
	Port    string
	IsHttps	bool
}

func (t *WebScan)SetIpList(ips []string){
	t.ipList = ips
}

func (t *WebScan)SetTimeOut(timeOut int){
	t.timeOut = timeOut
}

func (t *WebScan)SetPortMap(portMap *PortSet){
	t.portMap = portMap
}

func (t *WebScan)SetTasknum(tasknum int){
	t.tasknum = tasknum
}

func (t *WebScan)SetIsHttps(isHttps bool){
	t.IsHttps = isHttps
}

func (t *WebScan)SetResultsOutput(resultsOutput string) {
	t.resultsOutput = resultsOutput
	file,err :=os.Open(resultsOutput)
	defer file.Close()
	if err !=nil && os.IsNotExist(err) {
		file ,_= os.Create(resultsOutput)
	}
}

func (t *WebScan)Validate() error {
	switch {
	case len(t.ipList)< 1 :
		return errors.New("WebScan ipList is nil.")
	case t.tasknum ==  0:
		t.tasknum = 1000
	case t.timeOut == 0:
		t.timeOut = 10
	case t.portMap == nil :
		return errors.New("WebScan portMap is nil.")
	}
	return nil
}

func (t *WebScan)RunScan() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	//结果写入文件
	t.resultChan = make(chan string,1000)
	go t.writeResultToFile()


	tasks := make(chan WebScanTaskInfo,taskload)
	t.scanWg.Add(t.tasknum)
	for gr:=1;gr<= t.tasknum;gr++ {
		go t.worker(tasks)
	}

	for _,host := range t.ipList {

		for Port,_ := range t.portMap.Port {
			//fmt.Println("扫描地址:",host,Port)
			task := WebScanTaskInfo{
				Host:host,
				Port:Port,
				IsHttps: t.IsHttps,
			}
			tasks <- task
		}
	}
	close(tasks)
	t.scanWg.Wait()
	close(t.resultChan)
	return nil
}

func (t *WebScan)worker(tasks chan WebScanTaskInfo){
	defer t.scanWg.Done()
	for {
		task,ok := <- tasks
		if !ok {
			return
		}
		var url string
		if task.IsHttps {
			url = fmt.Sprintf("https://%s:%s/",task.Host,task.Port)
		}else{
			url = fmt.Sprintf("http://%s:%s/",task.Host,task.Port)
		}
		t.getTitle(url)
	}
}

// 写文件
func (t *WebScan)writeResultToFile() {
	var f *os.File
	var err error

	f, err = os.OpenFile(t.resultsOutput, os.O_RDWR|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for  {
		res,ok:=<- t.resultChan
		if !ok {
			fmt.Println("resultChan !ok")
			return
		}

		_, err = f.WriteString(res)
		if err !=nil{
			fmt.Println(err)
		}
	}

}

func (t *WebScan)getTitle(url string) {
	//fmt.Println("getTitle:" ,url)
	timeout := time.Duration(t.timeOut) * time.Second
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 设置请求头
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(request)
	if err != nil {
	//	fmt.Println(err)
		return
	}

	lv := response.Header.Get("Location")
	if lv != ""  {
		fmt.Println(lv)
	}

	defer response.Body.Close()
	statusCode := response.StatusCode
	//fmt.Println(url,statusCode)
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find("title").Text()

	var result string
	// 判断title是否为gbk编码
	detector := chardet.NewTextDetector()
	charset, _ := detector.DetectBest([]byte(title))

	if charset.Charset=="UTF-8" {
		result = fmt.Sprintf("%s\t%d\t%s\n", url, statusCode, title)

	} else {
		gbkTitle, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(title))
		result = fmt.Sprintf("%s\t%d\t%s\n", url, statusCode, string(gbkTitle))
	}
	// if isGBK([]byte(title)) {
	//      gbkTitle, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(title))
	//      output = fmt.Sprintf("%s\t%d\t%s", url, statusCode, string(gbkTitle))
	// } else {
	//      output = fmt.Sprintf("%s\t%d\t%s", url, statusCode, title)
	// }
	//fmt.Println(result)
	t.resultChan <- result
}
