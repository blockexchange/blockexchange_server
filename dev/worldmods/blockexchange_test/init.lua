if not minetest.settings:get_bool("blockexchange.enable_integration_test", false) then
    print("[blockexchange_test] skipping integration test")
end

print("[blockexchange_test] executing integration test")

minetest.register_on_mods_loaded(function()
    minetest.request_shutdown("test done")
    -- TODO
end)