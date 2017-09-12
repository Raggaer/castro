---
name: base64
---

# Base64 metatable

Provides access to base64 encoding methods

- [base64:encode(string)](#encode)
- [base64:decode(string)](#decode)

# encode

Returns the base64 encoded representation of the given string.

```lua
local encoded = base64:encode("hello")
-- encoded = aGVsbG8
```

# decode

Returns the base64 decoded representation of the given string.

```lua
local decoded = base64:decode("aGVsbG8")
-- decoded = hello
```
