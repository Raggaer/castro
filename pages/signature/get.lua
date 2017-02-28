function get()
    local character = db:singleQuery("SELECT id, name, level, vocation, town_id FROM players WHERE name = ?", http.getValues.name)

    if character == nil then
        http:redirect("/")
        return
    end

    if file:exists("public/images/signature/" .. character.name .. ".png") then
       local up = file:mod("public/images/signature/" .. character.name .. ".png") + 5 * 60
       if up > os.time() then
           http:serveFile("public/images/signature/" .. character.name .. ".png")
           return
       end
    end

    local online = db:singleQuery("SELECT 1 FROM players_online WHERE player_id = ?", character.id)

    local img = image:new(500, 150)

    img:setBackground("public/images/signature-bg.png")

    if online == nil then
        img:writeText(character.name, "#C64342", 19, 20, 20)
    else
        img:writeText(character.name, "#34BC41", 19, 20, 20)
    end

    img:writeText("Level: " .. character.level, "#000000", 16, 20, 60)
    img:writeText("Vocation: " .. xml:vocationByID(character.vocation).Name, "#000000", 16, 20, 80)
    img:writeText("Town: " .. otbm:townByID(character.town_id).Name, "#000000", 16, 20, 100)
    img:save("public/images/signature/" .. character.name .. ".png")

    http:serveFile("public/images/signature/" .. character.name .. ".png")
end