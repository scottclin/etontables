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
			switch mesg.Type {
			case "magnet":
				tor := loadTorrentMagnet(client, mesg.Message)
				<- tor.GotInfo()
			case "start", "kill":
				result := actionOnTorrent(client, mesg.Message, mesg.Type)
				if ! result {
					fmt.Fprint(os.Stderr, "Failed to find torrent %s", mesg.Message)
				}
			case "add_dir":
				AddWatchDir(mesg.Message)
			case "remove_dir":
				result := RemoveWatchDir(mesg.Message)
				if ! result {
					fmt.Fprint(os.Stderr, "Failed to remove directory %s", mesg.Message)
				}
				
			}
		case mesg2 := <- watchDirChannel:
			if mesg2.Type == "new_torrent_file"{
				tor := loadTorrentFile(client, mesg2.Message)
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
