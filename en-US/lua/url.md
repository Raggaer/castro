---
Name: url
---

# Url metatable

Provides access to url manipulation functions

- [url:decode(uri)](#decode)
- [url:encode(raw)](#encode)

# decode

Decodes the given encoded url value. Usually you should decode all non-numeric `GET` params

```lua
local value = url:decode("My+name+is+raggaer")
-- value = "My name is raggaer"
```

# encode

Encodes the given value. Usually you must encode all non-numeric `GET` params

```lua
local encoded = url:encode("My name is raggaer")
-- encoded = "My+name+is+raggaer"
```