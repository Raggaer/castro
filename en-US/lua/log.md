---
name: log
---

# Log metatable

Provides access to the application logging interface.

- [log:info(data)](#info)
- [log:error(data)](#error)
- [log:fatal(data)](#fatal)

# info

Logs a message with the info level.

```lua
log:info("This is a message")
```

# error

Logs a message with the error level.

```lua
log:error("This is a message")
```

# fatal

Logs a message with the fatal level. The fatal level will cause Castro to quit.

```lua
log:fatal("This is a message")
```