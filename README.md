# leviator
Multithreaded multi-instanced LuaJIT Supervisor/Toolkit. Or something else.

TBD.

This thing runs lua files, each in a seperate instance, and using IPC you can comunicate between each instance.

# Running:
`leviator file1.lua file2.lua`

# Custom Variables:
This is easy, there is only `state_id`, which returns the ID of the state.
0 will always be the first file, 1 the second, etc...
You use this for IPC, explained further down.

# Custom Functions:
All the [luar](https://github.com/vifino/luar) function are still there, don't panic.

`println(value)` is `fmt.Println` from Go, simple.

`sleep(seconds)` is sleep. Does what you expect.

`time_time()` returns the time in format `15:04:05 MST`.

`time_date()` returns the date in format `_2.1.2006`, `_` being expanded if the day is below 10, 9 becomes 09, etc..

`time_fulldate` returns a combination between the two above: `_2.1.2006 15:04:05`

`log(message)` logs messages in a friendly format: `[02.1.2006 15:04:05] [State <id>: state_file.lua]: <message>`

## Regex:

`regexp(regex)` compiles regex expressions, it's regexp.Compile. You normally don't have to use this.

`regmatch(regex, string, maxMatches)` is like gmatch, just with regexes. `regex` is either a compiled regex expression or a normal one, which will get compiled then. maxMatches is obvious and optional.

`regsub(regex, source, replacement)` is a regex replace function. `regex` is the same as above, `source` is the string you want to replace something in and `replacement` is obvious.

Example:


```lua
p = regexp '(\\w+)' -- Compiles a regexp.
for k,v in pairs(regmatch(p, 'Hello there friend!',3)) do print(k,v) end
```
This should print:
```
1       Hello
2       there
3       friend
```

But wait! There is more!

## IPC:
This is where the fun stuff begins!

` ipc_send(id, from, message)` sends `message` to `id`, which is the state id explained above. Message must be a string.

`ipc_read(id)` reads the IPC queue produced by `ipc_send`. Returns the id and message.

`ipc_readNB(id)` is the same as above, except it doesn't block until a message is there, it just returns 0 and "" instead.

## LineNoise:
Ever wanted to do a fancy repl for your project? Well, now you can.

`ln_read(prompt)` reads a line from stdin and returns a string. `prompt` is the prompt.

`ln_addhistory(string)` adds `string` to the prompt history, accessable using the arrow keys.

`ln_clear()` clears the screen.

# License:
MIT
