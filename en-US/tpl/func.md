---
name: Functions
---

# Functions

Castro provides a list of already defined functions you can use on your templates:

- [vocation](#vocation)
- [serverName](#servername)
- [serverMotd](#serverMotd)
- [nl2br](#nl2br)

# vocation

Returns the name of the given vocation identifier.

```html
<p>{{ vocation 12 }}</p>
```

# serverName 

Returns the `config.lua` server name

```html
{{ serverName }}
```

# serverMotd

Returns the `config.lua` server motd

```html
{{ serverMotd }}
```

#nl2br

Similar to the `PHP` function. It converts all newlines to `<br>`. Useful for textareas for example.

```html
{{ nl2br .text }}
```

