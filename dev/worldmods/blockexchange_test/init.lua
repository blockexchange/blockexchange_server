if not minetest.settings:get_bool("blockexchange.enable_integration_test", false) then
    print("[blockexchange_test] skipping integration test")
end

print("[blockexchange_test] executing integration test")

minetest.register_on_mods_loaded(function()
    minetest.after(1, function()
        blockexchange.api.get_token("Testuser", "default"):next(function()
            print("emerging")
            return blockexchange.emerge("singleplayer", {x=-50, y=-50, z=-50}, {x=50,y=50,z=50})
        end):next(function()
            minetest.request_shutdown("test done")
        end):catch(function(http_code)
            error("http: " .. http_code)
        end)
    end)
end)