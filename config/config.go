package config

//TODO: Write a proper config, with this as the layer on top. 

func GetServerIP() string {
	return "127.0.0.1:4638"
}

func GetWatchDir() string {
	return "/home/tox/Downloads" //TODO: Set to current dir
}
func GetWatchFrequency() int {
	return 5
}
