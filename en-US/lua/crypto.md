---
name: crypto
---

# Crypto metatable

Provides access to cryptography methods

- [crypto:sha1(string)](#sha1)
- [crypto:sha256(string)](#sha256)
- [crypto:hmacsha256(secret, message)](#sha1)
- [crypto:md5(string)](#md5)
- [crypto:randomString(length)](#randomstring)
- [crypto:qr(code)](#qr)
- [crypto:qrKey()](#qrkey)

# sha1

Returns the sha1 hash of the given string.

```lua
local hash = crypto:sha1("hello")
-- hash = aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
```

# sha256

Returns the sha256 hash of the given string.

```lua
local hash = crypto:sha256("hello")
-- hash = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```

# hmacsha256

Returns the HMAC-sha256 hash of the given secret and message.

```lua
local hash = crypto:hmacsha256("secret", "hello")
-- hash = 88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b
```

# md5

Returns the md5 hash of the given string.

```lua
local hash = crypto:md5("hello")
-- hash = 5d41402abc4b2a76b9719d911017c592
```

# randomString

Creates and returns a random string with the given length. The string will use any of this characters:

`abcdefghijklmnropqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`

```lua
local r = crypto:randomString(3)
-- r = aT1
```

# qr

Generates a base64 encoded QR image representing the given data.

```lua
local img = crypto:qr("Test")
```

You can pass it to a template and render it as an image.

```html
<img src="{{ .img }}">
```

# qrKey

Creates and returns a 16 length random string using the any of this characters:

`abcdefghijklmnropqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`

```lua
local r = crypto:generateSecretKey()
-- r = abHjclkOp18Jh7fg
```
