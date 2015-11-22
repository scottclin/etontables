package serverside

import (
	"github.com/anacrolix/torrent"
	"../util"
	"fmt"
	"os"
)

func Control(client torrent.Client){
	userEventChannel := util.GetChannel("userEventChannel")
	watchDirChannel := util.GetChannel("watchDirChannel")
	
	if userEventChannel == nil {
		fmt.Fprintf(os.Stderr, "Failed to get the userEventChannel, I think it should exist.\n")
	}
	if watchDirChannel == nil {
		fmt.Fprintf(os.Stderr, "Failed to get the watchDirChannel, I think it should exist.\n")
	}

	for {
		select {
		case mesg := <- userEventChannel:
			newUserAction := mesg.(util.Event)
			if newUserAction.Type == "magnet"{
				tor := loadTorrentMagnet(client, newUserAction.Message)
				<- tor.GotInfo()
				tor.DownloadAll()
			}
		case mesg2 := <- watchDirChannel:
			newFile := mesg2.(util.Event)
			if newFile.Type == "new_torrent_file"{
				tor := loadTorrentFile(client, newFile.Message)
				<- tor.GotInfo()
				tor.DownloadAll()
			}
		}
	}
}
