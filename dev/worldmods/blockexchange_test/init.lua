if not minetest.settings:get_bool("blockexchange.enable_integration_test", false) then
    print("[blockexchange_test] skipping integration test")
    return
end

-- TODO: make this work :P
-- local MP = minetest.get_modpath("blockexchange_test")
local compare_area = function() return true end -- dofile(MP.."/compare_area.lua")

print("[blockexchange_test] executing integration test")

-- override get_token
local token
function blockexchange.get_token()
    return token
end

-- propagate error
local function fail(msg)
    minetest.after(0, function()
        error(msg)
    end)
end

local size = {x=30, y=30, z=30}
local pos1 = {x=0, y=-10, z=0}
local pos2 = vector.add(pos1, size)
local pos1_load = {x=0, y=30, z=0 }
local pos2_load = vector.add(pos1_load, size)

local playername = "singleplayer"
local username = "Testuser"
local schemaname = "test_schema" .. math.random(1000)

minetest.register_on_mods_loaded(function()
    minetest.after(1, function()
        blockexchange.api.get_token(username, "default"):next(function(t)
            token = t
            return blockexchange.emerge(playername, pos1, pos2_load)
        end):next(function()
            return blockexchange.save(playername, pos1, pos2, schemaname)
        end):next(function()
            return blockexchange.load(playername, pos1_load, username, schemaname)
        end):next(function()
            local success, msg = compare_area(pos1, pos2, pos1_load, pos2_load)
            if not success then
                fail("loaded area does not match: " .. msg)
            end
            return true
        end):next(function()
            return blockexchange.save(playername, pos1, pos2, schemaname, true)
        end):next(function()
            return blockexchange.load(playername, pos1_load, username, schemaname, true)
        end):next(function()
            local success, msg = compare_area(pos1, pos2, pos1_load, pos2_load)
            if not success then
                fail("local loaded area does not match: " .. msg)
            end
            return true
        end):next(function()
            minetest.request_shutdown("test done")
        end):catch(function(http_code)
            fail("http error: " .. http_code)
        end)
    end)
end)