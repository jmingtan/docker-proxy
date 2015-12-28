package main

import (
	"github.com/fsouza/go-dockerclient"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func copyData(src net.Conn, dst net.Conn) {
	_, _ = io.Copy(src, dst)
	src.Close()
	dst.Close()
}

func makeProxy(port, address string) {
	ln, err := net.Listen("tcp", port)
	if err == nil {
		go func() {
			for {
				srcConn, srcErr := ln.Accept()
				dstConn, dstErr := net.Dial("tcp", address)
				if srcErr == nil && dstErr == nil {
					go copyData(srcConn, dstConn)
					go copyData(dstConn, srcConn)
				} else {
					log.Println("src:", srcErr, ", dst:", dstErr)
				}
			}
		}()
	}
}

func main() {
	host := os.Getenv("DOCKER_HOST")
	hosturl, _ := url.Parse(host)
	ip := strings.Split(hosturl.Host, ":")[0]
	log.Printf("docker host address: %s", ip)
	quit := make(chan int)
	client, _ := docker.NewClientFromEnv()
	containers, _ := client.ListContainers(docker.ListContainersOptions{})
	for _, c := range containers {
		for _, p := range c.Ports {
			port := ":" + strconv.FormatInt(p.PublicPort, 10)
			address := ip + port
			log.Printf("discovered container port: %s", address)
			makeProxy(port, address)
		}
	}
	<-quit
}
