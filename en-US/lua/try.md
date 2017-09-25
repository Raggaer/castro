---
name: try
---

# Try and catch

Castro supports try/catch blocks using lua functions. The error is passed to teh catch function.

```lua
try(tryfunction(), catchfunction(err))
```

Below is a simple example of the try/catch function:

```lua
try(
    function()
        print(1 + "aa")
    end,
    function(err)
        print("Something went wrong: " .. err)
    end
)
```