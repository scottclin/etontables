package main

import (
	"net"
	"os"
	"fmt"
	"./config"
	"./util"
	"./clientside/console"
)

func main(){

	server_ip, err := net.ResolveTCPAddr("tcp", config.GetServerIP())
	util.CheckError(err)
	
	conn, err := net.DialTCP("tcp", nil, server_ip)

	util.CheckError(err)

//	conn.SetkeepAlive(true)

	fmt.Println("I have connected")

	util.CheckError(err)
//Thread these things
	go util.SendMessage(conn)

	go console.ClientConsole()
//This reads a message from the server prints it and then dies. Yeah bad but it from the primodial soup code...
	var buf [512]byte

	fmt.Println("Reading...")
	
	_, err = conn.Read(buf[0:])
	if err != nil {
		return
	}
	
	fmt.Println(string(buf[0:]))

	os.Exit(0)	
}
