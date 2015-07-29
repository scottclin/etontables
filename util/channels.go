package util

var registeredChannels map[string] chan interface{}

//Setup for keeping a register of channels
func SetupRegister(){
	registeredChannels = make(map[string] chan interface{})
}

//Adds the channel to the map returns false if one already exists with that name
func RegisterChannel(name string,channel chan interface{}) bool{
	if _, ok := registeredChannels[name]; ok {
		return false
	}

	registeredChannels[name] = channel
	return true
}

func CheckForChannel(name string)bool {
	_, ok := registeredChannels[name] 
	return ok
}

//For getting the channel by name returns nil if no channel found
func GetChannel(name string) chan interface{}{
	channel, ok := registeredChannels[name]
	
	if ok {
		return channel
	}

	return nil
}

func RemoveChannel(name string){
	delete(registeredChannels, name)
}