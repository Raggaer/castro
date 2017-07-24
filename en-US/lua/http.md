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
- [http:curl(data)](#curl)

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

The function will return the request response, headers (as a table) and the status code.

```lua
local request = {}

request.method = "get"
request.url = "www.google.com"

response, headers, status = http:curl(request)

--[[
response = <html>...</html>
headers = { "Google-Header": "Test", ... }
status = 200
]]--
```

These are the data fields you can pass to the cURL client:

- [timeout](#curl.timeout)
- [method](#curl.method) - mandatory
- [url](#curl.url) - mandatory
- [data](#curl.data)
- [headers](#curl.headers)
- [authentication](#curl.authentication)

# curl.timeout

Time to wait before the client times out during a request, by default this field is 0.

# curl.method

The method to use while performing the request. This field is mandatory.

# curl.url

The destination site. This field is mandatory.

# curl.data

The data that is going to be sent with the request. This data needs to be a lua table, the table is converted to a slice of URL values.

```lua
local request = {}

request.data = {}

equest.data.username = "Raggaer"
equest.data.password = "Test1234"
```

# curl.headers

Table of headers to use during the request. Each element of the table defines a header where the table key is the header key and the table value is the header value.

```lua
local request = {}

request.headers = {}

request.headers["My-Custom-Header"] = "Hello World"
```

# curl.authentication

Basic HTTP authentication headers. You must provide an username and a password field using a lua table.

```lua
local request = {}

request.authentication = {}

request.authentication.username = "Raggaer"
request.authentication.password = "Test1234"
```