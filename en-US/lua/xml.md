---
name: xml
---

# Xml metatable

Provides access to XML manipulation functions.

- [xml:marshal(data)](#marshal)
- [xml:unmarshal(string)](#unmarshal)
- [xml:unmarshalFile(filename)](#unmarshalfile)

# marshal

Convers the given lua table to a valid XML string.

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

# unmarshalFile

Converts a valid XML file to a lua table. This function acts almost like [unmarshal](#unmarshal) but takes a file path instead of a string.