package util

/*
The package util is the package for utility functions and for types found thoughout the application, or need to be known application wide.

Author: Clinton Scott brutii@gmail.com
*/

import (
	"os"
	"fmt"
	"net"
	"encoding/json"
)
//The type expected by the server to be sent to it. 
type Message struct{
	Host, Info, Id, Tag string
	Level int
}
//The type expected to sent and recieved in channels
type Event struct{
	Type, Message string
}
//Helper to check and throw an error if required
func CheckError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

//Const messages to be sent around so it easier to refactor
const (
	NewTorrentMagnet = "new_torrent_magnet"
	Start = "start"
	Kill = "kill"
	AddWatchDir = "add_dir"
	RemoveWatchDir = "remove_dir"
	NewTorrentFile = "new_torrent_file"
)
/*
Will flesh out later into something more useful later or write another one to be more useful.
*/

func SendMessage (connection net.Conn){
	var messagechannel chan Event
	
	if CheckForChannel("messagechannel") {
		messagechannel = make(chan Event, 10)
		RegisterChannel("messagechannel", &messagechannel)
	}else{
		messagechannel = *GetChannel("messagechannel")
	}
	
	for {
		m := <- messagechannel
		enc := json.NewEncoder(connection)	
		enc.Encode(m)
	}
}
