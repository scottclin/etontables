package serverside

import (
	"net"
	"encoding/json"
	"fmt"
	"../util"
)


/**
This handles communications with the clients, recieves messages and sends them to the right channels 
*/
func HandleClient(conn net.Conn) {
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
