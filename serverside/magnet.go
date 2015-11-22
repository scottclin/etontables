package serverside

import (
	"../util"
	"github.com/anacrolix/torrent"
)

/**
Load a magnet and send to be watched and modified
*/
func LoadTorrentMagnet(client *torrent.Client, magnet string) {
	torrentMagnet, err := client.AddMagnet(magnet)

	util.CheckError(err)
}
