local function fetchValues(voc)
    local v = app.Custom.NewCharacterValues[voc]
    if v.inherit_from ~= nil then
        v = app.Custom.NewCharacterValues[v.inherit_from]
    end
    return v
end

function post()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    if http.postValues["character-name"]:len() < 5 or http.postValues["character-name"]:len() > 12 then
        session:setFlash("validation-error", "Invalid character name. Names can only have 5 to 12 characters")
        http:redirect()
        return
    end

    if not validator:validGender(http.postValues["character-gender"]) then
        session:setFlash("validation-error", "Invalid character gender. Gender not found")
        http:redirect()
        return
    end

    if not validator:validVocation(http.postValues["character-vocation"]) then
        session:setFlash("validation-error", "Invalid character vocation. Vocation not found")
        http:redirect()
        return
    end

    if not validator:validTown(http.postValues["character-town"]) then
        session:setFlash("validation-error", "Invalid character town. Town not found")
        http:redirect()
        return
    end

    if not validator:validUsername(http.postValues["character-name"]) then
        session:setFlash("validation-error", "Invalid character name format. Only letters A-Z and spaces allowed")
        http:redirect()
        return
    end

    if db:query("SELECT id FROM players WHERE name = ?", http.postValues["character-name"]) ~= nil then
        session:setFlash("validation-error", "Character name already in use")
        http:redirect()
        return
    end

    local account = session:loggedAccount()

    if db:singleQuery("SELECT COUNT(*) as total FROM players WHERE account_id = ?", account.ID).total > 5 then
        session:setFlash("validation-error", "You can only have 5 characters")
        http:redirect()
        return
    end

    ncv = fetchValues(tonumber(xml:vocationByName(http.postValues["character-vocation"]).ID)+1) -- new char values

    local eStr = "" -- extra string
    local vStr = "" -- value string

    if ncv.extra ~= nil then
        for k, value in pairs(ncv.extra) do
            eStr = eStr ..tostring(k)..", "
            vStr = vStr .. value .. ", "
        end
    end

    db:execute(
        "INSERT INTO `players` ("..eStr.."name, sex, account_id, vocation, town_id, conditions, level, health, healthmax, mana, manamax, cap, soul) VALUES ("..vStr.."?,?,?,?,?,'',?,?,?,?,?,?,?)",
        http.postValues["character-name"],
        http.postValues["character-gender"],
        account.ID,
        xml:vocationByName(http.postValues["character-vocation"]).ID,
        otbm:townByName(http.postValues["character-town"]).ID,
        ncv.level,
        ncv.health,
        ncv.healthmax,
        ncv.mana,
        ncv.manamax,
        ncv.cap,
        ncv.soul
    )

    session:setFlash("success", "Character created")
    http:redirect("/subtopic/account/dashboard")
end
