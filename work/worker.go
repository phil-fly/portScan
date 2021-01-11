package work

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"portScan/utils/ipParse"
	"strings"
	"time"
)

func scan(s ServerAndPort,iplist []string,ResultsOutput string,Tasknum int) error {
	if s.Enable == false {
		return nil
	}
	switch s.Server {
	case TCP :
		tcpscan := &TcpScan{}
		tcpscan.SetIpList(iplist)
		tcpscan.SetPortMap(String2PortMap(s.ServerPort))
		tcpscan.SetResultsOutput(ResultsOutput)
		tcpscan.SetTasknum(Tasknum)
		err := tcpscan.Validate()
		if err != nil {
			return err
		}
		tcpscan.RunScan()
		return nil
	case HTTPS:
		httpsScan := &WebScan{}
		httpsScan.SetIpList(iplist)
		httpsScan.SetPortMap(String2PortMap(s.ServerPort))
		httpsScan.SetResultsOutput(ResultsOutput)
		httpsScan.SetTasknum(Tasknum)
		httpsScan.SetTimeOut(s.TimeOut)
		httpsScan.IsHttps = true
		err := httpsScan.Validate()
		if err != nil {
			return err
		}
		httpsScan.RunScan()
		return nil
	case HTTP:
		httpsScan := &WebScan{}
		httpsScan.SetIpList(iplist)
		httpsScan.SetPortMap(String2PortMap(s.ServerPort))
		httpsScan.SetResultsOutput(ResultsOutput)
		httpsScan.SetTasknum(Tasknum)
		httpsScan.SetTimeOut(s.TimeOut)
		httpsScan.IsHttps = false
		err := httpsScan.Validate()
		if err != nil {
			return err
		}
		httpsScan.RunScan()
		return nil
	case NBNS:
		nbnsScan := &NbnsScan{}
		nbnsScan.SetIpList(iplist)
		nbnsScan.SetResultsOutput(ResultsOutput)
		nbnsScan.SetTasknum(Tasknum)
		nbnsScan.SetTimeOut(s.TimeOut)
		err := nbnsScan.Validate()
		if err != nil {
			return err
		}
		nbnsScan.RunScan()
		return nil
	}
	return nil
}

func ScanEngine(ss ScanServerAndPort) error {
	//log.Println(ips)
	startTime := time.Now()
	err := ss.Validate()
	if err != nil {
		return err
	}
	iplist := file2Iplist(true,ss.TargetFile)

	ipscan := NewIpScan()
	ipscan.SetIpList(iplist)
	ipscan.tasknum = 10000

	ipscan.RunScan()
	take1 := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("ipscan completed, taking %s.\n\n", take1)
	fmt.Println("总地址个数:",len(ipscan.ipList))
	fmt.Println("存活主机地址个数:",len(ipscan.aliveIpList))

	for _,v := range ss.ServerAndPorts{
		scan(v,ipscan.aliveIpList,ss.ResultsFile,ss.Tasknum)
	}
	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("Scanning completed, taking %s.\n\n", takes)
	return nil
}

func file2Iplist(fileMode bool,ipStr string) []string {
	var IPset []string
	var err error
	if fileMode {
		request,err := ioutil.ReadFile(ipStr)
		if err !=nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		ipStr := strings.Replace(string(request), "\r", "", -1 )
		ports := strings.Split(ipStr, "\n")
		for _,v := range ports {
			if v != ""{
				ips, err := ipParse.Parse(v)
				if err != nil {
					IPset = append(IPset, v)
				}
				IPset = append(IPset, ips...)
			}
		}
	}else{
		IPset, err = ipParse.Parse(ipStr)
		if err != nil {
			log.Fatal(err)
		}
	}
	return IPset
}