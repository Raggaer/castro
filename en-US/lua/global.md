---
name: global
---

# Global metatable

Easy access to global variables across your system. Values are saved into the database. Currently only lua tables are supported.

- [global:set(key, value)](#set)
- [global:get(key)](#get)
- [global:delete(key)](#delete)

# set

Saves a global value into the database. You need to provide a key and a table to be saved. If the key already exists it will be updated.

```lua
local data = {}
data.Test = 10

global:set("Test", data)
```

# get

Retrieves a global value from the database. You need to provide a key. If the key does not exist this method will return nil.

```lua
local data = global:get("Test")
--- data.Test = 10
```

```lua
local data = global:get("NoExists")
--- data = nil
```

# delete

Deletes a global value from the database. You need to provide a key. If the key does not exist nothing will happen.

```lua
global:delete("Test")

local data = global:get("Test")
--- data = nil
```
