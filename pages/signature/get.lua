function get()
    local character = Player(http.getValues.name)

    if character == nil then
        http:redirect("/")
        return
    end

    if file:exists("public/images/signature/" .. character:getName() .. ".png") and not app.Mode == "dev" then
       local up = file:mod("public/images/signature/" .. character:getName() .. ".png") + 5 * 60
       if up > os.time() then
           http:serveFile("public/images/signature/" .. character:getName() .. ".png")
           return
       end
    end

    local online = db:singleQuery("SELECT 1 FROM players_online WHERE player_id = ?", character.id)
    local img = image:new(500, 150)

    img:setBackground("public/images/signature-bg.png")

    if online == nil then
        img:writeText(character:getName(), "#FF0000", 19, 20, 20)
    else
        img:writeText(character:getName(), "#34BC41", 19, 20, 20)
    end

    img:writeText("Level: " .. character:getLevel(), "#000000", 16, 20, 60)
    img:writeText("Vocation: " .. character:getVocation().Name, "#000000", 16, 20, 80)
    img:writeText("Town: " .. character:getTown().Name, "#000000", 16, 20, 100)
    img:save("public/images/signature/" .. character:getName() .. ".png")

    http:serveFile("public/images/signature/" .. character:getName() .. ".png")
end