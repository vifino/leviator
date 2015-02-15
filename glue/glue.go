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
`

func BasicGlue(id int, state *lua.State) {
	luar.Register(state, "", luar.Map{
		"state_id":      id,
		"regexp":        regexp.Compile, // Regex
		"println":       fmt.Println,    // Println, just fmt.Println
		"sleep":         sleep,
		"_LN_READ":      linenoise.Line, // Line noise binding, for better repls and user input.
		"ln_addhistory": linenoise.AddHistory,
	})
	instance.Eval(state, luacode)
}

// Sleep
func sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}