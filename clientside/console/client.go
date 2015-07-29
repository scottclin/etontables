package console

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"../../util"
)

func ClientConsole(){

	messagechannel := util.GetChannel("messagechannel")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Client console -> command:")

		usertag, err := reader.ReadString('\n')
		util.CheckError(err)

		fmt.Print("Client console -> info:")

		userinfo, err := reader.ReadString('\n')
		util.CheckError(err)
		
		m := util.Message{Host: "local", Id: "local:" + time.Now().String(), Level: 0, Info: userinfo, Tag: usertag}

		messagechannel <- m
	}

}
