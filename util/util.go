package util

/*
The package util is the package for utility functions and for types found thoughout the application, or need to be known application wide.

Author: Clinton Scott brutii@gmail.com
*/

import (
	"os"
	"fmt"
	"net"
	"encoding/json"
)
//The type expected by the server to be sent to it. 
type Message struct{
	Host, Info, Id, Tag string
	Level int
}
//The type expected to sent and recieved in channels
type Event struct{
	Type, Message string
}
//Helper to check and throw an error if required
func CheckError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

//Get the string I am looking for 
func GetString(m map[string]interface{}, k string) string{
	v,ok := m[k]
	if ok {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}

//Get the slice of strings I am looking for
func GetSliceString(m map[string]interface{}, k string) (resultSlice [][]string){
	v,ok := m[k]
	if ok {
		w, _ := v.([]interface{})
		for _, x := range w {
			var stringSlice []string
			y := x.([]interface{})
			for _, z := range y {
				a, ok := z.(string)
				if ok {
					stringSlice = append(stringSlice, a)
				}				
			}
			resultSlice = append(resultSlice, stringSlice)
		}		
	}
	return
}

/*
Will flesh out later into something more useful later or write another one to be more useful.
*/

func SendMessage (connection net.Conn){
	var messagechannel chan interface{}
	
	if CheckForChannel("messagechannel") {
		messagechannel = make(chan interface{}, 10)
		RegisterChannel("messagechannel", messagechannel)
	}else{
		messagechannel =  GetChannel("messagechannel")
	}
	
	for {
		m := <- messagechannel
		enc := json.NewEncoder(connection)	
		enc.Encode(m)
	}
}

func Exists(filepath string) bool {
    _, err := os.Stat(filepath)
    if err != nil {
		if  os.IsNotExist(err) {
			return false
		} else {
			fmt.Println("Something went wrong with dir:", filepath)
    		return false
		}
    }
	return true
}

func IsDir(filepath string) bool { //Sort of like, "Exists && IsDir"
    src, err := os.Stat(filepath)
    if err != nil {
		if  os.IsNotExist(err) {
			return false
		} else {
			fmt.Println("Something went wrong with dir:", filepath)
    		return false
		}
    }
	return src.IsDir()		
}
