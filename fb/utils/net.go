package utils

import (
	"net/http"
	"io/ioutil"
"github.com/labstack/echo"
"net"
	"fmt"
	"strings"
)

func HttpGetString(url string)(int, []byte, error){
	res, err := http.Get(url)
	if err != nil {
		return 500, nil,err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode,nil,err
	}
	return  res.StatusCode, result,nil
}

//获取访问ip地址
func GetRemoteIP(r *http.Request)string{

	remoteAddr := r.RemoteAddr

	if ip := r.Header.Get(echo.HeaderXRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get(echo.HeaderXForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	fmt.Println("ii:",remoteAddr)
	return remoteAddr
}


//获取外网IP
func GetExternalIP() (string,error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	return strings.TrimSpace(string(result)),nil
}

//域名获取ip
//例子 Domain2IP("google.com:80")
func Domain2IP(domain string)string{
	conn, err := net.Dial("udp", domain)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	return  strings.Split(conn.LocalAddr().String(), ":")[0]
}