function post()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    -- Install extension
    if http.postValues.install_extension then
        local extension = json:unmarshalFile(string.format("extensions/%s/extension.json", http.postValues.install_extension))
        local successMessage = extension.id .. " has been installed."

        if extension.author == nil then
            extension.author = "Anonymous"
        end

        if extension.description == nil then
            extension.description = "-"
        end

        -- Run install script
        if file:exists(string.format("extensions/%s/install.lua", extension.id)) then
            local success = false
            try(
                function ()
                    -- Load file
                    dofile(string.format("extensions/%s/install.lua", extension.id))

                    success, message = install()
                    if success == false then
                        error(message or "install function explicitly returned false but did not provide any message.")
                        return
                    end

                    successMessage = string.format("Installed %s: %s", extension.id, message)
                end,
                -- Error function
                function (err)
                    session:setFlash("validationError", string.format("Failed to install %s: %s", extension.id, err))
                end
            )

            -- Abort if install.lua failed
            if not success then
                http:redirect("/subtopic/admin/extensions/install")
                return
            end
        end

        -- Install extension base
        db:execute("INSERT INTO castro_extensions (name, id, version, description, author, type, installed, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, b'1', ?, ?)", extension.name, extension.id, extension.version, extension.description, extension.author, extension.type, os.time(), os.time())

        -- Install Lua hooks
        if extension.hooks then
            for hook, script in pairs(extension.hooks) do
                db:execute("INSERT INTO castro_extension_hooks (extension_id, type, script, enabled) VALUES (?, ?, ?, b'1')", extension.id, hook, script)
            end
        end

        -- Install template hooks
        if extension.templateHooks then
            for hook, template in pairs(extension.templateHooks) do
                db:execute("INSERT INtO castro_extension_templatehooks (extension_id, type, template, enabled) VALUES (?, ?, ?, b'1')", extension.id, hook, template)
            end
        end

        -- Install pages
        if file:exists(string.format("extensions/%s/pages", extension.id)) then
            db:execute("INSERT INTO castro_extension_pages (extension_id, enabled) VALUES (?, b'1')", extension.id)
        end

        -- Install widgets
        if file:exists(string.format("extensions/%s/widgets", extension.id)) then
            db:execute("INSERT INTO castro_extension_widgets (extension_id, enabled) VALUES (?, b'1')", extension.id)
        end

        session:setFlash("success", successMessage)
        http:redirect("/subtopic/admin/extensions/install")
        return
    end

    -- Uninstall extension
    if http.postValues.uninstall_extension then
        local id = http.postValues.uninstall_extension

        -- Run uninstall.lua
        if file:exists(string.format("extensions/%s/uninstall.lua", id)) then
            try(
                function ()
                    -- Load file
                    dofile(string.format("extensions/%s/uninstall.lua", id))

                    local success, message = uninstall()
                    if success == false then
                        error(message or "uninstall function explicitly returned false but did not provide any message.")
                        return
                    end
                    session:setFlash("success", string.format("Uninstalled %s: %s", id, message or ""))
                end,
                -- Error function
                function (err)
                    session:setFlash("validationError", string.format("Failed to execute uninstall script for %s: %s. The extension have been removed anyway.", id, err))
                end
            )
        end

        -- Remove extension
        db:execute("DELETE FROM castro_extensions WHERE id = ?", id)

        http:redirect("/subtopic/admin/extensions/install")
        return
    end
end
