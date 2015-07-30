package main

import (
	"net"
	"./util"
	"runtime"
	"./serverside"
)


func main(){
	runtime.GOMAXPROCS(2)
	util.SetupRegister()

	//Threads for file stuff I think is it the right way to go
	go serverside.CheckForfile()
	go serverside.LoadTorrentFile()
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":4638")
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

