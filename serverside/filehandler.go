package serverside

import (
	"github.com/jackpal/bencode-go"
	"../util"
	"../config"
	"crypto/sha1"
	"bytes"
	"os"
	"regexp"
	"io/ioutil"
	"time"
	"strings"
	"fmt"
)

//Vars for things
var seen_torrent_files map[string]bool
var watch_dirs []string
var watch_duration int

/**
Initalise some settings
*/

func init(){
	seen_torrent_files = make(map[string]bool)
	watch_dirs = config.GetWatchDirs()
	watch_duration = config.GetWatchFrequency()
	fmt.Printf("Setup looking for files in %s\n", watch_dirs)
}

/**
Checks the folder for new files based which folders? are sent in the watchdir channel or in the config.
*/
func scanDirForTorrents(watch_dir string, watchdirchannel chan interface{}) {
	//Lets look in the dir and see if there is a new file
	files, err := ioutil.ReadDir(watch_dir)
	if err != nil {
		fmt.Println("Error scanning dir:", watch_dir)
		return
	}
	
	for _, fil := range files {
		if seen_torrent_files[fil.Name()] { //We have dealt with this torrent file before ignore
			continue
		}
		
		match, err := regexp.MatchString("torrent$", fil.Name()) // Is torrent file?
		util.CheckError(err)
		if match == false {
			continue
		}

		//Send out a message saying there is a new file and where it is
		fmt.Println("Sending the message to the channel")
		watchdirchannel <- util.Event{Type: "new_torrent_file", Message: watch_dir + "/" + fil.Name()}
		
		seen_torrent_files[fil.Name()] = true		
	}
}

func CheckForfile(){	
	watchdirchannel := util.GetChannel("watchdirchannel")
	if watchdirchannel == nil {
		watchdirchannel = make(chan interface{}, 5)	
		if !util.RegisterChannel("watchdirchannel", watchdirchannel) {
			fmt.Println("Failed to register channel.")
		}
	}

	for {
		for _, watchdir := range watch_dirs {
			scanDirForTorrents(watchdir, watchdirchannel)			
		}	
		time.Sleep(time.Duration(watch_duration) * time.Second)
	}
}


/**
Load up the newly found file that has been sent though the channel
*/
func LoadTorrentFile () {
	//Let get the channel if we cant then nothing we can do
	watchdirchannel := util.GetChannel("watchdirchannel")
	if watchdirchannel == nil {
		fmt.Printf("Failed to get channel dying")
		return
	}

	//Lets load up any file that comes though the channel
	for {
		newevent := <- watchdirchannel
		newfileevent := newevent.(util.Event)
		fmt.Printf("An event was recieved: %s %s\n", newfileevent.Type, newfileevent.Message)
		if newfileevent.Type != "new_torrent_file" {
			continue
		}

		//Open the file and do some shit with it
		file, err := os.Open(newfileevent.Message)
		defer file.Close()
		util.CheckError(err)

		rawinfo, _ := bencode.Decode(file)

		readableinfo, _ := rawinfo.(map[string]interface{})

		var b bytes.Buffer

		//Decode the file
		err = bencode.Marshal(&b, readableinfo)
		util.CheckError(err)

		hash := sha1.New()
		hash.Write(b.Bytes())

		var finalinfo util.MetaInfo

		err = bencode.Unmarshal(&b, &finalinfo.Info)
		util.CheckError(err)

		//This is all I know how to get out so far
		finalinfo.InfoHash = string(hash.Sum(nil))
		finalinfo.Announce = util.GetString(readableinfo, "announce")
		finalinfo.AnnounceList = util.GetSliceString(readableinfo, "announce-list")
		finalinfo.CreationDate = util.GetString(readableinfo, "creation date")
		finalinfo.Comment = util.GetString(readableinfo, "comment")
		finalinfo.CreatedBy = util.GetString(readableinfo, "created by")
		finalinfo.Encoding = strings.ToUpper(util.GetString(readableinfo, "encoding"))
		
		fmt.Println(finalinfo)
	}
}
