package main

import (
	"github.com/jackpal/bencode-go"
	"net"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"io/ioutil"
	"./util"
	"./config"
	"time"
	"strings"
	"runtime"
	"crypto/sha1"
	"bytes"
)

var seen_torrent_files map[string]bool
var watch_dir string
var watch_duration int

func main(){
	runtime.GOMAXPROCS(2)
	util.SetupRegister()
	seen_torrent_files = make(map[string]bool)
	watch_dir = config.GetWatchDir()
	watch_duration = config.GetWatchFrequency()
	
	go checkForfile()
	go loadTorrentFile()
	fmt.Printf("Setup looking for files in %s\n", watch_dir)
	
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":4638")
	util.CheckError(err)
	
	ln, err := net.ListenTCP("tcp", tcpAddr)

	util.CheckError(err)
	
	for {
		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go handleClient(conn)
	}
}

func checkForfile(){
	watchdirchannel := make(chan util.Event, 5)
	result := false
	if watchdirchannel != nil {
		result = util.RegisterChannel("watchdirchannel", watchdirchannel)
	}
	if ! result {
		fmt.Println("Failed to register channel checking if one already exists")
		watchdirchannel = util.GetChannel("watchdirchannel")
		if watchdirchannel == nil {
			fmt.Println("Well that did not work, I will just do my thing without telling anything")
		}else{
			fmt.Println("Success one was already here")
		}
	}
	for {
		_, err := os.Stat(watch_dir)
		if os.IsNotExist(err) {
			fmt.Printf("The path %s does not exist\n", watch_dir)
		}else if err != nil{
			fmt.Printf("Something went wrong trying to stat the path: %s\n", watch_dir)
		}

		files, err := ioutil.ReadDir(watch_dir)
		for _, fil := range files {
			//See if this is a torrnet file
			match, err  := regexp.MatchString("torrent$", fil.Name())

			util.CheckError(err)
			
			if seen_torrent_files[fil.Name()] { //We have delt with this torrnet file before ignore
				continue
			}else if match == false { //Ignore if not a torrent file
				continue
			}

			fmt.Println("Sending the message to the channel")
			watchdirchannel <- util.Event{Type: "new_torrent_file", Message: watch_dir + "/" + fil.Name()}
			
			seen_torrent_files[fil.Name()] = true		
		}
		time.Sleep(time.Duration(watch_duration) * time.Second)
	}
}

func loadTorrentFile () {
	watchdirchannel := util.GetChannel("watchdirchannel")
	if watchdirchannel == nil {
		fmt.Printf("Failed to get channel dying")
		return
	}
	for {
		newfileevent := <- watchdirchannel
		fmt.Printf("An event was recieved: %s %s\n", newfileevent.Type, newfileevent.Message)
		if newfileevent.Type != "new_torrent_file" {
			continue
		}

		_ , err := os.Stat(newfileevent.Message)
		if err != nil{
			fmt.Println("I cannot stat the file")
		}
		
		file, err := os.Open(newfileevent.Message)
		defer file.Close()
		
		util.CheckError(err)

		rawinfo, _ := bencode.Decode(file)

		readableinfo, _ := rawinfo.(map[string]interface{})

		var b bytes.Buffer

		err = bencode.Marshal(&b, readableinfo)
		util.CheckError(err)

		hash := sha1.New()
		hash.Write(b.Bytes())

		var finalinfo util.MetaInfo

		err = bencode.Unmarshal(&b, &finalinfo.Info)
		util.CheckError(err)

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

func handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()

	dec := json.NewDecoder(conn)
	for {
		var m util.Message
		if err := dec.Decode(&m); err != nil {
			break
		}

		fmt.Printf("%s, %d, %s: %s %s\n", m.Host, m.Level, m.Id, m.Info, m.Tag)

		_, err := conn.Write([]byte("Message received"))

		switch m.Tag {
		case "set_watch_folder": watch_dir = m.Info
			fmt.Printf("Set the watch dir to be %s\n", watch_dir)
		default: fmt.Printf("Unreconised tag sent: %s\n", m.Tag)
		}
		util.CheckError(err)
				
	}
}
