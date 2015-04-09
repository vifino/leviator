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
state.eval = function(s, code, ...)
	if type(s) == "userdata" then
		local code = code
		if s ~= state.self then
			if type(code) == "function" then
				code = string.dump(code)
			end
			return state.dostring(s, code, ...)
		else
			if type(code) == "string" then
				local ret = s.DoString(code)
				if ret == nil then
					return true
				else
					return false, ret
				end
			else
				return false, "state.self does not support functions as code."
			end
		end
	else
		error("State not userdata!")
	end
end
state.eval_async = function(s, code)
	if type(s) == "userdata" then
		return state.eval_async_func(s, code)
	else
		error("State not userdata!")
	end
end
state.evalFile = function(s, file)
	if type(s) == "userdata" then
		return s.DoFile(file)
	else
		error("State not userdata!")
	end
end
state.loadbytecode = function(s, code, name)
	if type(s) == "userdata" then
		local res = s.LoadBuffer(code, #code, name or "str")
		if res~=0 then
				local err=s.ToString(-1)
				s.SetTop(-2)
				return false,err
			else
				return true
			end
	else
		error("State not userdata!")
	end
end
state.push=function(L, ...)
	if type(L) == "userdata" then
		local p={...}
		for i=1,select("#",...) do
			local v=p[i]
			local tpe=type(v)
			if tpe=="string" then
				L.PushString(v)
			elseif tpe=="number" then
				L.PushNumber(v)
			elseif tpe=="nil" then
				L.PushNil()
			elseif tpe=="boolean" then
				L.pushBoolean(v)
			elseif tpe=="function" then
				state.loadbytecode(L, string.dump(v))
			elseif tpe=="table" then
				error("Tables not implemented.")
			else
				error("Unsupported type: "..tpe)
			end
		end
	else
		error("State not userdata!")
	end
end
state.type = function(L, n)
	if type(L) == "userdata" then
		return go.LuaValType2int(L.Type(n))
	else
		error("State not userdata!")
	end
end
state.pop = function(L, n)
	if type(L) == "userdata" then
		n=n or 1
		local o={}
		for i=n,1,-1 do
			local tpe=state.type(L, -1)
			if tpe==0 then
				o[i]=nil
			elseif tpe==1 then
				o[i]=L.ToBoolean(-1)
			elseif tpe==2 or tpe==7 then
				o[i]=L.ToUserdata(-1)
			elseif tpe==3 then
				o[i]=L.ToNumber(-1)
			elseif tpe==4 then
				o[i]=L.ToString(-1)
			elseif tpe==5 then
				o[i]={}
			elseif tpe==6 then
				--local ou=""
				--local writer=ffi.cast("lua_Writer",function(L,p,sz,ud)
				--	ou=ou..ffi.string(p,sz)
				--	return 0
				--end)
				--LC.dump(writer,nil)
				--o[i]=assert(loadstring(ou))
			elseif tpe==8 then
				o[i]=L.ToThread(-1)
			end
			L.SetTop(-2)
		end
		return unpack(o,1,n)
	else
		error("State not userdata!")
	end
end
state.pcall = function(L, ...)
	if type(L) == "userdata" then
		local t=L.GetTop()
		state.push(L,...)
		local res=L.Pcall(select("#",...),-1,0)
		if res==0 then
			return true,state.pop(L, (L.GetTop()-t)+1)
		else
			local err=L.ToString(-1)
			L.SetTop(-2)
			return false,err
		end
	else
		error("State not userdata!")
	end
end
state.dostring = function(L, txt,...)
	if type(L) == "userdata" then
		local func,err=state.loadbytecode(L, txt)
		if not func then
			return false,err
		end
		return state.pcall(L, ...)
	else
		error("State not userdata!")
	end
end
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
		"regexp":        regexp.Compile, // Regex
		"println":       fmt.Println,    // Println, just fmt.Println
		"ln_read":       linenoise.Line, // Line noise binding, for better repls and user input.
		"ln_addhistory": linenoise.AddHistory,
		"ln_clear":      linenoise.Clear,
	})
	luar.Register(state, "state", luar.Map{
		"self":            state,
		"eval_async_func": state_async_eval,
		"byname_func":     state_byname,
		"loadstring_func": state_loadstring,
	})
	luar.Register(state, "time", luar.Map{
		"time":     time_time,
		"date":     time_date,
		"fulldate": time_fulldate,
		"sleep":    sleep,
	})
	luar.Register(state, "go", luar.Map{
		"LuaValType2int": LuaValType2int,
	})
	//state.Register("state_loadstring", state_loadstring)
	state.Register("runbg", runbg)
	instance.Eval(state, luacode)
}
func Glue(state *lua.State, id int) {
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

func state_async_eval(L *lua.State, code string) { // Do. Not. Use.
	scheduler.Schedule(func() {
		L.DoString(code)
	})
}

func state_byname(L *lua.State, name string) *luar.LuaObject {
	return luar.NewLuaObjectFromName(L, name)
}

func state_loadstring(L *lua.State, code string) error {
	if err := L.DoString(code); err != nil {
		return err
	} else {
		return nil
	}
}

// Go helper functions.

func LuaValType2int(V lua.LuaValType) int {
	return int(V)
}

// Async stuff.
func runbg(L *lua.State) int {
	L.CheckAny(1)
	ind := L.GetTop() - 1
	go L.Pcall(ind, -1, 0)
	return 0
}
