---
name: Pages
---

# Pages

Pages in an extension should be placed in a folder named `pages`. Folders inside `pages` will be mapped to a URL based on the name of the folder it is in.
A page in `my-extension/pages/path/to/page`, will correspond to `example.com/subtopic/path/to/page`.
You can also add a new page to an existing folder. For example if you want to add a new page to `subtopic/community`, you just have to place your page in `my-extension/pages/community/my-page`.

## Files
Pages consists of 2 types of files, html for templates and Lua for logic. The Lua and html files for extension pages works the same way as [custom pages](/docs/info/pages).

### Folder structure

```html
my_extension /
    / pages
        / my-page
            get.lua
            my-page.html
```

## Overriding defaults

In addition to adding new pages, you may also override any default page or template.

### Template
If you wish to override the template of a page only, all you have to do is put a html file with the exact same name as the file one you wish to replace in your `pages` folder. To replace the home page template, you would put your own template in `my-extension/pages` and name it `home.html`. Any variables passed by the default Lua files will be available to your template.

Keep in mind that since you can overwrite any existing template by adding your own version, you might also overwrite a template by accident if you use a too generic name for your template. To avoid accidental overwrites it is a good idea to include, for example, your extenion's id in the filename, `my-extension.my-page.html`.

### Logic

To replace the logic (Lua files) of a page, you need to put your own files in a folder corresponding to the path of the page you wish to override. To override the home page, you would put your Lua files in `my-extension/pages/home`.
