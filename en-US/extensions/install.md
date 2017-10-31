---
name: Install script
---

# Install script

The install script is optional. If your extension needs to execute some code before it can be used, this is the place to do it. To include an install script is easy. All you have to do is add a file named `install.lua` in your extension folder. When the extension is installed, Castro will automatically load the file and call the `install` function.

## Return values

```lua
function install()
    return success, message
end
```

The install function may return a boolean value representing the success state, plus an optional message string. The first return value should be `true` for success and `false` for failure. If the return values are omitted, i.e return is nil, success is assumed.
If present, the message will be displayed on the success or error alert on the installation page.

## Examples

A common use case will be to add custom config values or perform required alterations to the database befor the extension can be used.

```lua
function install()
    -- Run some code
    return true, "We have success!"
end
```

# Uninstall script

Works the same as the install script, except the file and function should be named `uninstall.lua` and `uninstall` respectively. The only difference is that an uninstall script can not abort uninstallation by returning false. If the function return false, the extension will be uninstalled from Castro but an error message will be shown to the user. This is to prevent an extension from becoming unremovable.

## Examples

The uninstall script will mainly be used to clean up things that are no longer needed when the extension is no longer installed. For example to remove custom config values that are no longer used.

```lua
function uninstall()
    -- Run some code
    return true, "Cleanup done!"
end
```
