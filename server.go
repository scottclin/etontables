package main

import (
	"./config"
	"./serverside"
	"./util"
	"net"
	"runtime"
	"strconv"
)

func init() {
	runtime.GOMAXPROCS(config.GetCores())
	util.SetupRegister()
}

func main() {
	client := serverside.Start()

	go serverside.InfoStart()
	go serverside.CheckForfile()
	go serverside.Control(client)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(config.GetPort()))
	util.CheckError(err)

	ln, err := net.ListenTCP("tcp", tcpAddr)

	util.CheckError(err)

	for {
		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go serverside.HandleClient(conn)
	}

}
