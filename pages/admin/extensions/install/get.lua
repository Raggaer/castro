require "util"

function get()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}
    data.validationError = session:getFlash("validationError")
    data.success = session:getFlash("success")

    for _, extensionDirectory in pairs(file:getDirectories("extensions")) do
        data.extensions = data.extensions or {}
        data.extensions[extensionDirectory] = json:unmarshalFile(string.format("extensions/%s/extension.json", extensionDirectory))
    end

    local installedExtensions = db:query("SELECT * FROM castro_extensions")
    if installedExtensions and data.extensions then
        for _, extension in pairs(installedExtensions) do
            data.extensions[extension.id].installed = extension.installed
        end
    end

    http:render("installextension.html", data)
end
