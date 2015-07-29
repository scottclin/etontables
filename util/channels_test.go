package util

import (
	"testing"
	"os"
)

func TestRegisterChannel(t *testing.T){

	test_chan := make(chan Event)
	result := RegisterChannel("test", test_chan)
	if ! result {
		t.Fail()
	}

	result = RegisterChannel("test", test_chan)
	if result {
		t.Fail()
	}

	RemoveChannel("test")
}

func TestGetChannel(t *testing.T){

	test_chan := make(chan Event)
	result := RegisterChannel("test", test_chan)
	if ! result {
		t.FailNow()
	}

	retrive_test_chan := GetChannel("test")
	if retrive_test_chan == nil {
		t.Fail()
	}

	ask_for_random := GetChannel("I like glasses")

	if ask_for_random != nil {
		t.Fail()
	}

	RemoveChannel("test")
}

func TestRemoveChannel(t * testing.T){
	
	test_chan := make(chan Event)
	result := RegisterChannel("test", test_chan)
	if ! result {
		t.FailNow()
	}
	
	RemoveChannel("test")

	retrived_chan := GetChannel("test")
	if retrived_chan != nil {
		t.Fail()
	}

}

func TestMain(m *testing.M){
	SetupRegister()
	os.Exit(m.Run())
}
