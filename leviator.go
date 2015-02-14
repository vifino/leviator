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
var instancenum int        // This will later be the amount of instances, which will be defined by number of files
var channels []chan string // Communication chan's.

func init_states(num int) []*lua.State {
	instances := instance.Init(num)
	channels = ipc.MakeChans(num)
	// Map functions.
	for i := range instances {
		luar.Register(instances[i], "", luar.Map{
			"state_id":      i,
			"regexp":        regexp.Compile, // Regex
			"println":       fmt.Println,    // Println, just fmt.Println
			"ipc_read":      ipc_read,
			"ipc_send":      ipc_send,
			"ipc_broadcast": ipc_broadcast,
			"sleep":         sleep,
		})
	}
	return instances
}

// IPC
func ipc_send(id int, msg string) {
	ipc.Send(channels[id], msg)
}

func ipc_read(id int) string {
	return ipc.Receive(channels[id])
}
func ipc_broadcast(msg string) {
	ipc.Broadcast(channels, msg)
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
