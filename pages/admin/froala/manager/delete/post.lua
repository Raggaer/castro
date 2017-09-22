function post()
    if not session:isAdmin() then
        return
    end

    if file:exists("public" .. http.postValues.src) then
        os.remove("public" .. http.postValues.src)
    end
end