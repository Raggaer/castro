---
Name: socket
---

# Socket metatable

Provides access to raw socket connections.

- [socket:get(protocol, address, message)](#get)

# get

Opens a connection with the given `protocol` to the given `address`. On successful connection, `message` is written to the socket and the response is returned as a string. `address` needs to be a resolvable hostname or ip address including the destination port in the format "address:port".
`protocol` should normally be either "tcp" or "udp" but any protocol supported by net.Dial is valid. See https://pkg.go.dev/net#Dial for more details.

```lua
local data, err = socket:get("tcp", "localhost:7171", "hello")
--- data = "world"
--- err = nil
```

If an error occurs this function returns nil followed by the error message.

```lua
local data, err = socket:get("tcp", "noservice:7171", "hello")
--- data = nil
--- err = "Failed to connect socket: dial tcp [::1]:7171: connect: connection refused"
```

## Example

In this example we will request XML formatted server info from the status protocol and print the player record to the console, or log the error if unsuccessful.

We read the server's `ip` and `statusProtocolPort` settings from config.lua and use these for the `address` to connect to our own server.
To request the XML data from the status protocol we need to send 4 bytes (6, 0, 255, 255) followed by "info" as the `message`.

```lua
-- Get ip and statusProtocolPort
local address = string.format("%s:%s", config:get("ip"), config:get("statusProtocolPort"))

-- 4 bytes 6, 0, 255, 255 followed by info
local message = string.format("%c%c%c%cinfo", 6, 0, 255, 255)

-- Open the socket
local data, err = socket:get("tcp", address, message)
if err then
    -- Log the error
    return log:error(err)
end

-- Parse the xml string
local xmlData = xml:unmarshal(data)

-- Print the player record
print("Player record:", xmlData["tsqp"]["players"]["-peak"])
```
