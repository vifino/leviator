package glue

import (
	instance "../instance"
	scheduler "../scheduler"
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"regexp"
	"time"
)

import "C"

// Adds regmatch and regsub for easy regex usage.
const luacode string = `
-- Regex
regex = {}
regex.match = function(regex, input, num)
	if type(regex) == "string" then
		regex = regexp(regex)
	end
	return luar.slice2table(regex.FindAllString(tostring(input), tonumber(num) or -1))
end
regex.sub = function(regex, src, replacement)
	if type(regex) == "string" then
		regex = regexp(regex)
	end
	return regex.ReplaceAllString(tostring(src), tostring(replacement))
end
-- Loggin
log = function(...)
	println("["..time.fulldate().."]".." [State "..(state_id or "X")..": "..state.filename.."]:".. " "..table.concat({...}, " "))
end
-- State
state = {}
state.new = state_new
state_new = nil
state.eval = function(s, code)
	if type(st) == "userdata" and type(file) == "string" then
		return s.DoString(code)
	else
		error("State not userdata!")
	end
end
state.eval_async_func = state_asynceval
state_asynceval = nil
state.eval_async = function(s, code)
	if type(state) == "userdata" and type(file) == "string" then
		return state.eval_async_func(s,code)
	else
		error("State not userdata!")
	end
end
state.evalFile = function(s, file)
	if type(s) == "userdata" and type(file) == "string" then
		return s.DoFile(file)
	else
		error("State not userdata!")
	end
end
state.self = state_self
state_self = nil
-- LineNoise
ln = {}
ln.addhistory = ln_addhistory
ln_addhistory = nil
ln.clear = ln_clear
ln_clear = nil
ln.read_func = ln_read
ln_read = nil
ln.read = function(prompt)
	if prompt == nil then
		return ln.read_func("")
	else
		return ln.read_func(tostring(prompt))
	end
end
`

func BasicGlue(state *lua.State) {
	luar.Register(state, "", luar.Map{
		"regexp":            regexp.Compile, // Regex
		"println":           fmt.Println,    // Println, just fmt.Println
		"ln_read":           linenoise.Line, // Line noise binding, for better repls and user input.
		"ln_addhistory":     linenoise.AddHistory,
		"ln_clear":          linenoise.Clear,
		"state_self":        state,
		"state_asynceval":   state_async,
	})
	luar.Register(state, "time", luar.Map{
		"time":     time_time,
		"date":     time_date,
		"fulldate": time_fulldate,
		"sleep":    sleep,
	})
	//state.Register("state_loadstring", state_loadstring)
	state.Register("runbg", runbg)
	instance.Eval(state, luacode)
}
func Glue(state *lua.State,id int) {
	luar.Register(state, "", luar.Map{
		"state_id": id,
	})
	BasicGlue(state)
}
func Args(state *lua.State, args []string) {
	luar.Register(state, "", luar.Map{
		"args": args,
	})
	state.DoString(`
		-- Args
		args = luar.slice2table(args)
	`)
}

// Time.
const time_format string = "15:04:05 MST"
const date_format string = "_2.1.2006"
const fulldate_format string = "_2.1.2006 15:04:05"

func sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
func time_time() string {
	return time.Now().Format(time_format)
}
func time_date() string {
	return time.Now().Format(date_format)
}
func time_fulldate() string {
	return time.Now().Format(fulldate_format)
}

// State

func state_async_eval(L *lua.State, code string){ // Do. Not. Use.
	scheduler.Schedule(func(){
		L.DoString(code)
	})
}

/*func (L *lua.State) pcall(nargs, nresults, errfunc int) int {
	return int(C.lua_pcall(L.s, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

func state_loadstring(L *lua.State) int {
	string := L.CheckString(2)
	state := L.ToUserdata(1)
	if r := state.LoadString(string); r != 0 {
		L.PushString(state.ToString(-1))
		return 1
	} else {
		L.PushNil()
		return 1
	}
}

func state_pcall(L *lua.State) int {
	state := L.touserdata(1)
	L.CheckAny(1)

}*/

// Async stuff.
func runbg(L *lua.State) int {
	L.CheckAny(1)
	ind := L.GetTop() - 1
	go L.Pcall(ind, -1, 0)
	return 0
}
