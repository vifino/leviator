while true do
	id, str = ipc_read(state_id)
	--print("Received: "..str)
	if id == 0 then
		act, code = str:match("^(.-):(.-)$")
		--print("code: "..code)
		if act == "eval" then
			--print("act is eval")
			func, err = loadstring(code)
			if func then
				--print(func())
				suc, res = pcall(func)
				--suc = true
				--res = func()
				--print(suc, res)
				ret = suc and "retTrue:" or "retFalse:"
				ret = ret..tostring(res)
				--print(ret)
				ipc_send(id, state_id, ret)
			
			else
				ipc_send(id, state_id, "retFalse:"..tostring(err))
			end
		end
	end
end
