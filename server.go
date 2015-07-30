package main

import (
	"net"
	"encoding/json"
	"fmt"
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

		go handleClient(conn)
	}
}


/**
This handles communications with the clients, recieves messages and sends them to the right channels 
*/
func handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()

	dec := json.NewDecoder(conn)
	for {
		var m util.Message
		if err := dec.Decode(&m); err != nil {
			break
		}

		fmt.Printf("%s, %d, %s: %s %s\n", m.Host, m.Level, m.Id, m.Info, m.Tag)

		_, err := conn.Write([]byte("Message received"))

		switch m.Tag {
		case "set_watch_folder": watch_dir := m.Info
			fmt.Printf("Set the watch dir to be %s\n", watch_dir)
		default: fmt.Printf("Unreconised tag sent: %s\n", m.Tag)
		}
		util.CheckError(err)
				
	}
}
