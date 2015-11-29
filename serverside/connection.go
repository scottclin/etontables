package serverside

import (
	"../util"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

/**
This handles communications with the clients, recieves messages and sends them to the right channels
*/
func HandleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()

	//Get user action channel
	userEventChannel := util.GetChannel("userEventChannel")
	if userEventChannel == nil {
		fmt.Fprintf(os.Stderr, "Failed to get the userEventChannel, I think it should exist.\n")
	}

	dec := json.NewDecoder(conn)
	for {
		var m util.Message
		if err := dec.Decode(&m); err != nil {
			break
		}

		fmt.Printf("%s, %d, %s: %s %s\n", m.Host, m.Level, m.Id, m.Info, m.Tag)

		_, err := conn.Write([]byte("Message received"))

		switch m.Tag {
		case "get_info":
			fmt.Printf("Asked for info, sending....\n")
			
			fmt.Printf("Info sent\n")			
		default:
			userEventChannel <- util.Event{Type: m.Tag, Message: m.Info}
			fmt.Printf("Event recieved: %s\n", m.Tag)
		}
		util.CheckError(err)

	}
}
