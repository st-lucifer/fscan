package common

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ParseIPErr  =errors.New("host parsing error\n" +
	"format: \n"+
	"192.168.1.1\n" +
	"192.168.1.1/8\n"+
	"192.168.1.1/16\n"+
	"192.168.1.1/24\n"+
	"192.168.1.1,192.168.1.2\n" +
	"192.168.1.1-255")

func ParseIP(ip string,filename string)(hosts []string,err error){

	if ip != ""{
		hosts,err = ParseIPs(ip)
	}
	if filename != ""{
		var filehost []string
		filehost,_ = Readipfile(filename)
		hosts = append(hosts,filehost...)
	}
	hosts = RemoveDuplicate(hosts)
	return hosts,err
}

func ParseIPs(ip string)(hosts []string,err error){
	if strings.Contains(ip,","){
		IPList:=strings.Split(ip,",")
		var ips []string
		for _,ip:=range IPList{
			ips,err = ParseIPone(ip)
			CheckErr(ip,err)
			hosts = append(hosts,ips...)
		}
		return hosts,err
	}else {
		hosts,err = ParseIPone(ip)
		CheckErr(ip,err)
		return hosts,err
	}
}

func ParseIPone(ip string)([]string,error){
	reg:=regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case strings.Contains(ip[len(ip)-3:len(ip)],"/24"):
		return ParseIPA(ip)
	case strings.Contains(ip[len(ip)-3:len(ip)],"/16"):
		return ParseIPD(ip)
	case strings.Contains(ip[len(ip)-2:len(ip)],"/8"):
		return ParseIPE(ip)
	case strings.Count(ip,"-")==1:
		return ParseIPC(ip)
	case reg.MatchString(ip):
		_, err := net.LookupHost(ip)
		if err != nil {
			return nil,err
		}
		return []string{ip},nil
	default:
		testIP:=net.ParseIP(ip)
		if testIP==nil{
			return nil,ParseIPErr
		}
		return []string{ip},nil
	}
}
//Parsing CIDR IP
func ParseIPA(ip string)([]string,error){
	realIP:=ip[:len(ip)-3]
	testIP:=net.ParseIP(realIP)

	if testIP==nil{
		return nil,ParseIPErr
	}
	IPrange:=strings.Join(strings.Split(realIP,".")[0:3],".")
	var AllIP []string
	for i:=0;i<=255;i++{
		AllIP=append(AllIP,IPrange+"."+strconv.Itoa(i))
	}
	return AllIP,nil
}

//Resolving multiple IPS, for example: 192.168.111.1,192.168.111.2
func ParseIPB(ip string)([]string,error){
	IPList:=strings.Split(ip,",")
	for _,i:=range IPList{
		testIP:=net.ParseIP(i)
		if testIP==nil{
			return nil,ParseIPErr
		}
	}
	return IPList,nil

}

//Resolving a range of IP,for example: 192.168.111.1-255
func ParseIPC(ip string)([]string,error){
	IPRange:=strings.Split(ip,"-")
	testIP:=net.ParseIP(IPRange[0])
	Range,err:=strconv.Atoi(IPRange[1])
	if testIP==nil || Range>255 || err!=nil{
		return nil,ParseIPErr
	}
	SplitIP:=strings.Split(IPRange[0],".")
	ip1,err1:=strconv.Atoi(SplitIP[3])
	ip2,err2:=strconv.Atoi(IPRange[1])
	PrefixIP:=strings.Join(SplitIP[0:3],".")
	var AllIP []string
	if ip1>ip2 || err1!=nil || err2!=nil{
		return nil,ParseIPErr
	}
	for i:=ip1;i<=ip2;i++{
		AllIP=append(AllIP,PrefixIP+"."+strconv.Itoa(i))
	}
	return AllIP,nil

}

func ParseIPD(ip string)([]string,error){
	realIP:=ip[:len(ip)-3]
	testIP:=net.ParseIP(realIP)

	if testIP==nil{
		return nil,ParseIPErr
	}
	IPrange:=strings.Join(strings.Split(realIP,".")[0:2],".")
	var AllIP []string
	for a:=0;a<=255;a++{
		for b:=0;b<=255;b++{
			AllIP=append(AllIP,IPrange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b))
		}
	}
	return AllIP,nil
}

func ParseIPE(ip string)([]string,error){
	realIP:=ip[:len(ip)-2]
	testIP:=net.ParseIP(realIP)

	if testIP==nil{
		return nil,ParseIPErr
	}
	IPrange:=strings.Join(strings.Split(realIP,".")[0:1],".")
	var AllIP []string
	for a:=0;a<=255;a++{
		for b:=0;b<=255;b++{
			AllIP=append(AllIP,IPrange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(1))
			AllIP=append(AllIP,IPrange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(254))
		}
	}
	return AllIP,nil
}

func Readipfile(filename string)([]string,error){
	file, err := os.Open(filename)
	if err!=nil{
		fmt.Println("Open %s error, %v", filename,err)
		os.Exit(0)
	}
	defer file.Close()
	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			host,err := ParseIPs(text)
			CheckErr(text,err)
			content=append(content,host...)
		}
	}
	return content,nil
}


func RemoveDuplicate(old []string) ([]string) {
	result := make([]string, 0, len(old))
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

