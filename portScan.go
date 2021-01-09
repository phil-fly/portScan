package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"portScan/work"
	"runtime"
)

const ScanConfig =  "config.json"


func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ParseFile2Json(jsonFile string)  (ServerAndPorts work.ScanServerAndPort,err error) {
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return ServerAndPorts,err
	}

	err = json.Unmarshal(byteValue, &ServerAndPorts)
	if err != nil {
		return ServerAndPorts,err
	}
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	scan := work.ScanServerAndPort{}
	if !FileExist(ScanConfig) {
		scan.Default(ScanConfig)
	}
	ServerAndPorts,err := ParseFile2Json(ScanConfig)
	if err != nil {

	}
	work.ScanEngine(ServerAndPorts)
}
