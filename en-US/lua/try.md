---
name: try
---

# Try and catch

Castro supports try/catch blocks using lua functions. The syntax is the following

```lua
try(tryfunction(), catchfunction())
```

Below is a simple example of the try/catch function:

```lua
try(
    function()
        print(1 + "aa")
    end,
    function()
        print("Something went wrong!")
    end
)
```