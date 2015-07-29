package main

import (
	"net"
	"os"
	"fmt"
	"./config"
	"./util"
	"./console"
)

func main(){

	server_ip, err := net.ResolveTCPAddr("tcp", config.GetServerIP())
	util.CheckError(err)
	
	conn, err := net.DialTCP("tcp", nil, server_ip)

	util.CheckError(err)

//	conn.SetkeepAlive(true)

	fmt.Println("I have connected")

	util.CheckError(err)

	go util.SendMessage(conn)

	go console.ClientConsole()

	var buf [512]byte

	fmt.Println("Reading...")
	
	_, err = conn.Read(buf[0:])
	if err != nil {
		return
	}
	
	fmt.Println(string(buf[0:]))

	os.Exit(0)	
}
