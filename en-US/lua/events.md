---
Name: events
---

# Events metatable

Allows execution of backrground tasks:

- [events:new(function)](#new)

# new

Starts a new event. All events run on a different thread.

```lua
events:new(
    function()
        print("Hello World")
    end
 )
```

### Full example

Below is the `online_chart` extension event. We do all the logic inside a  `while true` loop, we also call `sleep` to make some sort of tick for the loop.

```lua
events:new(
    function()
        local interval = time:parseDuration(app.Custom.OnlineChart.Interval)
        local result = db:singleQuery("SELECT COUNT(*) AS count FROM players_online")
        local count, time = result.count, os.time()
        db:execute("INSERT INTO castro_onlinechart (count, time) VALUES (?, ?)", count, time)
    
        local old = time - (interval * (app.Custom.OnlineChart.Display + 1))
        db:execute("DELETE FROM castro_onlinechart WHERE time < ?", old)
        sleep(app.Custom.OnlineChart.Interval)
    end
)
```