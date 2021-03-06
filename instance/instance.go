package instance

import (
	"fmt"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
)

func Init_State() *lua.State {
	state := luar.Init()
	luar.Register(state, "state", luar.Map{
		"new":     Init_State,
	})
	return state
}

func Eval(state *lua.State, code string) {
	err := state.DoString(code)
	if err != nil {
		errs := err.Error()
		fmt.Println(errs)
	}
}

func EvalFile(state *lua.State, filename string) {
	luar.Register(state, "", luar.Map{
		"state_filename": filename,
	})
	err := state.DoFile(filename)
	if err != nil {
		errs := err.Error()
		fmt.Println(errs)
	}
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
