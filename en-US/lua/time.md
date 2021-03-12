---
Name: time
---

# Metatable:time

Helper functions to work with time values

- [time:newDuration(nanoseconds)](#newduration)
- [time:parseUnix(timestamp)](#parseunix)
- [time:parseDuration(durationString)](#parseduration)
- [time:parseDate(dateString, dateLayout)](#parsedate)

# newDuration

Creates a duration table from the given time, the number you pass to the function needs to be in nanoseconds. The result is a table that contains the following table:

- Nanoseconds
- Seconds
- Minutes
- Hours

```lua
local t = time:newDuration(86400000000000)
--[[
t.Nanoseconds = 86400000000000
t.Seconds = 86400
t.Minutes = 1440
t.Hours = 24
]]--
```

Below is an example on how we parse some config values using `math.pow` to convert miliseconds to nanoseconds:

```lua
local timeToDecreaseFrags = time:newDuration(config:get("timeToDecreaseFrags") * math.pow(10, 6))
local whiteSkullTime = time:newDuration(config:get("whiteSkullTime") * math.pow(10, 6))
```

# parseUnix

Parses the given unix timestamp and returns a castro time table. This table contains the following fields

- Year
- Month
- Day
- Hour
- Minute
- Second
- Result

The result field is a time string with the following format `Mon Jan 2 15:04:05`

```lua
local parsed = time:parseUnix(1487974045)
--[[
parsed.Year = 2017
parsed.Month = 2
parsed.Day = 24
parsed.Hour = 23
parsed.Minute = 08
parsed.Second = 30
parsed.Result = Fri Feb 24 23:08:30
]]--
```

# parseDuration

Parses the given duration string and returns the duration seconds.

Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. For example `1h35m`.

```lua
local seconds = time:parseDuration("1h")
-- seconds = 3600
```

# parseDate

Parses the given date using a layout

```lua
local seconds = time:parseDate("Fri Feb 24 23:08:30", "Mon Jan 2 15:04:05")
-- seconds = 1487974045
```