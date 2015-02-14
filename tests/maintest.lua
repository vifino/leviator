print("sup d00d, state 0 here")
id, msg = ipc_read(state_id)
print("Main Got from State ".. tostring(id)..": "..msg)
