---
Name: events
---

# Events metatable

Allows execution of tick events

- [events:tick(file, interval)](#tick)

# tick

Starts a tick event. Events must be loaded from files with a given tick rate.

Event files must have the `function run()` declaration.

```lua
events:tick("engine/onlinechart.lua", "2m")
```

### test.lua

```lua
-- Players online chart background task
local interval = time:parseDuration(app.Custom.OnlineChart.Interval)

function run()
	local result = db:singleQuery("SELECT COUNT(*) AS count FROM players_online")
	local count, time = result.count, os.time()
	db:execute("INSERT INTO castro_onlinechart (count, time) VALUES (?, ?)", count, time)

	local old = time - (interval * (app.Custom.OnlineChart.Display + 1))
	db:execute("DELETE FROM castro_onlinechart WHERE time < ?", old)
end
```

This example illustrates how the online chart gets its data. Tick rate strings can contain `s,m,h,d,x,y` for example `1h35m`

### Stopping an event

Castro adds the `event` metatable to each event file.

- [event:stop()](#stop)

# stop

Stops the execution of the current event.

```lua
function run()
    local i = 0
    for i = 0, 10, 1 do
        if i == 5 then
            event:stop()
        end    
    end
end
```
