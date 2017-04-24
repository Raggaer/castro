---
Name: env
---

# Metatable:env

Provides access environment variables

- [env:set(key, value)](#set)
- [env:get(key)](#get)

# set

Sets the given environment variable

```lua
env:set("my_value", "test")
```

# get

Retrieves the given environment variable

```lua
local data = env:get("my_value")
--[[ data = test ]]--
```