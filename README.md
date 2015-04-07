# leviator
Multithreaded multi-instanced LuaJIT Supervisor/Toolkit. Or something else.

TBD.

This thing runs lua files, each in a seperate instance, and using IPC you can comunicate between each instance.

# Running:
`leviator file.lua args`

# Custom Variables:
This is easy, there is only `state_file`. 

`state_file`, is the file the state is currently running.

# Custom Functions:
All the [luar](https://github.com/vifino/luar) function are still there, don't panic.

`println(value)` is `fmt.Println` from Go, simple.

`time.sleep(seconds)` is sleep. Does what you expect.

`time.time()` returns the time in format `15:04:05 MST`.

`time.date()` returns the date in format `_2.1.2006`, `_` being expanded if the day is below 10, 9 becomes 09, etc..

`time.fulldate` returns a combination between the two above: `_2.1.2006 15:04:05`

`log(message)` logs messages in a friendly format: `[02.1.2006 15:04:05] [State <id>: state_file.lua]: <message>`

## Regex:

`regexp(regex)` compiles regex expressions, it's regexp.Compile. You normally don't have to use this.

`regex.match(regex, string, maxMatches)` is like gmatch, just with regexes. `regex` is either a compiled regex expression or a normal one, which will get compiled then. maxMatches is obvious and optional.

`regex.sub(regex, source, replacement)` is a regex replace function. `regex` is the same as above, `source` is the string you want to replace something in and `replacement` is obvious.

Example:


```lua
p = regexp '(\\w+)' -- Compiles a regexp.
for k,v in pairs(regex.match(p, 'Hello there friend!',3)) do print(k,v) end
```
This should print:
```
1       Hello
2       there
3       friend
```

But wait! There is more!

## LineNoise:
Ever wanted to do a fancy repl for your project? Well, now you can.

`ln.read(prompt)` reads a line from stdin and returns a string. `prompt` is the prompt.

`ln.addhistory(string)` adds `string` to the prompt history, accessable using the arrow keys.

`ln.clear()` clears the screen.

So yeah.. You can make pretty repl's now.
There is more, though :)

# States
States are basically lua vm's that you can spin up and do stuff with.

`state.new()` returns a new state. The state returned by this will be called `s` below.

`state.self` returns the state that is currently running this. It's dangerous, you can get the state to crash if you aren't careful.

`state.eval(s, code)` eval's code in a state.

`state.evalFile(s, filepath)` is the same as above, but instead of supplying the code in the string, you supply a filepath.

`state.eval_async(s, code)` behaves like `state.eval(s, code)`, but runs the code in the background.


# License:
MIT
