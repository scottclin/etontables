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
		*userEventChannel = make(chan util.Event, 5)
		if !util.RegisterChannel("userEventChannel", userEventChannel) {
			fmt.Println("Failed to create user event channel.")
		}
	}
	if watchDirChannel == nil {
		fmt.Fprintf(os.Stderr, "Failed to get the watchDirChannel, I think it should exist.\n")
	}

	for {
		select {
		case mesg := <- *userEventChannel:
			switch mesg.Type {
			case util.NewTorrentMagnet:
				tor := loadTorrentMagnet(client, mesg.Message)
				<- tor.GotInfo()
			case util.Start, util.Kill:
				result := actionOnTorrent(client, mesg.Message, mesg.Type)
				if ! result {
					fmt.Fprint(os.Stderr, "Failed to find torrent %s", mesg.Message)
				}
			case util.AddWatchDir :
				AddWatchDir(mesg.Message)
			case util.RemoveWatchDir :
				result := RemoveWatchDir(mesg.Message)
				if ! result {
					fmt.Fprint(os.Stderr, "Failed to remove directory %s", mesg.Message)
				}
				
			}
		case mesg2 := <- *watchDirChannel:
			if mesg2.Type == util.NewTorrentFile {
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
			case util.Start:
				torrent.DownloadAll()
				return true
			case util.Kill:
				torrent.Drop()
				return true
			}
		}
	}

	return false
}
