require "util"

function post()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    if http.postValues.install_extension then
        local extension = json:unmarshalFile(string.format("extensions/%s/extension.json", http.postValues.install_extension))
        db:execute("INSERT INTO castro_extensions (name, id, version, description, author, type, installed, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, b'1', ?, ?)", extension.name, extension.id, extension.version, extension.description, extension.author, extension.type, os.time(), os.time())

        -- Install hooks
        if extension.hooks then
            for hook, script in pairs(extension.hooks) do
                db:execute("INSERT INTO castro_extension_hooks (extension_id, type, script, enabled) VALUES (?, ?, ?, b'1')", extension.id, hook, script)
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

        session:setFlash("success", http.postValues.install_extension .. " has been installed.")
    elseif http.postValues.uninstall_extension then
        -- Remove extension
        db:execute("DELETE FROM castro_extensions WHERE id = ?", http.postValues.uninstall_extension)
        session:setFlash("success", http.postValues.uninstall_extension .. " has been uninstalled.")
    end

    http:redirect("/subtopic/admin/extensions/install")
end
