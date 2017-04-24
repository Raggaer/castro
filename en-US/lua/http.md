---
Name: http
---

# Http metatable

Provides access to HTTP related functions.

- [http:redirect(url, header)](#redirect)
- [http:render(template, data)](#render)
- [http:write(string)](#write)
- [http:serveFile(path)](#servefile)
- [http:get(url)](#get)
- [http:postForm(url, data)](#postform)
- [http:setHeader(key, value)](#setheader)
- [http:getHeader(key)](#getheader)
- [http:getRemoteAddress()](#getremoteaddress)

# redirect

Redirects the user to the given location. You can provide an optional header. By default all redirects are done using a `302` header.

```lua
http:redirect("/subtopic/test")
```

Redirecting an user does not stop the execution of the page. You must return on each redirect.

# render

Renders the given template. You can pass a lua table for the template data. This data is HTML escaped by the templating engine.

```lua
local data = {}

data.name = "Raggaer"

http:render("test.html", data)
```

Rendering a template always set the status header as `200`. Rendering does not stop the execution of the page.

# write

Outputs a raw string to the response writer. You can use this to output JSON data (or even make a REST API).

```lua
http:write("Some data")
```

Writing does not stop the execution of the page.

# serveFile

Serves the given file. Serving a file does not stop the execution of the page.

```lua
http:serveFile("path/to/my/file.png")
```

# get

Performs a HTTP GET request to the given destination. Very basic method, for more control over the request use [http:curl(data)](#curl).

```lua
local response = http:get("www.google.com")
-- response = "<html>...</html>"
```

# postForm

Performs a HTTP POST request to the given destination. 

```lua
local data = {}

data.name = "Raggaer"

local response = http:postForm("www.test.com", data)
```

# setHeader

Sets a header for the current running request.

```lua
http:setHeader("Engine", "Castro")
```

# getHeader

Retrieves a header from the current running request.

```lua
local value = http:getHeader("Engine")
-- value = "Castro"
```

# getRemoteAddress

Retrieves the remote address of the current running request.

```lua
local addr = http:getRemoteAddress()
-- addr = "127.0.0.1"
```

# curl

Provides access to an extensible request creator. You can set headers, authentication and data. These are the table fields:

- timeout: request timeout in seconds. Optional.
- method: request method. 
- url: request destination.
- data: request content table. Optional.
- headers: request headers table. Optional.
- authentication: request authentication table. Optional.

```lua
local data = {}

data.timeout = 10
data.method = "get"
data.url = "www.test.com"

data.data = {}
data.data.name = "Raggaer"

data.headers = {}
data.headers.Engine = "Castro"

data.authentication = {}
data.username = "Raggaer"
data.password = "password"

local response, headers, status = http:curl(data)
--[[
response = "<html>...</html>"
headers = {}
status = 200
]]--
```

The function will return the request response, headers (as a table) and the status code.