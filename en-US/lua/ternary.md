---
Name: ternary
---

# Ternary operator

Castro supports some sort of ternary condition using the function `ternary(condition, success, fail)`.

Below is an example of the ternary operator function:

```lua
local test = false
local v = ternary(test, 0, 100)
--- v = 100
```

This method is useful for form values for example:

```lua
local testValue = ternary(http.postValues["test"] == nil, "no", "yes")
```