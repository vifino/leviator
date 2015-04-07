package main

import (
	glue "./glue"
	instance "./instance"
	ipc "./ipc"
	scheduler "./scheduler"
	"fmt"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"os"
//	"strconv"
)

// Initialization n stuff.
//var instances []*lua.State
var instancenum int             // This will later be the amount of instances, which will be defined by number of files
var channels []chan ipc.IPCData // Communication chan's.

func state_register(id int, state *lua.State) {
	luar.Register(state, "", luar.Map{ // IPC
		"ipc_read":      ipc_read,
		"ipc_readNB":    ipc_readNB,
		"ipc_send":      ipc_send,
		"ipc_broadcast": ipc_broadcast,
	})
	glue.Glue(state, id)
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
	instancenum = 0
	if len(args) > 0 {
		go scheduler.RunScheduler()
		state := instance.Init_State()
		glue.Glue(state, 0)
		glue.Args(state, os.Args)
		instance.EvalFile(state, args[0])
	} else {
		fmt.Println("help.png")
	}
}
