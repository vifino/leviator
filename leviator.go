package main

import (
	instance "./instance"
	"strconv"
	"time"
)

var instancenum int = 10

// Lua functions
func main() {
	instances := instance.Init(instancenum)
	defer instance.Close(instances)
	for i := range instances {
		go instance.Eval(instances[i], `print("Hooray, I'm instance number `+strconv.Itoa(i)+`!")`)
	}
	time.Sleep(30 * time.Second)
}
