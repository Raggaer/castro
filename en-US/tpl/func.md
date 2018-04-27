---
name: Functions
---

# Functions

Castro provides a list of already defined functions you can use on your templates:

- [url](#url)
- [vocation](#vocation)
- [serverName](#servername)
- [serverMotd](#servermotd)
- [nl2br](#nl2br)
- [urlEncode](#urlencode)
- [urlDecode](#urldecode)
- [isDev](#isdev)
- [str2html](#str2html)
- [str2url](#str2url)
- [menuPages](#menupages)
- [eq](#eq)
- [eqNumber](#eqnumber)
- [gtNumber](#gtnumber)
- [lsNumber](#lsnumber)
- [unixToDate](#unixtodate)

# url

Creates an absolute URL for Castro. You should use this when creating links inside Castro.

```html
<a href="{{ url "subtopic" "login" }}></a> 
```

You can pass as many strings as you want to create the URL.

# vocation

Returns the name of the given vocation identifier.

```html
<p>{{ vocation 12 }}</p>
```

# serverName 

Returns the `config.lua` server name.

```html
{{ serverName }}
```

# serverMotd

Returns the `config.lua` server motd.

```html
{{ serverMotd }}
```

# nl2br

Similar to the `PHP` function. It converts all newlines to `<br>`. Useful for textareas for example.

```html
{{ nl2br .text }}
```

# urlEncode

Encodes the given value so its safe to use inside an URL.

```html
<a href="/test?name={{ urlEncode .name }}">Test</a>
```

# urlDecode

Decodes the given encoded URL string.

```html
<p>Decoded value is: {{ urlDecode "My+encoded+url" }}</p>
```

# isDev

Checks if Castro runs on development mode.

```html
{{ if isDev }}
<p>Development</p>
{{ end }}
```

# str2html

Converts the given string to HTML output. By default Castro sanitizes all HTML output for variables.

```html
{{ str2html .text }}
```

So in order to be able to use HTML elements from your variables (like for example a news article where you use HTML tags) you need to call this function.

# str2url

Converts the given string to a safe URL output. By default Castro sanitizes URL values to prevent any kind of attacks.

```html
<a href="/test/{{ str2url "hello }}">
```

# eq

Compares two objects and return true if they are equal.

```html
{{ if eq "hello" "world" }}
    <p>Not equal</p>
{{ end }}
```

# eqNumber

Compares two numbers a  == b.

```html
{{ if eqNumber 5 4 }}
    <p>Not equal</p>
{{ end }}
```

# gtNumber

Compares two numbers a > b.

```html
{{ if gtNumber 5 1 }}
    <p>Greater</p>
{{ end }}
```

# lsNumber

Compares two numbers a < b.

```html
{{ if lsNumber 1 5 }}
    <p>Less</p>
{{ end }}
```

# unixToDate

Returns a date string for the given UNIX timestamp.

```html
{{ $d := unixToDate 1521285384 }}
<p>Current date: {{ $d }}</p>
```