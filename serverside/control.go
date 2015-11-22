package serverside

import (
	"github.com/anacrolix/torrent"
)

func Control(client *torrent.Client){
	userEventChannel := util.GetChannel("userEventChannel")
	if userEventChannel == nil {
		userEventChannel = make(chan interface{}, 5)
		if !util.RegisterChannel("userEventChannel", userEventChannel) {
			fmt.Fprintf("Failed to register a userEventChannel.\n")
		}
	}

}
