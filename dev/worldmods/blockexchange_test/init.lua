if not minetest.settings:get_bool("blockexchange.enable_integration_test", false) then
    print("[blockexchange_test] skipping integration test")
    return
end

local MP = minetest.get_modpath("blockexchange_test")
local compare_area = dofile(MP.."/compare_area.lua")

print("[blockexchange_test] executing integration test")

local size = {x=30, y=30, z=30}
local pos1 = {x=0, y=-10, z=0}
local pos2 = vector.add(pos1, size)
local pos1_load = {x=0, y=30, z=0 }
local pos2_load = vector.add(pos1_load, size)

local playername = "singleplayer"
local username = "Testuser"
local schemaname = "test_schema" .. math.random(1000)

Promise.async(function(await)
    await(Promise.mods_loaded())
    await(Promise.after(1))
    local token = await(blockexchange.api.get_token(username, "default"))

    local player_settings = blockexchange.get_player_settings(playername)
    player_settings.token = token
    blockexchange.set_player_settings(playername, player_settings)

    await(blockexchange.emerge(playername, pos1, pos2_load))
    await(blockexchange.save(playername, pos1, pos2, schemaname))
    await(blockexchange.load(playername, pos1_load, username, schemaname))
    local success, msg = compare_area(pos1, pos2, pos1_load, pos2_load, {}) -- TODO: fix/enable full checks
    if not success then
        error("loaded area does not match: " .. msg, 0)
    end

    await(blockexchange.save_local(playername, pos1, pos2, schemaname))
    await(blockexchange.load_local(playername, pos1_load, schemaname))
    success, msg = compare_area(pos1, pos2, pos1_load, pos2_load)
    if not success then
        error("local loaded area does not match: " .. msg, 0)
    end

    minetest.request_shutdown("test done")
end):catch(function(e)
    minetest.after(0, function()
        error(e)
    end)
end)
