---
Name: cache
---

# Metatable:cache

Provides access to the application cache object

- [cache:set(key, value, duration)](#set)
- [cache:get(key)](#get)
- [cache:delete(key)](#delete)

# set

Saves the given object into the application cache. For the duration you must use a valid string such as "1h", "1s", "2h20m".

```lua
cache:set("test", 12, "4h")
```

# get

Retrieves a value from the cache object.

```lua
local data = cache:get("test")
--[[ data = 12 ]]--
```

# delete

Deletes a cache item.

```lua
cache:delete("test")
```