---
name: crypto
---

# Crypto metatable

Provides access to cryptography methods

* [crypto:sha1(string)](#sha1)

# sha1

Returns the sha1 hash of the given string

```lua
local hash = crypto:sha1("hello")
--[[ hash = aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d ]]--
```