package work

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ServerType	string

var (
	TCP		ServerType	= "TCP"
	UDP		ServerType	= "UDP"


	HTTPS	ServerType	= "HTTPS"
	HTTP	ServerType	= "HTTP"
	FTP	ServerType	= "FTP"
	SMB	ServerType	= "SMB"
	NBNS	ServerType	= "NBNS"
)

func (t ServerType)String() string {
	switch t {
	case HTTPS:
		return "HTTPS"
	case HTTP:
		return "HTTP"
	case FTP:
		return "FTP"
	case SMB:
		return "SMB"
	case NBNS:
		return "NBNS"

	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	}
	return ""
}

func String2ServerType(serverType string) ServerType {
	switch serverType {
	case "HTTPS":
		return "HTTPS"
	case "HTTP":
		return "HTTP"
	case "FTP":
		return "FTP"
	case "SMB":
		return "SMB"
	case "NBNS":
		return "NBNS"

	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	}
	return ""
}

type ServerAndPort struct {
	Enable	bool	`json:"Enable"`
	Server	ServerType	`json:"ServerType"`
	ServerPort string	`json:"ServerPort"`
	TimeOut	int	`json:"TimeOut"`
}

type ScanServerAndPort struct {
	TargetFile	string	`json:"TargetFile"`
	ResultsFile	string	`json:"ResultsFile"`
	Tasknum	int	`json:"Tasknum"`
	ServerAndPorts []ServerAndPort 	`json:"ServerAndPorts"`
}

const ScanList =  "list.txt"

func (t *ScanServerAndPort)Validate() error {
	switch {
	case t.Tasknum == 0:
		t.Tasknum = 1000
	case t.TargetFile == "":
		t.TargetFile = ScanList
	}
	file,err :=os.Open(t.ResultsFile)
	defer file.Close()
	if err !=nil && os.IsNotExist(err) {
		file ,_= os.Create(t.ResultsFile)
	}

	file,err =os.Open(t.TargetFile)
	defer file.Close()
	if err !=nil && os.IsNotExist(err) {
		file ,_= os.Create(t.TargetFile)
	}
	return nil
}


func (t *ScanServerAndPort)Default(file string) {
	t.Tasknum = 1000
	t.TargetFile = ScanList
	t.ResultsFile = "results.txt"

	tcpSan := ServerAndPort{
		Enable: true,
		Server: TCP,
		ServerPort: "1-65535",
		TimeOut: 1,
	}
	t.ServerAndPorts = append(t.ServerAndPorts, tcpSan)

	tcpSan.Enable = false
	tcpSan.ServerPort = "21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017"
	t.ServerAndPorts = append(t.ServerAndPorts, tcpSan)



	httpSan := ServerAndPort{
		Enable: false,
		Server: HTTP,
		ServerPort: "80,8080",
		TimeOut: 10,
	}
	t.ServerAndPorts = append(t.ServerAndPorts, httpSan)

	httpsSan := ServerAndPort{
		Enable: true,
		Server: HTTPS,
		ServerPort: "443,8443",
		TimeOut: 10,
	}
	t.ServerAndPorts = append(t.ServerAndPorts, httpsSan)

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("JSON ERR:", err)
	}
	ioutil.WriteFile(file,b,666)
}

const (
	taskload		    = 60000
)