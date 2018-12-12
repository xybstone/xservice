package main

import (
	"github.com/xybstone/xservice/service"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":6001")
	service.Run(listener)
}
