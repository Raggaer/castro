---
Name: cache
---

# Cache metatable

Provides access to the application cache instance.

- [cache:set(key, value, duration)](#set)
- [cache:get(key)](#get)
- [cache:delete(key)](#delete)

# set

Saves the given object into the application cache. For the duration you must use a valid string such as "1h", "1s", "2h20m".

The cache is stored on your system memory.

```lua
cache:set("test", 12, "4h")
```

# get

Retrieves a value from the cache object.

```lua
local data = cache:get("test")
--[[ data = 12 ]]--
```

If the item does not exists this function will return nil

```lua
local data = cache:get("does_not_exist")
--- data = nil
```

# delete

Deletes a cache item.

```lua
cache:delete("test")
```