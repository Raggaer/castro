---
Name: captcha
---

# Captcha metatable

Provides access to Google reCAPTCHA functions. You must have a valid private and public key on your `config.toml`.

- [captcha:isEnabled()](#isenabled)
- [captcha:verify(data)](#verify)

# isEnabled

Checks if the captcha system is enabled.

```lua
local enabled = captcha:isEnabled()
-- enabled = true
```

# verify

Verifies that the given answer is valid.

```lua
local good = captcha:verify(answer_text)
-- good = true
```

Usually the captcha answer comes from a `POST` form. Check `register.lua` to see a live example.