---
name: xml
---

# Xml metatable

Provides access to XML manipulation functions.

- [xml:vocationList](#vocationlist)
- [xml:vocationByName](#vocationbyname)
- [xml:vocationByID](#vocationbyid)
- [xml:monsterList](#monsterlist)
- [xml:monsterByName(name)](#monsterbyname)
- [xml:marshal(data)](#marshal)
- [xml:unmarshal(string)](#unmarshal)
- [xml:unmarshalFile(filename)](#unmarshalfile)

# vocationList

Returns a list of server vocations.

```lua
local vocs = xml:vocationList()
```

# vocationByName

Returns a vocation by the name.

```lua
local voc = xml:vocationByName("Sorcerer")
```

# vocationByID

Returns a vocation by its identifier.

```lua
local voc = xml:vocationByID(4)
```

# monsterList

Returns a list of monsters, retrieved from monsters.xml and each monster file.

```lua
local list = xml:monsterList()
local first = list[0].Name
--- first = "Amazon"
```

# monsterByName

Retrieves a monster by its name.

```lua
local monster = xml:monsterByName("Demon")
```

# marshal

Converts the given lua table to a valid XML string.

```lua
local data = {}

data.name = "Raggaer"
data.level = 10

local text = xml:marshal(data)
--- text = <doc><level>10</level><name>Raggaer</name></doc>
```

# unmarshal

Converts a valid XML string to a lua table.

```lua
local text = "<doc><level>10</level><name>Raggaer</name></doc>"

local data = xml:unmarshal(text)
--[[
data.name = "Raggaer"
data.level = 10
]]--
```

Inline element keys follow this structure `-element`.

# unmarshalFile

Converts a valid XML file to a lua table. This function acts almost like [unmarshal](#unmarshal) but takes a file path instead of a string.

This function is very solid, you can easily parse all your server XML files using this method, below is an example on how spells are parsed:

```lua
function get()
    local data = {}

    data.list = xml:unmarshalFile(app.Main.Datapack .. "/data/spells/spells.xml")

    for i, spell in pairs(data.list.spells.instant) do
        if spell["-script"] == nil  then
            data.list.spells.instant[i] = nil
        else
            if string.find(spell["-script"], "monster/", 1) then
                data.list.spells.instant[i] = nil
            end
        end
    end

    http:render("spells.html", data)
end
```

Inline element keys follow this structure: `-element`.

Its recommended that you cache the results, parsing an XML file on each request is not the way to go, however, [xml:unmarshalFile(filename)](#unmarshalfile) will save the XML result in the cache by default.
