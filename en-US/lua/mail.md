---
Name: mail
---

# Mail metatable

Provides access to mail sending functions. You must have configured a mail server in your `config.toml` file.

- [mail:send(info_table)](#send)

# send

Sends an email with the provided information. The information must be provided on a lua table with the following content:

- to: email destination. Who to send the email to.
- subject: email subject.
- body: email body. You can use HTML.

```lua
local data = {}

data.to = "test@test.com"
data.subject = "Welcome!"
data.body = "<h1>Welcome<h1><p>Hello!</p>"

mail:send(data)
```