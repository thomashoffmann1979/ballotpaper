package api

import (
	"io"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
	"encoding/json"
	"log"

)

var Cookies []http.Cookie
var Jar *cookiejar.Jar

var timeout = time.Duration(10 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func InitJar() {
	if Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{ })
		if err != nil {
			log.Fatal(err)
		}
		Jar = jar
		fmt.Println("Jar initialized")
	}
}

func Get(url string) (string, error) {
	InitJar()
	transport := http.Transport{
		Dial: dialTimeout,
	}
	client := http.Client{
		Transport: &transport,
		Jar: Jar,
	}
	var resp *http.Response
	var err error
	var body []byte
	resp, err = client.Get(url)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Post(url string, data string) (string, error) {
	InitJar()
	transport := http.Transport{
		Dial: dialTimeout,
	}
	client := http.Client{
		Transport: &transport,
		Jar: Jar,
	}
	var resp *http.Response
	var err error
	var body []byte
	resp, err = client.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Login(url string, username string, password string) (LoginResponse, error) {
	var loginResponse LoginResponse
	sb,err := Post(url, "forcelogin=1&username="+username+"&password="+password+"")
	json.Unmarshal([]byte(sb), &loginResponse)
	return loginResponse, err
}


func Ping(url string) (PingResponse, error) {
	var response PingResponse
	sb,err := Get(url+"dashboard/ping")
	json.Unmarshal([]byte(sb), &response)
	return response, err
}


func GetKandidaten(url string) (KandidatenResponse, error) {
	var response KandidatenResponse
	sb,err := Get(url+"ds/kandidaten/read")
	json.Unmarshal([]byte(sb), &response)
	return response, err
}
