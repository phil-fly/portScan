package ipParse

import (
	"strings"
	"strconv"
	"regexp"
	"net"
	"errors"
	"fmt"
)

//  192.168.0.1
//  192.168.0.1-100
//  192.168.0.0/24

var (
	rePureIP  *regexp.Regexp
	reCIDR    *regexp.Regexp
	reIPRange *regexp.Regexp
)

func init() {
	rePureIP = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	reCIDR = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{2}$`)
	reIPRange = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}-\d{1,3}$`)
}

// IP2Long convert a ip to a uint.
func IP2Long(ip net.IP) uint {
	return (uint(ip[12]) << 24) + (uint(ip[13]) << 16) + (uint(ip[14]) << 8) + uint(ip[15])
}

// Long2IP convert a uint to a ip.
func Long2IP(long uint) net.IP {
	return net.IPv4(byte(long>>24), byte(long>>16), byte(long>>8), byte(long))
}

// CIDR2IPS returns a ip list of CIDR.
func CIDR2IPS(ipr string) (ips []string, err error) {
	ip, ipnet, err := net.ParseCIDR(ipr)
	if err != nil {
		return nil, err
	}
	if !ip.Equal(ipnet.IP) {
		return nil, errors.New("Invalid CIDR")
	}

	var mask int
	mask, _ = strconv.Atoi(strings.Split(ipr, "/")[1])
	if mask < 16 {
		return nil, errors.New("Invalid Mask(too small)")
	}

	ipstart := IP2Long(ip)
	hostCount := (1 << uint(32-mask)) - 1
	for i := 0; i <= hostCount; i++ {
		ips = append(ips, Long2IP(ipstart+uint(i)).String())
	}
	return ips, nil
}

// Range2IPS returns a ip list of a iprange
func Range2IPS(ipr string) (ips []string, err error) {
	if !reIPRange.MatchString(ipr) {
		return nil, errors.New("Invalid ip-range")
	}
	tmp := strings.Split(ipr, "-")[0]
	if net.ParseIP(tmp) == nil {
		return nil, errors.New("Invalid start ip")
	}

	tmpIP := strings.Split(ipr, ".")
	prefix := tmpIP[0] + "." + tmpIP[1] + "." + tmpIP[2] + "."
	start, _ := strconv.Atoi(strings.Split(tmpIP[3], "-")[0])
	end, _ := strconv.Atoi(strings.Split(tmpIP[3], "-")[1])
	if end < start || end > 255 {
		return nil, errors.New("Invalid end ip")
	}
	for i := start; i <= end; i++ {
		ip := prefix + strconv.Itoa(i)
		ips = append(ips, ip)
	}
	return ips, nil
}

// Parse parse a ip such as ip,CIDR,ip-range and returns ip list
func Parse(ip string) ([]string, error) {
	if rePureIP.MatchString(ip) {
		return []string{ip}, nil
	} else if reCIDR.MatchString(ip) {
		return CIDR2IPS(ip)
	} else if reIPRange.MatchString(ip) {
		return Range2IPS(ip)
	} else {
		return nil, errors.New("Invalid IP")
	}
}


func main(){
	iplist := "10.10.1.20-255"
	ips,err := Parse(iplist)
	if err !=nil {
		fmt.Println(err)
		return
	}
	ips,err = Parse("10.10.1.0/24")
	if err !=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ips)
}
