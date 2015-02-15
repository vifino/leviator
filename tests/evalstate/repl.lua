while true do
	str = ln_read("> ")
	ipc_send(1, state_id, "eval:"..str)
	ln_addhistory(str)
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
