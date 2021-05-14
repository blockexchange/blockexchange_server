if not minetest.settings:get_bool("blockexchange.enable_integration_test", false) then
    print("[blockexchange_test] skipping integration test")
end

print("[blockexchange_test] executing integration test")

-- override get_token
local token
function blockexchange.get_token()
    return token
end

local pos1 = {x=0, y=0, z=0}
local pos2 = {x=20,y=20,z=20}
local playername = "singleplayer"
local username = "Testuser"
local schemaname = "test_schema" .. math.random(1000)

minetest.register_on_mods_loaded(function()
    minetest.after(1, function()
        blockexchange.api.get_token(username, "default"):next(function(t)
            token = t
            return blockexchange.emerge(playername, pos1, pos2)
        end):next(function()
            return blockexchange.save(playername, pos1, pos2, schemaname)
        end):next(function()
            return blockexchange.load(playername, pos1, username, schemaname)
        end):next(function()
            minetest.request_shutdown("test done")
        end):catch(function(http_code)
            minetest.after(0, function()
                error("http: " .. http_code)
            end)
        end)
    end)
end)