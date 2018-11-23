package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

/**
定义一个指针结构体指向自定义的ServeHTTP
 */
type MyMux struct {
}

//ip结构体
type SIp struct {
	//上一次检查的时间
	lastCheckTime int
	//ip地址
	ip string
	//检查次数
	checkTimes int
	//状态
	available bool
}

//ip结构体列表
var gIpSlice = make([]SIp, 0)

//通过http请求的方式获取数据
func httpMethod(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Println("param:", request.Form)    //调试输出请求参数信息
	fmt.Println("path:", request.URL.Path) //调试输出请求url信息

	httpClient := http.Client{};
	req, err := http.NewRequest("GET", "http://"+getAvaiableIp()+request.URL.Path, nil)
	if (err != nil) {
		fmt.Println(err)
	}
	req.Host = "steam.zabite.net"
	resp, err := httpClient.Do(req)
	if (err != nil) {
		fmt.Println(err)
	}
	defer resp.Body.Close();
	b, _ := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		fmt.Println("pass\n")
	case 429:
		fmt.Println("429\n")
	case 500:
		fmt.Println("500\n")
	}
	fmt.Fprint(responseWriter, string(b))

}

func getAvaiableIp() (ip string) {
	for _, value := range gIpSlice {
		if true == value.available {
			return value.ip
		}
	}
	return ""
}

//处理某个ip 429的情况
func process429(ip string) {

}

//处理某个ip 500的情况
func process500(ip string) {

}

func httpGet() {

}

//通过tcp获取数据,未实现
func tcpMethod(responseWriter http.ResponseWriter, request http.Request) {
	request.ParseForm()
	fmt.Println("param:", request.Form)    //调试输出请求参数信息
	fmt.Println("path:", request.URL.Path) //调试输出请求url信息

	for k, v := range request.Form {
		fmt.Println("key:", k)
		fmt.Println("value", strings.Join(v, "")) //将value处理成string
	}

	client, err := net.Dial("tcp", "steam.zabite.net:http") //47.90.110.198:80  127.0.0.1:8849
	if err != nil {
		fmt.Print(err)
		return
	}

	defer client.Close()
	buf := make([]byte, 16)
	str := "GET /inventory/76561198828839992/570/2 HTTP/1.1\r\n" +
		"Host: steam.zabite.net\r\n" +
		"User-Agent: curl/7.62.0r\n" +
		"Accept: */*\r\n"
	io.CopyBuffer(client, bytes.NewBuffer([]byte(str)), buf)

	//获取服务端数据
	resp := bytes.NewBuffer(make([]byte, 0))
	io.CopyBuffer(resp, client, buf)
	fmt.Fprint(responseWriter, resp)
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//访问路径为根目录
	httpMethod(w, r)
	return
}

func main() {
	mux := &MyMux{}
	go checkLoop()
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func init() {
	temp := SIp{0, "47.90.110.198", 0, true}
	gIpSlice = append(gIpSlice, temp)
	fmt.Println(len(gIpSlice))
}

//检查线程死循环
func checkLoop() {
	for {
		fmt.Println("i am checking")
		time.Sleep(1 * time.Second)
		fmt.Println(len(gIpSlice))
		for i := 0; i < len(gIpSlice); i++ {
			sIp := gIpSlice[i]
			fmt.Println(sIp.ip, sIp.lastCheckTime)
		}

	}
}

//使用此ip去尝试访问steam,获取状态码和返回内容
func get(url string,method string)(httpCode int,body string) {
	httpClient := http.Client{};
	req, err := http.NewRequest("GET", "http://"+ip+"/inventory/76561198828839992/570/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Host = "steam.zabite.net"
	resp, err := httpClient.Do(req)
	if (err != nil) {
		fmt.Println(err)
	}
	defer resp.Body.Close();
	b, _ := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		fmt.Println("pass\n")
	case 429:
		fmt.Println("429\n")
	case 500:
		fmt.Println("500\n")
	}
	return resp.StatusCode,string(b)

}
