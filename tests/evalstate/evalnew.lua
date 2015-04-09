while true do
	local code = ln.read("> ")
	ln.addhistory(code)
	local func, err = loadstring(code)
	if func then
		local suc, res = pcall(func)
		--suc = true
		--res = func()
		--print(suc, res)
		--print(ret)
		if suc then
			print(res)
		else
			print("Error :".. tostring(res))
		end
	else
		print("Error: "..tostring(err))
	end
end
