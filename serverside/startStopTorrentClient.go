package serverside

import (
	"../util"
	"github.com/anacrolix/torrent"
)

var client torrent.Client

func Start() *torrent.Client {
	client, err := torrent.NewClient()
	util.CheckError(err)
	return &client
}

func Stop() {
	client.Close()
}
