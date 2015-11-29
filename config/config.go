package config

import (
	"strconv"
	"io/ioutil"
	"../util"
	"runtime"
	"os/user"
	"os"
	"fmt"
	"strings"
	"gopkg.in/yaml.v2"
)

type Conf struct { // This can be split into sections see https://code.google.com/p/gcfg/
	host string
	watch_dirs []string
	port, watch_freq, cores int	
}

var conf Conf

func init() {
	content, err := ioutil.ReadFile("./conf.ini")
	util.CheckError(err)
	err = yaml.Unmarshal(content, &conf)
	util.CheckError(err)
	
	if conf.cores < 1 {
		conf.cores = runtime.NumCPU()
	}
	fmt.Println(conf.watch_dirs)
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Warning error obtaining current username.", err)	
		// Because golang: user: Currently not implemented on linux/386 ... are you SERIOUS go!? 
		
		tmp := make([]string, len(conf.watch_dirs))
		for i, s := range conf.watch_dirs {
			tmp[i] = strings.Replace(s, "$user", usr.Username, -1)			
			if info, _ := os.Stat(tmp[i]); ! info.IsDir() {
				fmt.Printf("Watch dir %s not found/not a dir.", tmp[i])
			}
		}
		conf.watch_dirs = tmp	
	}
}

func GetCores() int {
	return conf.cores
}

func GetPort() int {
	return conf.port
}

func GetServerIP() string {
	return conf.host + ":" + strconv.Itoa(conf.port)
}

func GetWatchDirs() []string {
	return conf.watch_dirs //TODO: Set to current dir // Really? 
}

func GetWatchFrequency() int {
	return conf.watch_freq
}
