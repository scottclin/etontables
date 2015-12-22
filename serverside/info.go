package serverside

import (
	"../util"
	//	"../config"
	"fmt"
)

/**
This is for setting up and starting all things information related like the infoChannel.
Then the info control is kicked off which handles updating information as well as reporting. This should mean there is consistancy in the information reported. 
*/
func InfoStart() {
	var infoChannel chan util.Event
	infoChannel = make(chan util.Event, 50)
	if !util.RegisterChannel("infoChannel", &infoChannel) {
		fmt.Println("Failed to create info channel.")
	}	
}
