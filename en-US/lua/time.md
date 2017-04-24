---
Name: time
---

# Metatable:time

Helper functions to work with time values

- [time:parseUnix(timestamp)](#parseunix)
- [time:parseDuration(durationString)](#parseduration)
- [time:parseDate(dateString, dateLayout)](#parsedate)

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

Duration strings can contain `s,m,h,x,d,y` for example `1h35m`

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