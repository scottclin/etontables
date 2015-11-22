package serverside

import (
	"../util"
	"github.com/anacrolix/torrent"
)

/**
Load a magnet and send to be watched and modified
*/
func loadTorrentMagnet(client torrent.Client, magnet string) torrent.Torrent {
	torrentMagnet, err := client.AddMagnet(magnet)

	util.CheckError(err)

	return torrentMagnet
}
