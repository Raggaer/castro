---
Name: http
---

# Http metatable

Provides access to HTTP related functions.

- [http.method](#method)
- [http.subtopic](#subtopic)
- [http.body](#body)
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
- [http:formFile(name)](#formfile)
- [http:setCookie(name, value, expiration)](#setcookie)
- [http:getCookie(name)](#getcookie)
- [http:getRelativeURL()](#getrelativeurl)

# method

Holds the incoming request method.

```lua
local method = http.method
-- method = "GET"
```

# subtopic

Holds the current subtopic uri.

```lua
-- example.com/subtopic/test

local subtopic = http.subtopic
-- subtopic = "/subtopic/test"
```

# body

Holds the incoming request body, useful for creating a JSON API. Will be an empty string if there is no body attached.

```lua
local body = http.body
-- body = "{data}"
```

This can be used to handle JSON or XML requests, for example the Tibia 11 client webservice sends some JSON data to the server, we can work with that data like this:

```lua
function get()
  local data = json:unmarshal(http.body)
  print(data.password)
end
```

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

# formFile

Retrieves a file from a `multipart/form-data` encoded form. If the file is not present this function will return `nil`.

```lua
local file = http:formFile("guild-image")
```

This function returns a metatable with the following functions:

- [formFile:contentType()](#contentype)
- [formFile:isValidExtension(type)](#isvalidextension)
- [formFile:isValidPNG()](#isvalidpng)
- [formFile:saveFile(destination)](#destination)
- [formFile:saveFileAsPNG(destination, width, height)](#savefileaspng)
- [formFile:saveFileAsJPEG(destination, quality, width, height)](#savefileasjpeg)
- [formFile:getFile()](#getfile)

# contentType

Returns the file content type.

```lua
local file = http:formFile("guild-image")

local c = file:contentType()
-- c = "image/png"
```

# isValidExtension

Checks if the current form file is a valid file extension. File extensions are actually `MIME types`.

```lua
local file = http:formFile("guild-image")

local isPNG = file:isValidExtension("image/gif")
-- isPNG = false
```

# isValidPNG

Checks if the current form file is a valid `png` image.

```lua
local file = http:formFile("guild-image")

local isPNG = file:isValidPNG()
-- isPNG = true
```

# saveFile

Saves the current form file to the given destination. This method can be used for any file extension,
if the file is an image you should use one of the save as image method instead (this methods are more secure to use for images since files are decoded/encoded before beeing saved)

```lua
local file = http:formFile("guild-image")

file:saveFile("public/images/guild-image.png")
```

# saveFileAsPNG

Saves the current form file as a `png` image to the given destination. You can pass a width and a height as optional values to resize the image.

```lua
local file = http:formFile("guild-image")

file:saveFileAsPNG("public/images/guild-image.png", 64, 64)
```

# saveFileAsJPEG

Saves the current form file as a `jpeg` image to the given destination. You can pass the quality you want the image to be saved (as a default value images are saved with 50 quality) you can also pass a width and a height as optional values to resize the image.

```lua
local file = http:formFile("guild-image")

file:saveFileAsPNG("public/images/guild-image.png", 100, 64, 64)
```

# getFile

Returns the file byte array as a lua string.

```lua
local file = http:formFile("guild-image")

local f = file:getFile()
-- f = <file byte array>
```

# setCookie

Sets a HTTP cookie. You need to pass a name, a value (string) and an expiration value (date in seconds).

```lua
http:setCookie("Hello", "World", os.time() + 5 * 60)
```

The example above will set a cookie named `Hello` with a value `World` that will expire in 5 minutes.

# getCookie

Returns a HTTP cookie by the given name.

```lua
local c = http:getCookie("Hello")
-- c = "World"
```

If the cookie is not set (or expired) this method will return `nil`.

```lua
local c = http:getCookie("Test")
-- c = nil
```

# getRelativeURL

Returns the current relative URL of the request (similar to `http.subtopic`):

```lua
local u = http:getRelativeURL()
-- u = "/subtopic/test?test=test"
```

