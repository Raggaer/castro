local list = {
    "onStartup",
    "onNewArticle",
    "onEditArticle"
}

function listExtensionHooks()
    return list
end

-- Execute a hook by type
function executeHook(hookType, ...)
    local hooks = db:query("SELECT * FROM castro_extension_hooks WHERE type = ? AND enabled = 1", hookType, (app.Mode ~= "dev")) -- no cache on dev mode
    if hooks then
        for _, hook in pairs(hooks) do
            local script, err = loadfile(string.format("extensions/%s/%s", hook.extension_id, hook.script))
            if script then
                script()
                if _G[hookType] then
                    _G[hookType](...)
                    -- Remove function from global table
                    _G[hookType] = nil
                else
                    print(string.format("Warning: extensions/%s/%s\nFunction %s is missing.", hook.extension_id, hook.script, hookType))
                    log:error(string.format("Warning: extensions/%s/%s\nFunction %s is missing.", hook.extension_id, hook.script, hookType))
                end
            else
                print(string.format("Failed to run plugin %s, file: extensions/%s/%s\n%s", hook.extension_id, hook.extension_id, hook.script, err))
                log:error(string.format("Failed to run plugin %s, file: extensions/%s/%s\n%s", hook.extension_id, hook.extension_id, hook.script, err))
            end
        end
    end
end
