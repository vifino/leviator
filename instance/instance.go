package instance

import (
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
)

func Init_State() *lua.State {
	return luar.Init()
}

func Eval(state *lua.State, code string) {
	state.DoString(code)
}

func Init(instances int) []*lua.State {
	retval := make([]*lua.State, instances)
	for i := range retval {
		retval[i] = Init_State()
	}
	return retval
}
func Close(instances []*lua.State) {
	for i := range instances {
		instances[i].Close()
	}
}
