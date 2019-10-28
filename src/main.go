package main

import (
	"flag"
	"net"
	"log"
	"fmt"
	neturl "net/url"
	"strings"
	"bufio"
)

func main()  {
	var streamUrl string
	flag.StringVar(&streamUrl, "url", "rtmp://127.0.0.1/live/xafanda", "use -url to specify stream url, default is rtmp://127.0.0.1/live/xafanda")
	flag.Parse()

	localAddr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		log.Println(fmt.Sprintf("resolve local addr failed, err is %v", err))
		return
	}

	u, err := neturl.Parse(streamUrl)
	if err != nil {
		log.Println(fmt.Sprintf("parse stream url:%v failed, err is %v", streamUrl, err))
		return
	}
	port := ":1935"
	host := u.Host
	id := strings.Index(host, ":")
	if id != -1 {
		port = host[id:]
		host = host[:id]
	}

	hp := host + port
	remoteAddr, err := net.ResolveTCPAddr("tcp", hp) // remote host:port, default port is 1935
	if err != nil {
		log.Println(fmt.Sprintf("resolve remote addr:%v failed, err is %v", hp, err))
		return
	}
	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		log.Println(fmt.Sprintf("dial failed, err is %v", err))
		return
	}
	log.Println(fmt.Sprintf("dial tcp from %v to %v success", localAddr.String(), remoteAddr.String()))
	defer conn.Close()


	size := 4 * 1024
	rw := bufio.NewReadWriter(bufio.NewReaderSize(conn, size), bufio.NewWriterSize(conn, size))
	// tcp handshake
	if err := Handshake(rw); err != nil {
		log.Println(fmt.Sprintf("do handshake failed, err is %v", err))
		return
	}

	// command message
	WriteConnect(rw)
	// conn
	// create stream
	// play
	// recv streams
}