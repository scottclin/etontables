package serverside

import (
	"../config"
	"../util"
	"fmt"
	"github.com/anacrolix/torrent"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

//Vars for things
var seenTorrentFiles map[string]bool
var watchDirs []string
var watchDuration int
var loadFileChannel chan interface{}

/**
Initalise some settings
*/

func init() {

	seenTorrentFiles = make(map[string]bool)
	watchDirs = config.GetWatchDirs()
	watchDuration = config.GetWatchFrequency()
	fmt.Printf("Setup looking for files in %s\n", watchDirs)
}

/**
Checks the folder for new files based which folders? are sent in the watchdir channel or in the config.
*/
func scanDirForTorrents(watchDir string, watchDirChannel chan interface{}) {
	//Lets look in the dir and see if there is a new file
	files, err := ioutil.ReadDir(watchDir)
	if err != nil {
		fmt.Println("Error scanning dir:", watchDir)
		return
	}

	for _, fil := range files {
		if seenTorrentFiles[fil.Name()] { //We have dealt with this torrent file before ignore
			continue
		}

		match, err := regexp.MatchString("torrent$", fil.Name()) // Is torrent file?
		util.CheckError(err)
		if match == false {
			continue
		}

		//Send out a message saying there is a new file and where it is
		fmt.Println("Sending the message to the channel")
		watchDirChannel <- util.Event{Type: "new_torrent_file", Message: watchDir + "/" + fil.Name()}

		seenTorrentFiles[fil.Name()] = true
	}
}

/**
Checks to see if a new file has appeared, keeps track of the files we know about
*/
func CheckForfile() {
	watchDirChannel := util.GetChannel("watchDirChannel")
	if watchDirChannel == nil {
		watchDirChannel = make(chan interface{}, 5)
		if !util.RegisterChannel("watchDirChannel", watchDirChannel) {
			fmt.Println("Failed to register channel.")
		}
	}

	for {
		for _, watchdir := range watchDirs {
			scanDirForTorrents(watchdir, watchDirChannel)
		}
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

/**
Load up the newly found file that has been sent though the channel
Sends for watching and mofication
*/
func loadTorrentFile(client torrent.Client, torrentFileString string) torrent.Torrent{

	//Lets load up any file that comes though the channel

	if _, err := os.Stat(torrentFileString); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The file handler was unable to stat the file: %s\n", torrentFileString)
	}

	torrentFile, err := client.AddTorrentFromFile(torrentFileString)

	util.CheckError(err)

	return torrentFile
}
