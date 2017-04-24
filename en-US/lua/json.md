---
Name: json
---

# Json metatable

Provides access to json manipulation functions.

- [json:marshal(table)](#marshal)
- [json:unmarshal(string)](#unmarshal)

# marshal

Converts the given lua table to a valid JSON string.

```lua
local data = {}

data.name = "Raggaer"
data.level = "80"

local text = json:marshal(data)

-- text = {"level":"80","name":"Raggaer"}
```

# unmarshal

Converts a valid JSON string to a lua table.

```lua
local data = json:unmarshal("{\"level\":\"80\",\"name\":\"Raggaer\"}")
--[[
data.level = 80
data.name = "Raggaer"
]]--
```