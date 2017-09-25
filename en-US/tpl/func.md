---
name: Functions
---

# Functions

Castro provides a list of already defined functions you can use on your templates:

- [url](#url)
- [vocation](#vocation)
- [serverName](#servername)
- [serverMotd](#serverMotd)
- [nl2br](#nl2br)
- [urlEncode](#urlencode)
- [isDev](#isdev)
- [str2html](#str2html)

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