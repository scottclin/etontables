package config

import (
	"strconv"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"../util"
	"runtime"
)

type Conf struct {
	host, watch_dir string
	port, watch_freq, cores int	
}

var conf Conf

func init() {
	content, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		util.CheckError(err)
	}
	err = yaml.Unmarshal(content, &conf)
	
	if conf.cores < 1 {
		conf.cores = runtime.NumCPU()
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

func GetWatchDir() string {
	return conf.watch_dir //TODO: Set to current dir
}

func GetWatchFrequency() int {
	return conf.watch_freq
}
