package main

import (
	"net"
	"./util"
	"runtime"
	"./serverside"
	"strconv"
	"./config"
)

func init() {
	runtime.GOMAXPROCS(config.GetCores())
	util.SetupRegister()		
}

func main(){
	go serverside.CheckForfile()
	go serverside.LoadTorrentFile()
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":" + strconv.Itoa(config.GetPort()))
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

