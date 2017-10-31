---
name: Manifest
---

# Extension manifest

The manifest file contains some information about the extension and should be named `extension.json` and placed at the root of the extension's folder.

The manifest can contain the following fields:

- author: Name of the extension author.
- name: The extension's user-friendly name.
- id: Should be a uniquely identifiable name for the extension and must match the name of the extension's folder.
- version: Version of the extension, to be used for extension updating.
- description: A description of the extension.
- type: A type/category for your extension.
- url (optional): A URL associated with your extension.
- hooks (optional): An object where each key is the hook name and the value is the script to execute.
- templateHooks (optional): same as hooks except the file should be a html template.

## Example extension.json
```json
{
    "author":"Castro",
    "name":"Hello world page",
    "id":"page.helloworld",
    "version":"1",
    "description":"This is a sample page extension.",
    "type":"page",
    "url":"https://docs.castroaac.org",
    "hooks":{
        "onStartup":"my-script.lua"
    },
    "templateHooks":{
        "head":"my-template.html"
    }
}
```
