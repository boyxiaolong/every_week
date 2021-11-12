
local idPlayer = 0;

local runtest = function ()
	setplayerid(idPlayer)

	-- print(idPlayer .. " begin sleep")

	sleep(3000)

	-- print(idPlayer .. " end sleep")

	logout()

	-- print(idPlayer .. " logout")

	_G["allroutine"][idPlayer] = nil;
end

idPlayer = getplayerid()
local thread = coroutine.create(runtest);

-- 存放到全局环境中，否则LUA将可能自动回收thread
_G["allroutine"][idPlayer] = thread
coroutine.resume(thread)

-- for n in pairs(_G["allroutine"]) do print(n) end
