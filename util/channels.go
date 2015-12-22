package util

var registeredChannels map[string] *chan Event

//Setup for keeping a register of channels
func SetupRegister(){
	registeredChannels = make(map[string] *chan Event)
}

//Adds the channel to the map returns false if one already exists with that name
func RegisterChannel(name string,channel *chan Event) bool{
	if _, ok := registeredChannels[name]; ok {
		return false
	}

	registeredChannels[name] = channel
	return true
}

func CheckForChannel(name string) (ok bool) {
	_, ok = registeredChannels[name] 
	return 
}

//For getting the channel by name returns nil if no channel found
func GetChannel(name string) *chan Event{
	channel, ok := registeredChannels[name]
	if ok {
		return channel
	}
	return nil
}

func RemoveChannel(name string) bool {
	if CheckForChannel(name) {
		delete(registeredChannels, name)
		return true
	}
	return false
}
