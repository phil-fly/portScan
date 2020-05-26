package main

import (
	"runtime"
	"fmt"
	"flag"
	"portScan/utils/ipParse"
	"log"
	"portScan/work"
	"time"
	"strings"
	"io/ioutil"
	"os"
)
var (
	maxThread     int
	fullMode      bool
	fileMode	  bool
	specifiedPort string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	banner := `

                               $$\      $$$$$$\                               
                               $$ |    $$  __$$\                              
 $$$$$$\   $$$$$$\   $$$$$$\ $$$$$$\   $$ /  \__| $$$$$$$\ $$$$$$\  $$$$$$$\  
$$  __$$\ $$  __$$\ $$  __$$\\_$$  _|  \$$$$$$\  $$  _____|\____$$\ $$  __$$\ 
$$ /  $$ |$$ /  $$ |$$ |  \__| $$ |     \____$$\ $$ /      $$$$$$$ |$$ |  $$ |
$$ |  $$ |$$ |  $$ |$$ |       $$ |$$\ $$\   $$ |$$ |     $$  __$$ |$$ |  $$ |
$$$$$$$  |\$$$$$$  |$$ |       \$$$$  |\$$$$$$  |\$$$$$$$\\$$$$$$$ |$$ |  $$ |
$$  ____/  \______/ \__|        \____/  \______/  \_______|\_______|\__|  \__|
$$ |                                                                          
$$ |                                                                          
\__|                                                                          

`
	ip := ""
	specifiedPort = ""
	fileMode = false
	fullMode = false
	fmt.Println(banner)
	flag.StringVar(&ip, "ip", "", "IP to be scanned, supports three formats:\n192.168.0.1 \n192.168.0.1-8 \n192.168.0.0/24")
	flag.BoolVar(&fullMode, "full", false, "Scan all TCP and UDP ports in full scan mode. The default is off. By default, only common TCP ports are scanned.")
	flag.BoolVar(&fileMode, "file", false, "Use file mode to specify ip address .")
	flag.IntVar(&maxThread, "t", 10000, "Maximum number of threads")
	flag.StringVar(&specifiedPort, "p", "", "Port to be scanned, supports three formats:\n22,80 \n22-65535")
	flag.Parse()

	if ip == "" {
		flag.Usage()
		return
	}

	ips:=ipHandle(fileMode,ip)
	//log.Println(ips)
	startTime := time.Now()
	if len(specifiedPort) > 0 {
		if fullMode {
			fmt.Println("Multi-host mode does not support full scan")
			return
		}
		Ports :=work.PortHandle(specifiedPort)
		//log.Println(Ports)
		work.Task(ips,Ports,maxThread)
	} else {
		var commonPorts string
		if fullMode {
			commonPorts = "1-65535"
		} else {
			commonPorts = "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"

		}
		Ports :=work.PortHandle(commonPorts)
		work.Task(ips,Ports,maxThread)
	}
	takes := time.Since(startTime).Truncate(time.Millisecond)
	fmt.Printf("Scanning completed, taking %s.\n\n", takes)
}

func ipHandle(fileMode bool,ipStr string) []string {
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
					log.Fatal(err)
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