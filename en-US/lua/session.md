---
name: session
---

# Session metatable

Provides access to the current request session interface.

- [session:isLogged()](#islogged)
- [session:isAdmin()](#isadmin)
- [session:setFlash(key, value)](#setflash)
- [session:getFlash(key)](#getflash)
- [session:set(key, value)](#set)
- [session:get(key)](#get)
- [session:destroy()](#destroy)
- [session:loggedAccount()](#loggedaccount)

# isLogged

Checks if the current user is logged-in or not.

```lua
local logged = session:isLogged()
-- logged = false
```

# isAdmin

Checks if the current user has admin privileges.

```lua
local admin = session:isAdmin()
-- admin = false
```

# setFlash

Sets a session flash value. Flash values are cleared when accessing them. You usually use these for form validation errors for example.

```lua
session:setFlash("error", "Invalid username")
```

# getFlash

Retrieves a session flash value. Flash values are cleared when accessing them. You usually use these for form validation errors for example.

```lua
local error = session:getFlash("error")
-- error = "Invalid username"
```

# set

Sets a session value.

```lua
session:set("name", "Raggaer")
```

# get

Retrieves a session flash value

```lua
local name = session:get("name")
-- name = "Raggaer"
```

# destroy

Kills a session. This will clear all values from it.

```lua
session:destroy()
```

# loggedAccount

Returns the current logged account table.

```lua
local acocunt = session:loggedAccount()
--[[
account.ID = 1
account.Name = "xxxx"
account.Password = "xxx"
account.Premium_ends_at = 1200000
account.Creation = 1200000
account.Secret = "xx"

account.Castro.ID = 1
account.Castro.Points = 10
account.Castro.Admin = false
]]--
```