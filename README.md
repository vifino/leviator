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

`regexp(regex)` compiles regexes, it's regexp.Compile.
```lua
p = regexp '(\\w+)' -- Compiles a regexp.
for k,v in pairs(luar.slice2table(p.FindAllString('Hello there friend!',3))) do print(k,v) end
```
This should print:
```
1       Hello
2       there
3       friend
```

There is also a bit more...

# IPC:
This is where the fun stuff begins!

` ipc_send(id, message)` sends `message` to `id`, which is the state id explained above. Message must be a string.

`ipc_read(id)` reads the IPC queue produced by `ipc_send`. Returns the message.

`ipc_readNB(id)` is the same as above, except it doesn't block until a message is there, it just returns "" instead.
