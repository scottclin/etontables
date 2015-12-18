package console

import (
	"bufio"
	"fmt"
	"os"
	//	"time"
	"../../util"  // Is it wise to have relative imports? (GOPATH+= bootstrap?)
)

/**
This is the beginings of the clientside console based interface no flashy GUI shit here. 
Yay! Down with UI
*/

func ClientConsole(){
	//Get the message channel DO NOT READ MESSAGES FROM IT
	messagechannel := util.GetChannel("messagechannel")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Client console -> command:")

		usertag, err := reader.ReadString('\n')
		util.CheckError(err)

		fmt.Print("Client console -> info:")

		userinfo, err := reader.ReadString('\n')
		util.CheckError(err)

		//Create the message to send to the server using the message struct which is shared as it in the util
		m := util.Event{Type: "Message", Message: userinfo + usertag}

		//Send the message to the channel so the thread can pick it up and send it
		messagechannel <- m
	}

}
