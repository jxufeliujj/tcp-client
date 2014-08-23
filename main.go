package main

import (
	"flag"
	"net"
	"tcp-client/socket"
)

func main() {
	flag.Parse()
	service := flag.Arg(0)

	if service == "" {
		service = "127.0.0.1:8080"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	socket.Start(conn)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
