package main

import (
	instance "./instance"
	ipc "./ipc"
	"fmt"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Initialization n stuff.
//var instances []*lua.State
var instancenum int             // This will later be the amount of instances, which will be defined by number of files
var channels []chan ipc.IPCData // Communication chan's.

func init_state() *lua.State { // Not normally used, but can be used to add new states at runtime.
	state := instance.Init_State()
	state_register(instancenum+1, state)
	return state
}

func state_register(id int, state *lua.State) {
	luar.Register(state, "", luar.Map{
		"state_id":      id,
		"regexp":        regexp.Compile, // Regex
		"println":       fmt.Println,    // Println, just fmt.Println
		"ipc_read":      ipc_read,
		"ipc_readNB":    ipc_readNB,
		"ipc_send":      ipc_send,
		"ipc_broadcast": ipc_broadcast,
		"sleep":         sleep,
	})
}

func init_states(num int) []*lua.State {
	instances := instance.Init(num)
	channels = ipc.MakeChans(num)
	// Map functions.
	for i := range instances {
		state_register(i, instances[i])
	}
	return instances
}

// IPC
func ipc_send(to int, from int, msg string) {
	ipc.Send(channels[to], from, msg)
}

func ipc_read(id int) (int, string) {
	return ipc.Receive(channels[id])
}

func ipc_readNB(id int) (int, string) {
	return ipc.ReceiveNonBlocking(channels[id])
}
func ipc_broadcast(from int, msg string) {
	ipc.Broadcast(channels, from, msg)
}

// Sleep
func sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Main thing.
/*
func main() {
	var exitchan chan bool
	for i := range instances {
		if i == 0 {
			go func() {
				instance.Eval(instances[0], `print("Hooray, I'm instance number `+strconv.Itoa(i)+`!")`)
				exitchan <- true
			}()
		} else {
			go instance.Eval(instances[i], `print("Hooray, I'm instance number `+strconv.Itoa(i)+`!")`)
		}
	}

}*/
func main() {
	args := os.Args[1:]
	instancenum = len(args)
	if len(args) > 0 {
		c := make(chan bool)
		instances := init_states(len(args))
		for i := range args {
			if i == 0 {
				fmt.Println("State 0 executing File: " + args[0])
				go func() {
					instance.EvalFile(instances[0], args[0])
					c <- true
				}()
			} else {
				fmt.Println("State " + strconv.Itoa(i) + " executing File: " + args[i])
				go instance.EvalFile(instances[i], args[i])
			}
		}
		//instance.EvalFile(instances[0], args[0])
		<-c
	}
}
