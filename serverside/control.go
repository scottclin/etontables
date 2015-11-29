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
			switch newUserAction.Type {
			case "magnet":
				tor := loadTorrentMagnet(client, newUserAction.Message)
				<- tor.GotInfo()
			case "start", "kill":
				result := actionOnTorrent(client, newUserAction.Message, newUserAction.Type)
				if ! result {
					fmt.Fprint(os.Stderr, "Failed to find torrent %s", newUserAction.Message)
				}				
			}
		case mesg2 := <- watchDirChannel:
			newFile := mesg2.(util.Event)
			if newFile.Type == "new_torrent_file"{
				tor := loadTorrentFile(client, newFile.Message)
				<- tor.GotInfo()
			}
		}
	}
}

func actionOnTorrent(client torrent.Client,torrentName string, action string)bool{
	
	allTorrents := client.Torrents()
	for _, torrent := range allTorrents{
		if torrent.Name() == torrentName {
			switch action {
			case "start":
				torrent.DownloadAll()
				return true
			case "kill":
				torrent.Drop()
				return true
			}
		}
	}

	return false
}
