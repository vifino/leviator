while true do
	io.write("> ")
	str = io.read()
	--print(str)
	ipc_send(1, state_id, "eval:"..str)
	id, res = ipc_read(state_id)
	suc, ret = res:match("^ret(.-):(.-)$")
	if suc == "True" then
		suc = true
	else
		suc = false
	end
	if suc == true then
		print("T: "..ret)
	else
		print("F: "..ret)
	end
end
