package main

import (
	"github.com/fsouza/go-dockerclient"
	"os"
	"strconv"
	"strings"
	"net/http"
	"log"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	proxy *httputil.ReverseProxy
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.proxy.ServeHTTP(w, req)
}

func NewProxyHandler(proxyUrl *url.URL) *ProxyHandler {
	return &ProxyHandler{httputil.NewSingleHostReverseProxy(proxyUrl)}
}

func newURL(urlStr string) *url.URL {
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal("url.Parse: ", err)
	}
	return url
}

func main() {
	host := os.Getenv("DOCKER_HOST")
	hosturl, _ := url.Parse(host)
	ip := strings.Split(hosturl.Host, ":")[0]
	log.Printf("docker host address: %s", ip)
	client, _ := docker.NewClientFromEnv()
	containers, _ := client.ListContainers(docker.ListContainersOptions{})
	quit := make(chan int)
	for _, c := range containers {
		for _, p := range c.Ports {
			port := ":" + strconv.FormatInt(p.PublicPort, 10)
			address := "http://" + ip + port
			log.Printf("discovered container port: %s", address)

			proxy := NewProxyHandler(newURL(address))
			go func() {
				err := http.ListenAndServe(port, proxy)
				if err != nil {
					log.Fatal("ListenAndServe: ", err)
				}
				quit <- 0
			}()
		}
	}
	<-quit
}
