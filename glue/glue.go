package glue

import (
	instance "../instance"
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	lua "github.com/vifino/golua/lua"
	luar "github.com/vifino/luar"
	"regexp"
	"time"
)

// Adds regmatch and regsub for easy regex usage.
const luacode string = `
regmatch = function(regex, input, num)
	if type(regex) == "string" then
		regex = regexp(regex)
	end
	return luar.slice2table(regex.FindAllString(tostring(input), tonumber(num) or -1))
end
regsub = function(regex, src, replacement)
	if type(regex) == "string" then
		regex = regexp(regex)
	end
	return regex.ReplaceAllString(tostring(src), tostring(replacement))
end
ln_read = function(prompt)
	if prompt == nil then
		return _LN_READ("")
	else
		return _LN_READ(tostring(prompt))
	end
end
log = function(...)
	println("["..time_fulldate().."]".." [State "..state_id..": "..state_filename.."]:".. " "..table.concat({...}, " "))
end
`

func BasicGlue(id int, state *lua.State) {
	luar.Register(state, "", luar.Map{
		"state_id":      id,
		"regexp":        regexp.Compile, // Regex
		"println":       fmt.Println,    // Println, just fmt.Println
		"sleep":         sleep,
		"_LN_READ":      linenoise.Line, // Line noise binding, for better repls and user input.
		"ln_addhistory": linenoise.AddHistory,
		"ln_clear":      linenoise.Clear,
		"time_time":     time_time,
		"time_date":     time_date,
		"time_fulldate": time_fulldate,
	})
	instance.Eval(state, luacode)
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
